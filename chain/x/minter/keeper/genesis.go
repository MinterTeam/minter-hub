package keeper

import (
	"encoding/binary"
	"github.com/MinterTeam/mhub/chain/x/minter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.setParams(ctx, data.Params)
	ctx.KVStore(keeper.storeKey).Set(types.MinterNonce, sdk.Uint64ToBigEndian(data.StartMinterNonce))
}

func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	p := k.GetParams(ctx)

	bz := ctx.KVStore(k.storeKey).Get(types.MinterNonce)
	var startMinterNonce uint64 = 1
	if bz != nil {
		startMinterNonce = binary.BigEndian.Uint64(bz)
	}

	return types.GenesisState{
		Params:           &p,
		StartMinterNonce: startMinterNonce,
	}
}
