package keeper

import (
	"github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryCurrentEpoch = "currentEpoch"
	QueryPrices       = "prices"
	QueryEthFee       = "eth_fee"
	QueryCoins       = "coins"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryCurrentEpoch:
			return queryCurrentEpoch(ctx, keeper)
		case QueryPrices:
			return queryPrices(ctx, keeper)
		case QueryEthFee:
			return queryEthFee(ctx, keeper)
		case QueryCoins:
			return queryCoins(ctx, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint", types.ModuleName)
		}
	}
}

func queryCurrentEpoch(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	currentEpoch := types.Epoch{
		Nonce: keeper.GetCurrentEpoch(ctx),
		Votes: nil,
	}

	att := keeper.GetAttestation(ctx, currentEpoch.Nonce, &types.MsgPriceClaim{})
	votes := att.GetVotes()
	for _, valaddr := range votes {
		priceClaim := keeper.GetClaim(ctx, valaddr, currentEpoch.Nonce).(*types.GenericClaim).GetPriceClaim()
		currentEpoch.Votes = append(currentEpoch.Votes, &types.Vote{
			Oracle: valaddr,
			Claim:  priceClaim,
		})
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, currentEpoch)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryPrices(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, keeper.GetPrices(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryCoins(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, keeper.GetCoins(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryEthFee(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	gasPrice, err := keeper.GetEthGasPrice(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "gas price")
	}

	ethPrice, err := keeper.GetEthPrice(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "eth price")
	}

	response := types.QueryEthFeeResponse{
		Min:  gasPrice.Mul(ethPrice).MulRaw(EthMaxExecutionGas).QuoRaw(gweiInEth).QuoRaw(keeper.GetGasUnits()),
		Fast: gasPrice.Mul(ethPrice).MulRaw(EthMaxFastExecutionGas).QuoRaw(gweiInEth).QuoRaw(keeper.GetGasUnits()),
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, &response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
