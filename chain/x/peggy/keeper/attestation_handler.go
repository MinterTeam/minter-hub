package keeper

import (
	minterkeeper "github.com/MinterTeam/mhub/chain/x/minter/keeper"
	oracletypes "github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/MinterTeam/mhub/chain/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var minDepositAmount = sdk.NewInt(100)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper       Keeper
	bankKeeper   types.BankKeeper
	minterKeeper minterkeeper.Keeper
}

// Handle is the entry point for Attestation processing.
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.EthereumClaim) error {
	switch claim := claim.(type) {
	case *types.MsgDepositClaim:
		amount := a.keeper.oracleKeeper.ConvertFromEthValue(ctx, claim.TokenContract, claim.Amount)
		if amount.LT(minDepositAmount) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is too small to be deposited")
		}

		token := types.ERC20Token{
			Amount:   amount,
			Contract: claim.TokenContract,
		}
		coin := token.PeggyCoin(ctx, a.keeper.oracleKeeper)
		if _, err := types.ValidatePeggyCoin(coin, ctx, a.keeper.oracleKeeper); err != nil {
			return sdkerrors.Wrap(err, "invalid coin")
		}

		vouchers := sdk.Coins{coin}
		if err := a.bankKeeper.MintCoins(ctx, types.ModuleName, vouchers); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", vouchers)
		}

		addr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid reciever address")
		}

		// pay commissions
		{
			valset := a.minterKeeper.GetCurrentValset(ctx)
			commission := sdk.NewCoin(coin.Denom, coin.Amount.ToDec().Mul(a.keeper.oracleKeeper.GetCommissionForDemon(ctx, coin.Denom)).RoundInt()) // total commission
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
				_, err := a.minterKeeper.AddToOutgoingPool(ctx, sdk.AccAddress{}, val.MinterAddress, "#commission", sdk.NewCoin(commission.Denom, amount))
				if err != nil {
					return sdkerrors.Wrap(err, "commission withdrawal")
				}
			}

			a.minterKeeper.BuildOutgoingTXBatch(ctx, minterkeeper.OutgoingTxBatchSize)
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

		a.keeper.oracleKeeper.SetTxStatus(ctx, claim.TxHash, oracletypes.TX_STATUS_DEPOSIT_RECEIVED, "")

	case *types.MsgSendToMinterClaim:
		amount := a.keeper.oracleKeeper.ConvertFromEthValue(ctx, claim.TokenContract, claim.Amount)
		if amount.LT(minDepositAmount) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount is too small to be deposited")
		}

		receiver := sdk.AccAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		err := a.Handle(ctx, att, &types.MsgDepositClaim{
			EventNonce:     claim.EventNonce,
			TokenContract:  claim.TokenContract,
			Amount:         claim.Amount,
			EthereumSender: claim.EthereumSender,
			CosmosReceiver: receiver.String(),
			Orchestrator:   claim.Orchestrator,
		})
		if err != nil {
			return sdkerrors.Wrap(err, "deposit claim")
		}

		denom, err := a.keeper.OracleKeeper().GetCoins(ctx).GetDenomByEthereumAddress(claim.TokenContract)
		if err != nil {
			return sdkerrors.Wrap(err, "coin not found")
		}

		commission := sdk.NewCoin(denom, amount.ToDec().Mul(a.keeper.oracleKeeper.GetCommissionForDemon(ctx, denom)).RoundInt())
		_, err = a.minterKeeper.AddToOutgoingPool(ctx, receiver, claim.MinterReceiver, claim.TxHash, sdk.NewCoin(denom, amount).Sub(commission))
		if err != nil {
			return sdkerrors.Wrap(err, "withdraw")
		}
	case *types.MsgWithdrawClaim:
		a.keeper.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce, claim.TxSender, claim.TxHash)

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
