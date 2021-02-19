package keeper

import (
	"fmt"
	"strconv"

	"github.com/MinterTeam/mhub/chain/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const OutgoingTxBatchSize = 100
const gweiInEth = 1e9

// BuildOutgoingTXBatch starts the following process chain:
// - find bridged denominator for given voucher type
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) BuildOutgoingTXBatch(ctx sdk.Context, contractAddress string, maxElements int) (*types.OutgoingTxBatch, error) {
	if maxElements == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}

	selectedTx := k.pickUnbatchedTX(ctx, contractAddress, maxElements)
	if len(selectedTx) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "empty batch")
	}

	totalCommission := sdk.NewInt(0)
	for _, tx := range selectedTx {
		totalCommission = totalCommission.Add(tx.Erc20Fee.Amount)
	}

	coinId, err := k.oracleKeeper.GetCoins(ctx).GetMinterIdByDenom(selectedTx[0].Erc20Fee.PeggyCoin(ctx, k.oracleKeeper).Denom)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "coin not found")
	}

	price, err := k.oracleKeeper.GetMinterPrice(ctx, coinId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "coin price")
	}

	gasPrice, err := k.oracleKeeper.GetEthGasPrice(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "gas price")
	}

	ethPrice, err := k.oracleKeeper.GetEthPrice(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "eth price")
	}

	totalUsdCommission := totalCommission.Mul(price).Quo(k.oracleKeeper.GetPipInBip())
	totalUsdGas := gasPrice.Mul(ethPrice).MulRaw(int64(k.oracleKeeper.GetMinBatchGas(ctx))).QuoRaw(gweiInEth).QuoRaw(k.oracleKeeper.GetGasUnits())
	if totalUsdCommission.LT(totalUsdGas) {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "not enough gas yet")
	}

	for _, tx := range selectedTx {
		if err := k.removeFromUnbatchedTXIndex(ctx, tx.Erc20Fee.PeggyCoin(ctx, k.oracleKeeper), tx.Id); err != nil {
			return nil, sdkerrors.Wrap(err, "fee")
		}
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	batch := &types.OutgoingTxBatch{
		BatchNonce:    nextID,
		Transactions:  selectedTx,
		TokenContract: contractAddress,
	}
	k.StoreBatch(ctx, batch)

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
func (k Keeper) OutgoingTxBatchExecuted(ctx sdk.Context, tokenContract string, nonce uint64, txSender string, txHash string) error {
	b := k.GetOutgoingTXBatch(ctx, tokenContract, nonce)
	if b == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}

	totalFee := sdk.NewInt64Coin(b.Transactions[0].Erc20Fee.PeggyCoin(ctx, k.oracleKeeper).Denom, 0)
	// cleanup outgoing TX pool
	for _, tx := range b.Transactions {
		totalFee = totalFee.Add(tx.Erc20Fee.PeggyCoin(ctx, k.oracleKeeper))
		k.removePoolEntry(ctx, tx.Id)

		k.SubLockedCoins(ctx, sdk.Coins{tx.Erc20Token.PeggyCoin(ctx, k.oracleKeeper)})
	}
	commissionKeeperAddress := sdk.AccAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	vouchers := sdk.Coins{totalFee}
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, vouchers); err != nil {
		return sdkerrors.Wrapf(err, "mint vouchers coins: %s", vouchers)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, commissionKeeperAddress, vouchers); err != nil {
		return sdkerrors.Wrap(err, "transfer vouchers")
	}

	k.minterKeeper.AddToOutgoingPool(ctx, commissionKeeperAddress, "Mx"+txSender[2:], txHash, totalFee)

	// Iterate through remaining batches
	k.IterateOutgoingTXBatches(ctx, func(key []byte, iterBatch *types.OutgoingTxBatch) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		// TODO: iterate only over batches we need to iterate over
		if iterBatch.TokenContract == b.TokenContract && iterBatch.BatchNonce < b.BatchNonce {
			k.CancelOutgoingTXBatch(ctx, tokenContract, iterBatch.BatchNonce)
		}
		return false
	})

	// Delete batch since it is finished
	k.DeleteBatch(ctx, *b)

	batchEventExecuted := sdk.NewEvent(
		types.EventTypeOutgoingBatchExecuted,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(nonce)),
		sdk.NewAttribute(types.AttributeKeyBatchTxHash, txHash),
	)

	for _, tx := range b.Transactions {
		batchEventExecuted = batchEventExecuted.AppendAttributes(sdk.NewAttribute(types.AttributeKeyTxHash, tx.TxHash))
	}

	ctx.EventManager().EmitEvent(batchEventExecuted)

	return nil
}

// StoreBatch stores a transaction batch
func (k Keeper) StoreBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	// set the current block height when storing the batch
	batch.Block = uint64(ctx.BlockHeight())
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshalBinaryBare(batch))
}

// StoreBatchUnsafe stores a transaction batch w/o setting the height
func (k Keeper) StoreBatchUnsafe(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshalBinaryBare(batch))
}

// DeleteBatch deletes an outgoing transaction batch
func (k Keeper) DeleteBatch(ctx sdk.Context, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce))
}

// pickUnbatchedTX find TX in pool and remove from "available" second index
func (k Keeper) pickUnbatchedTX(ctx sdk.Context, contractAddress string, maxElements int) []*types.OutgoingTransferTx {
	var selectedTx []*types.OutgoingTransferTx
	k.IterateOutgoingPoolByFee(ctx, contractAddress, func(txID uint64, tx *types.OutgoingTx) bool {
		txOut := &types.OutgoingTransferTx{
			Id:          txID,
			Sender:      tx.Sender,
			DestAddress: tx.DestAddr,
			Erc20Token:  types.NewERC20Token(tx.Amount.Amount, contractAddress),
			Erc20Fee:    types.NewERC20Token(tx.BridgeFee.Amount, contractAddress),
			TxHash:      tx.TxHash,
		}
		selectedTx = append(selectedTx, txOut)
		return len(selectedTx) == maxElements
	})
	return selectedTx
}

// GetOutgoingTXBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTXBatch(ctx sdk.Context, tokenContract string, nonce uint64) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(tokenContract, nonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	var b types.OutgoingTxBatch
	k.cdc.MustUnmarshalBinaryBare(bz, &b)
	// TODO: figure out why it drops the contract address in the ERC20 token representation
	for _, tx := range b.Transactions {
		tx.Erc20Token.Contract = tokenContract
		tx.Erc20Fee.Contract = tokenContract
	}
	return &b
}

// CancelOutgoingTXBatch releases all TX in the batch and deletes the batch
func (k Keeper) CancelOutgoingTXBatch(ctx sdk.Context, tokenContract string, nonce uint64) error {
	batch := k.GetOutgoingTXBatch(ctx, tokenContract, nonce)
	if batch == nil {
		return types.ErrUnknown
	}
	for _, tx := range batch.Transactions {
		tx.Erc20Fee.Contract = tokenContract
		k.prependToUnbatchedTXIndex(ctx, tx.Erc20Fee.PeggyCoin(ctx, k.oracleKeeper), tx.Id)
	}

	// Delete batch since it is finished
	k.DeleteBatch(ctx, *batch)

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

// GetOutgoingTxBatches returns the outgoing tx batches
func (k Keeper) GetOutgoingTxBatches(ctx sdk.Context) (out []*types.OutgoingTxBatch) {
	k.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		out = append(out, batch)
		return false
	})
	return
}
