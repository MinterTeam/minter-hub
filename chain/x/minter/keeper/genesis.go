package keeper

import (
	"github.com/althea-net/peggy/module/x/minter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.setParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	p := k.GetParams(ctx)
	return types.GenesisState{
		Params: &p,
	}
}
