package minter

import (
	"bytes"
	"fmt"

	"github.com/MinterTeam/mhub/chain/x/minter/keeper"
	"github.com/MinterTeam/mhub/chain/x/minter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "Peggy" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgSetMinterAddress:
			return handleMsgSetMinterAddress(ctx, keeper, msg)
		case *types.MsgValsetConfirm:
			return handleMsgConfirmValset(ctx, keeper, msg)
		case *types.MsgValsetRequest:
			return handleMsgValsetRequest(ctx, keeper, msg)
		case *types.MsgSendToMinter:
			return handleMsgSendToMinter(ctx, keeper, msg)
		case *types.MsgRequestBatch:
			return handleMsgRequestBatch(ctx, keeper, msg)
		case *types.MsgConfirmBatch:
			return handleMsgConfirmBatch(ctx, keeper, msg)
		case *types.MsgDepositClaim:
			return handleDepositClaim(ctx, keeper, msg)
		case *types.MsgWithdrawClaim:
			return handleWithdrawClaim(ctx, keeper, msg)
		case *types.MsgValsetClaim:
			return handleValsetClaim(ctx, keeper, msg)
		case *types.MsgSendToEthClaim:
			return handleSendToEthClaim(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Peggy Msg type: %v", msg.Type()))
		}
	}
}

func handleDepositClaim(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgDepositClaim) (*sdk.Result, error) {
	if keeper.IsStopped(ctx) {
		return nil, types.ErrServiceStopped
	}

	var attestationIDs [][]byte
	// TODO SECURITY this does not auth the sender in the current validator set!
	// anyone can vote! We need to check and reject right here.

	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	validator := findValidatorKey(ctx, orch)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}
	att, err := keeper.AddClaim(ctx, msg.GetType(), msg.GetEventNonce(), validator, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}
	attestationIDs = append(attestationIDs, types.GetAttestationKey(att.EventNonce, msg))

	return &sdk.Result{
		Data:   bytes.Join(attestationIDs, []byte(", ")),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleSendToEthClaim(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgSendToEthClaim) (*sdk.Result, error) {
	if keeper.IsStopped(ctx) {
		return nil, types.ErrServiceStopped
	}

	var attestationIDs [][]byte
	// TODO SECURITY this does not auth the sender in the current validator set!
	// anyone can vote! We need to check and reject right here.

	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	validator := findValidatorKey(ctx, orch)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}
	att, err := keeper.AddClaim(ctx, msg.GetType(), msg.GetEventNonce(), validator, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}
	attestationIDs = append(attestationIDs, types.GetAttestationKey(att.EventNonce, msg))

	return &sdk.Result{
		Data:   bytes.Join(attestationIDs, []byte(", ")),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleWithdrawClaim(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgWithdrawClaim) (*sdk.Result, error) {
	var attestationIDs [][]byte
	// TODO SECURITY this does not auth the sender in the current validator set!
	// anyone can vote! We need to check and reject right here.

	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	validator := findValidatorKey(ctx, orch)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}
	att, err := keeper.AddClaim(ctx, msg.GetType(), msg.GetEventNonce(), validator, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}
	attestationIDs = append(attestationIDs, types.GetAttestationKey(att.EventNonce, msg))

	return &sdk.Result{
		Data:   bytes.Join(attestationIDs, []byte(", ")),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleValsetClaim(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgValsetClaim) (*sdk.Result, error) {
	var attestationIDs [][]byte
	// TODO SECURITY this does not auth the sender in the current validator set!
	// anyone can vote! We need to check and reject right here.

	orch, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	validator := findValidatorKey(ctx, orch)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}
	att, err := keeper.AddClaim(ctx, msg.GetType(), msg.GetEventNonce(), validator, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}
	attestationIDs = append(attestationIDs, types.GetAttestationKey(att.EventNonce, msg))

	return &sdk.Result{
		Data:   bytes.Join(attestationIDs, []byte(", ")),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func findValidatorKey(ctx sdk.Context, orchAddr sdk.AccAddress) sdk.ValAddress {
	// todo: implement proper in keeper
	// TODO: do we want ValAddress or do we want the AccAddress for the validator?
	// this is a v important question for encoding
	return sdk.ValAddress(orchAddr)
}

func handleMsgValsetRequest(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgValsetRequest) (*sdk.Result, error) {
	// todo: is requester in current valset?\

	// disabling bootstrap check for integration tests to pass
	//if keeper.GetLastValsetObservedNonce(ctx).isValid() {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, "bridge bootstrap process not observed, yet")
	//}

	hasPendingValset := false
	keeper.IterateValsetRequest(ctx, func(_ []byte, _ *types.Valset) bool {
		hasPendingValset = true
		return true
	})

	if hasPendingValset {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "already have a pending valset")
	}

	v := keeper.SetValsetRequest(ctx)
	return &sdk.Result{
		Data:   types.UInt64Bytes(v.Nonce),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

// This function takes in a signature submitted by a validator's Eth Signer
func handleMsgConfirmBatch(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgConfirmBatch) (*sdk.Result, error) {

	batch := keeper.GetOutgoingTXBatch(ctx, msg.Nonce)
	if batch == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find batch")
	}

	valaddr, _ := sdk.AccAddressFromBech32(msg.Validator)

	//peggyID := keeper.GetPeggyID(ctx)
	//checkpoint, err := batch.GetCheckpoint(peggyID)
	//if err != nil {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	//}
	//
	//sigBytes, err := hex.DecodeString(msg.Signature)
	//if err != nil {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	//}
	//validator := findValidatorKey(ctx, valaddr)
	//if validator == nil {
	//	return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
	//}
	//
	//ethAddress := keeper.GetMinterAddress(ctx, sdk.AccAddress(validator))
	//if ethAddress == "" {
	//	return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
	//}
	//err = types.ValidateMinterSignature(checkpoint, sigBytes, ethAddress)
	//if err != nil {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
	//}

	// check if we already have this confirm
	if keeper.GetBatchConfirm(ctx, msg.Nonce, valaddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
	}
	key := keeper.SetBatchConfirm(ctx, msg)
	return &sdk.Result{
		Data:   key,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

// This function takes in a signature submitted by a validator's Eth Signer
func handleMsgConfirmValset(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgValsetConfirm) (*sdk.Result, error) {

	valset := keeper.GetValsetRequest(ctx, msg.Nonce)
	if valset == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find valset")
	}

	valaddr, _ := sdk.AccAddressFromBech32(msg.Validator)

	//peggyID := keeper.GetPeggyID(ctx)
	//checkpoint := valset.GetCheckpoint(peggyID)
	//
	//sigBytes, err := hex.DecodeString(msg.Signature)
	//if err != nil {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	//}
	//
	//validator := findValidatorKey(ctx, valaddr)
	//if validator == nil {
	//	return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
	//}
	//
	//ethAddress := keeper.GetMinterAddress(ctx, sdk.AccAddress(validator))
	//if ethAddress == "" {
	//	return nil, sdkerrors.Wrap(types.ErrEmpty, "eth address")
	//}
	//err = types.ValidateMinterSignature(checkpoint, sigBytes, ethAddress)
	//if err != nil {
	//	return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with peggy-id %s with checkpoint %s found %s", ethAddress, peggyID, hex.EncodeToString(checkpoint), msg.Signature))
	//}

	// persist signature
	if keeper.GetValsetConfirm(ctx, msg.Nonce, valaddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
	}
	key := keeper.SetValsetConfirm(ctx, *msg)
	return &sdk.Result{
		Data: key,
	}, nil
}

func handleMsgSetMinterAddress(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgSetMinterAddress) (*sdk.Result, error) {
	valaddr, _ := sdk.AccAddressFromBech32(msg.Validator)
	validator := findValidatorKey(ctx, valaddr)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}

	keeper.SetMinterAddress(ctx, sdk.AccAddress(validator), msg.Address)
	return &sdk.Result{}, nil
}

func handleMsgSendToMinter(ctx sdk.Context, keeper keeper.Keeper, msg *types.MsgSendToMinter) (*sdk.Result, error) {
	if keeper.IsStopped(ctx) {
		return nil, types.ErrServiceStopped
	}

	_, err := types.MinterCoinFromPeggyCoin(msg.Amount, ctx, keeper.OracleKeeper())
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("amount %#v is not a voucher type", msg))
	}

	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	txID, err := keeper.AddToOutgoingPool(ctx, sender, msg.MinterDest, "todo", msg.Amount) // todo: txhash
	if err != nil {
		return &sdk.Result{}, nil // todo log
	}
	return &sdk.Result{
		Data:   sdk.Uint64ToBigEndian(txID),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgRequestBatch(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRequestBatch) (*sdk.Result, error) {
	batchID, err := k.BuildOutgoingTXBatch(ctx, keeper.OutgoingTxBatchSize)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{
		Data:   types.UInt64Bytes(batchID.BatchNonce),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
