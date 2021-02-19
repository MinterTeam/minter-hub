package keeper

import (
	"fmt"
	"github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	"math"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	StakingKeeper types.StakingKeeper

	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace paramtypes.Subspace

	cdc        codec.BinaryMarshaler // The wire codec for binary encoding/decoding.
	bankKeeper types.BankKeeper

	AttestationHandler interface {
		Handle(sdk.Context, types.Attestation, types.Claim) error
	}
}

// NewKeeper returns a new instance of the peggy keeper
func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, paramSpace paramtypes.Subspace, stakingKeeper types.StakingKeeper, bankKeeper types.BankKeeper) Keeper {
	k := Keeper{
		cdc:           cdc,
		paramSpace:    paramSpace,
		storeKey:      storeKey,
		StakingKeeper: stakingKeeper,
		bankKeeper:    bankKeeper,
	}
	k.AttestationHandler = AttestationHandler{
		keeper:        k,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return k
}

func (k Keeper) GetPrices(ctx sdk.Context) *types.Prices {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.CurrentPricesKey)

	if len(bytes) == 0 {
		return nil
	}

	var prices types.Prices
	k.cdc.MustUnmarshalBinaryBare(bytes, &prices)

	return &prices
}

func (k Keeper) GetPipInBip() sdk.Int {
	t, _ := sdk.NewIntFromString("1000000000000000000")
	return t
}

func (k Keeper) GetMinterPrice(ctx sdk.Context, id uint64) (sdk.Int, error) {
	return k.getPrice(ctx, fmt.Sprintf("minter/%d", id))
}

func (k Keeper) GetEthGasPrice(ctx sdk.Context) (sdk.Int, error) {
	return k.getPrice(ctx, "eth/gas")
}

func (k Keeper) GetEthPrice(ctx sdk.Context) (sdk.Int, error) {
	return k.getPrice(ctx, "eth/0")
}

func (k Keeper) getPrice(ctx sdk.Context, key string) (sdk.Int, error) {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.CurrentPricesKey)

	if len(bytes) == 0 {
		return sdk.Int{}, sdkerrors.ErrKeyNotFound
	}

	var prices types.Prices
	k.cdc.MustUnmarshalBinaryBare(bytes, &prices)

	for _, price := range prices.GetList() {
		if price.GetName() == key {
			return price.Value, nil
		}
	}

	return sdk.Int{}, sdkerrors.ErrKeyNotFound
}

func (k Keeper) GetCurrentEpoch(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.CurrentEpochKey)

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

func (k Keeper) setCurrentEpoch(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CurrentEpochKey, types.UInt64Bytes(nonce))
}

/////////////////////////////
//       PARAMETERS        //
/////////////////////////////

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, ps types.Params) {
	k.paramSpace.SetParamSet(ctx, &ps)
}

func (k Keeper) GetCoins(ctx sdk.Context) types.Coins {
	p := k.GetParams(ctx)

	return types.NewCoins(p.Coins)
}

func (k Keeper) GetMinSingleWithdrawGas(ctx sdk.Context) uint64 {
	p := k.GetParams(ctx)

	return p.MinSingleWithdrawGas
}

func (k Keeper) GetMinBatchGas(ctx sdk.Context) uint64 {
	p := k.GetParams(ctx)

	return p.MinBatchGas
}

// logger returns a module-specific logger.
func (k Keeper) logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ProcessCurrentEpoch(ctx sdk.Context) {
	currentEpoch := k.GetCurrentEpoch(ctx)
	k.setCurrentEpoch(ctx, currentEpoch+1)

	claim := &types.MsgPriceClaim{
		Epoch: currentEpoch,
	}
	att := k.GetAttestation(ctx, currentEpoch, claim)
	if att != nil {
		k.tryAttestation(ctx, att, claim)
	}
}

func (k Keeper) storePrices(ctx sdk.Context, prices *types.Prices) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CurrentPricesKey, k.cdc.MustMarshalBinaryBare(prices))
}

func (k Keeper) GetGasUnits() int64 {
	return 10
}

func (k Keeper) GetNormalizedValPowers(ctx sdk.Context) map[string]uint64 {
	validators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	bridgeValidators := map[string]uint64{}
	var totalPower uint64

	for _, validator := range validators {
		validatorAddress := validator.GetOperator()

		p := uint64(k.StakingKeeper.GetLastValidatorPower(ctx, validatorAddress))
		totalPower += p

		bridgeValidators[validatorAddress.String()] = p
	}

	// normalize power values
	for address, power := range bridgeValidators {
		bridgeValidators[address] = sdk.NewUint(power).MulUint64(math.MaxUint16).QuoUint64(totalPower).Uint64()
	}

	return bridgeValidators
}

// prefixRange turns a prefix into a (start, end) range. The start is the given prefix value and
// the end is calculated by adding 1 bit to the start value. Nil is not allowed as prefix.
// 		Example: []byte{1, 3, 4} becomes []byte{1, 3, 5}
// 				 []byte{15, 42, 255, 255} becomes []byte{15, 43, 0, 0}
//
// In case of an overflow the end is set to nil.
//		Example: []byte{255, 255, 255, 255} becomes nil
// MARK finish-batches: this is where some crazy shit happens
func prefixRange(prefix []byte) ([]byte, []byte) {
	if prefix == nil {
		panic("nil key not allowed")
	}
	// special case: no prefix is whole range
	if len(prefix) == 0 {
		return nil, nil
	}

	// copy the prefix and update last byte
	end := make([]byte, len(prefix))
	copy(end, prefix)
	l := len(end) - 1
	end[l]++

	// wait, what if that overflowed?....
	for end[l] == 0 && l > 0 {
		l--
		end[l]++
	}

	// okay, funny guy, you gave us FFF, no end to this range...
	if l == 0 && end[0] == 0 {
		end = nil
	}
	return prefix, end
}
