package keeper

import (
	"fmt"
	oracletypes "github.com/MinterTeam/mhub/chain/x/oracle/types"
	"strconv"

	"github.com/MinterTeam/mhub/chain/x/minter/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const OutgoingTxBatchSize = 100

// BuildOutgoingTXBatch starts the following process chain:
// - find bridged denominator for given voucher type
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) BuildOutgoingTXBatch(ctx sdk.Context, maxElements int) (*types.OutgoingTxBatch, error) {
	if maxElements == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}
	// TODO: figure out how to check for know or unknown denoms? this might not matter anymore
	selectedTx, err := k.pickUnbatchedTX(ctx, maxElements)
	if len(selectedTx) == 0 || err != nil {
		return nil, err
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	minterNonce := k.autoIncrementID(ctx, types.MinterNonce)
	batch := &types.OutgoingTxBatch{
		BatchNonce:   nextID,
		MinterNonce:  minterNonce,
		Transactions: selectedTx,
	}
	k.storeBatch(ctx, batch)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(nextID)),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(nextID)),
	)

	for _, tx := range selectedTx {
		batchEvent = batchEvent.AppendAttributes(sdk.NewAttribute(types.AttributeKeyTxHash, tx.TxHash))
	}

	ctx.EventManager().EmitEvent(batchEvent)
	return batch, nil
}

// OutgoingTxBatchExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It frees all the transactions in the batch, then cancels all earlier batches
func (k Keeper) OutgoingTxBatchExecuted(ctx sdk.Context, nonce uint64, hash string) error {
	b := k.GetOutgoingTXBatch(ctx, nonce)
	if b == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}

	// cleanup outgoing TX pool
	for _, tx := range b.Transactions {
		k.removePoolEntry(ctx, tx.Id)
	}

	// Iterate through remaining batches
	k.IterateOutgoingTXBatches(ctx, func(key []byte, iter_batch *types.OutgoingTxBatch) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		// TODO: iterate only over batches we need to iterate over
		if iter_batch.BatchNonce < b.BatchNonce {
			k.CancelOutgoingTXBatch(ctx, iter_batch.BatchNonce)
		}

		return false
	})

	// Delete batch since it is finished
	k.deleteBatch(ctx, *b)

	batchEventExecuted := sdk.NewEvent(
		types.EventTypeOutgoingBatchExecuted,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(nonce)),
	)

	for _, tx := range b.Transactions {
		batchEventExecuted = batchEventExecuted.AppendAttributes(sdk.NewAttribute(types.AttributeKeyTxHash, tx.TxHash))
		k.oracleKeeper.SetTxStatus(ctx, tx.TxHash, oracletypes.TX_STATUS_BATCH_EXECUTED, hash)
	}

	ctx.EventManager().EmitEvent(batchEventExecuted)

	return nil
}

func (k Keeper) ValsetExecuted(ctx sdk.Context, nonce uint64) error {
	b := k.GetValsetRequest(ctx, nonce)
	if b == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}

	k.setLastValset(ctx, b)
	k.deleteValset(ctx, b)

	return nil
}

func (k Keeper) storeBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshalBinaryBare(batch))
}

func (k Keeper) deleteBatch(ctx sdk.Context, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxBatchKey(batch.BatchNonce))
}

// pickUnbatchedTX find TX in pool and remove from "available" second index
func (k Keeper) pickUnbatchedTX(ctx sdk.Context, maxElements int) ([]*types.OutgoingTransferTx, error) {
	var selectedTx []*types.OutgoingTransferTx
	var err error
	k.IterateOutgoingPool(ctx, func(txID uint64, tx *types.OutgoingTx) bool {
		var mCoin *types.MinterCoin
		mCoin, err = types.MinterCoinFromPeggyCoin(tx.Amount, ctx, k.oracleKeeper)
		txOut := &types.OutgoingTransferTx{
			Id:          txID,
			Sender:      tx.Sender,
			DestAddress: tx.DestAddr,
			MinterToken: types.NewMinterCoin(tx.Amount.Amount, mCoin.CoinId),
			TxHash:      tx.TxHash,
		}
		selectedTx = append(selectedTx, txOut)
		err = k.removeFromUnbatchedTXIndex(ctx, txID)
		return err != nil || len(selectedTx) == maxElements
	})
	return selectedTx, err
}

// GetOutgoingTXBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTXBatch(ctx sdk.Context, nonce uint64) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(nonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	var b types.OutgoingTxBatch
	k.cdc.MustUnmarshalBinaryBare(bz, &b)

	return &b
}

// CancelOutgoingTXBatch releases all TX in the batch and deletes the batch
func (k Keeper) CancelOutgoingTXBatch(ctx sdk.Context, nonce uint64) error {
	batch := k.GetOutgoingTXBatch(ctx, nonce)
	if batch == nil {
		return types.ErrUnknown
	}
	for _, tx := range batch.Transactions {
		k.prependToUnbatchedTXIndex(ctx, tx.Id)
	}

	// Delete batch since it is finished
	k.deleteBatch(ctx, *batch)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.GetBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.GetBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(nonce)),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(nonce)),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nil
}

// IterateOutgoingTXBatches iterates through all outgoing batches in DESC order.
func (k Keeper) IterateOutgoingTXBatches(ctx sdk.Context, cb func(key []byte, batch *types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTXBatchKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.OutgoingTxBatch
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(iter.Key(), &batch) {
			break
		}
	}
}
