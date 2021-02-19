package keeper

import (
	"github.com/MinterTeam/mhub/chain/x/minter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const gweiInEth = 1e9

var minDepositAmount = sdk.NewInt(100)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper      Keeper
	bankKeeper  types.BankKeeper
	peggyKeeper types.PeggyKeeper
}

func (a *AttestationHandler) SetPeggyKeeper(keeper types.PeggyKeeper) {
	a.peggyKeeper = keeper
}

// Handle is the entry point for Attestation processing.
func (a *AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.MinterClaim) error {
	switch claim := claim.(type) {
	case *types.MsgDepositClaim:
		if claim.Amount.LT(minDepositAmount) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is too small to be deposited")
		}

		token := types.MinterCoin{
			Amount: claim.Amount,
			CoinId: claim.CoinId,
		}
		coin := token.PeggyCoin(ctx, a.keeper.oracleKeeper)
		if _, err := types.ValidatePeggyCoin(coin, ctx, a.keeper.oracleKeeper); err != nil {
			return sdkerrors.Wrapf(err, "coin is not valid")
		}
		vouchers := sdk.Coins{coin}
		if err := a.bankKeeper.MintCoins(ctx, types.ModuleName, vouchers); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", vouchers)
		}

		addr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid receiver address")
		}

		// pay commissions
		{
			valset := a.keeper.GetCurrentValset(ctx)
			commission := sdk.NewCoin(coin.Denom, coin.Amount.QuoRaw(100)) // total commission
			vouchers = sdk.Coins{coin.Sub(commission)}

			if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress{}, sdk.Coins{commission}); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}

			var totalPower uint64
			for _, val := range valset.Members {
				totalPower += val.Power
			}

			for _, val := range valset.Members {
				amount := commission.Amount.Mul(sdk.NewIntFromUint64(val.Power)).Quo(sdk.NewIntFromUint64(totalPower))
				_, err := a.keeper.AddToOutgoingPool(ctx, sdk.AccAddress{}, val.MinterAddress, "#commission", sdk.NewCoin(commission.Denom, amount))
				if err != nil {
					return sdkerrors.Wrap(err, "commission withdrawal")
				}
			}
		}

		if err = a.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, vouchers); err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}

		depositEvent := sdk.NewEvent(
			types.EventTypeBridgeDepositReceived,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTxHash, claim.TxHash),
		)
		ctx.EventManager().EmitEvent(depositEvent)

	case *types.MsgSendToEthClaim:
		if claim.Amount.LT(minDepositAmount) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is too small to be deposited")
		}

		receiver := sdk.AccAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		if err := a.Handle(ctx, att, &types.MsgDepositClaim{
			EventNonce:     claim.EventNonce,
			CoinId:         claim.CoinId,
			Amount:         claim.Amount,
			MinterSender:   claim.MinterSender,
			CosmosReceiver: receiver.String(),
			Orchestrator:   claim.Orchestrator,
			TxHash:         claim.TxHash,
		}); err != nil {
			return sdkerrors.Wrap(err, "deposit claim")
		}

		denom, err := a.keeper.oracleKeeper.GetCoins(ctx).GetDenomByMinterId(claim.CoinId)
		if err != nil {
			return sdkerrors.Wrap(err, "coin not found")
		}

		commission := sdk.NewCoin(denom, claim.Amount.QuoRaw(100))

		fee := sdk.NewCoin(denom, claim.Fee)
		feeIsOk := false

		coinPrice, err := a.keeper.oracleKeeper.GetMinterPrice(ctx, claim.CoinId)
		if err != nil {
			return sdkerrors.Wrap(err, "fee")
		}

		gasPrice, err := a.keeper.oracleKeeper.GetEthGasPrice(ctx)
		if err != nil {
			return sdkerrors.Wrap(err, "gas price")
		}

		ethPrice, err := a.keeper.oracleKeeper.GetEthPrice(ctx)
		if err != nil {
			return sdkerrors.Wrap(err, "eth price")
		}

		totalUsdCommission := fee.Amount.Mul(coinPrice).Quo(a.keeper.oracleKeeper.GetPipInBip())
		totalUsdGas := gasPrice.Mul(ethPrice).MulRaw(int64(a.keeper.oracleKeeper.GetMinSingleWithdrawGas(ctx))).QuoRaw(gweiInEth).QuoRaw(a.keeper.oracleKeeper.GetGasUnits())
		if totalUsdCommission.GTE(totalUsdGas) {
			feeIsOk = true
		}

		if claim.Amount.LTE(claim.Fee) {
			feeIsOk = false
		}

		if a.keeper.IsStopped(ctx) {
			feeIsOk = true
		}

		if feeIsOk {
			_, err := a.peggyKeeper.AddToOutgoingPool(ctx, receiver, claim.EthReceiver, claim.MinterSender, claim.TxHash, sdk.NewCoin(denom, claim.Amount).Sub(commission).Sub(fee), fee)
			if err != nil {
				return sdkerrors.Wrap(err, "withdraw")
			}
		} else {
			refundEvent := sdk.NewEvent(
				types.EventTypeRefund,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
				sdk.NewAttribute(types.AttributeKeyTxHash, claim.TxHash),
			)
			ctx.EventManager().EmitEvent(refundEvent)
			_, err := a.keeper.AddToOutgoingPool(ctx, receiver, claim.MinterSender, claim.TxHash, sdk.NewCoin(denom, claim.Amount).Sub(commission))
			if err != nil {
				return sdkerrors.Wrap(err, "refund")
			}
		}

	case *types.MsgWithdrawClaim:
		return a.keeper.OutgoingTxBatchExecuted(ctx, claim.BatchNonce)
	case *types.MsgValsetClaim:
		return a.keeper.ValsetExecuted(ctx, claim.ValsetNonce)

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
