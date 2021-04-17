package peggy

import (
	minterkeeper "github.com/MinterTeam/mhub/chain/x/minter/keeper"
	"github.com/MinterTeam/mhub/chain/x/peggy/keeper"
	"github.com/MinterTeam/mhub/chain/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// return coins from pools if they are too old
	hasRefunds := false
	coinsList := k.OracleKeeper().GetCoins(ctx).List()
	for _, coin := range coinsList {
		k.IterateOutgoingPoolByFee(ctx, coin.EthAddr, func(id uint64, tx *types.OutgoingTx) bool {
			if ctx.BlockTime().After(time.Unix(tx.ExpirationTime, 0)) {
				k.RefundOutgoingTx(ctx, id, tx)
				hasRefunds = true
			}

			return false
		})
	}

	if hasRefunds {
		k.MinterKeeper().BuildOutgoingTXBatch(ctx, minterkeeper.OutgoingTxBatchSize)
	}

	// valsets are sorted so the most recent one is first
	valsets := k.GetValsets(ctx)
	if len(valsets) == 0 {
		k.SetValsetRequest(ctx)
	}
}
