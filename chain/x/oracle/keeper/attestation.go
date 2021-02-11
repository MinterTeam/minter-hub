package keeper

import (
	"fmt"
	"github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddClaim starts the following process chain:
// - Records that a given validator has made a claim about a given ethereum event, checking that the event nonce is contiguous
//   (non contiguous eventNonces indicate out of order events which can cause double spends)
// - Either creates a new attestation or adds the validator's vote to the existing attestation for this event
// - Checks if the attestation has enough votes to be considered "Observed", then attempts to apply it to the
//   consensus state (e.g. minting tokens for a deposit event)
// - If so, marks it "Observed" and emits an event
func (k Keeper) AddClaim(ctx sdk.Context, details types.Claim) (*types.Attestation, error) {
	if err := k.storeClaim(ctx, details); err != nil {
		return nil, sdkerrors.Wrap(err, "claim")
	}

	att := k.voteForAttestation(ctx, details)

	k.tryAttestation(ctx, att, details)

	k.SetAttestation(ctx, att, details)

	return att, nil
}

// storeClaim persists a claim. Fails when a claim submitted by an Eth signer does not increment the event nonce by exactly 1.
func (k Keeper) storeClaim(ctx sdk.Context, details types.Claim) error {
	// Store the claim
	genericClaim, _ := types.GenericClaimfromInterface(details)
	store := ctx.KVStore(k.storeKey)
	cKey := types.GetClaimKey(details)
	store.Set(cKey, k.cdc.MustMarshalBinaryBare(genericClaim))
	return nil
}

// voteForAttestation either gets the attestation for this claim from storage, or creates one if this is the first time a validator
// has submitted a claim for this exact event
func (k Keeper) voteForAttestation(
	ctx sdk.Context,
	details types.Claim,
) *types.Attestation {
	// Tries to get an attestation with the same eventNonce and details as the claim that was submitted.
	att := k.GetAttestation(ctx, details.GetEpoch(), details)

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			Epoch:    details.GetEpoch(),
			Observed: false,
		}
	}

	sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(details.GetClaimer()))

	// Add the validator's vote to this attestation
	att.Votes = append(att.Votes, sval.GetOperator().String())

	return att
}

// tryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) tryAttestation(ctx sdk.Context, att *types.Attestation, claim types.Claim) {
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.
	if !att.Observed && k.GetCurrentEpoch(ctx) > claim.GetEpoch() {
		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
		requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
		attestationPower := sdk.NewInt(0)
		for _, validator := range att.Votes {
			val, err := sdk.ValAddressFromBech32(validator)
			if err != nil {
				panic(err)
			}
			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
			// Add it to the attestation power's sum
			attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))
			// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
			// process the attestation, set Observed to true, and break
			if attestationPower.GTE(requiredPower) {
				k.processAttestation(ctx, att, claim)
				att.Observed = true
				k.emitObservedEvent(ctx, att, claim)
				break
			}
		}
	}
}

// emitObservedEvent emits an event with information about an attestation that has been applied to
// consensus state.
func (k Keeper) emitObservedEvent(ctx sdk.Context, att *types.Attestation, claim types.Claim) {
	observationEvent := sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyAttestationType, string(claim.GetType())),
		sdk.NewAttribute(types.AttributeKeyAttestationID, string(types.GetAttestationKey(att.Epoch, claim))), // todo: serialize with hex/ base64 ?
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(att.Epoch)),
		// TODO: do we want to emit more information?
	)
	ctx.EventManager().EmitEvent(observationEvent)
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation, claim types.Claim) {
	// then execute in a new Tx so that we can store state on failure
	// TODO: It seems that the validator who puts an attestation over the threshold of votes will also
	// be charged for the gas of applying it to the consensus state. We should figure out a way to avoid this.
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, *att, claim); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(att.Epoch, claim),
			"epoch", fmt.Sprint(att.Epoch),
		)
	} else {
		commit() // persist transient storage

		// TODO: after we commit, delete the outgoingtxbatch that this claim references
	}
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, att *types.Attestation, claim types.Claim) {
	store := ctx.KVStore(k.storeKey)
	att.ClaimHash = claim.ClaimHash()
	att.Height = uint64(ctx.BlockHeight())
	aKey := types.GetAttestationKey(att.Epoch, claim)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(att))
}

// SetAttestationUnsafe sets the attestation w/o setting height and claim hash
func (k Keeper) SetAttestationUnsafe(ctx sdk.Context, att *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKeyWithHash(att.Epoch, att.ClaimHash)
	store.Set(aKey, k.cdc.MustMarshalBinaryBare(att))
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, epoch uint64, details types.Claim) *types.Attestation {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(epoch, details)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.Attestation
	k.cdc.MustUnmarshalBinaryBare(bz, &att)
	return &att
}

// DeleteAttestation deletes an attestation given an event nonce and claim
func (k Keeper) DeleteAttestation(ctx sdk.Context, att types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAttestationKeyWithHash(att.Epoch, att.ClaimHash))
}

// GetAttestationMapping returns a mapping of eventnonce -> attestations at that nonce
func (k Keeper) GetAttestationMapping(ctx sdk.Context) (out map[uint64][]types.Attestation) {
	out = make(map[uint64][]types.Attestation)
	k.IterateAttestaions(ctx, func(_ []byte, att types.Attestation) bool {
		if val, ok := out[att.Epoch]; !ok {
			out[att.Epoch] = []types.Attestation{att}
		} else {
			out[att.Epoch] = append(val, att)
		}
		return false
	})
	return
}

// IterateAttestaions iterates through all attestations
func (k Keeper) IterateAttestaions(ctx sdk.Context, cb func([]byte, types.Attestation) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte(types.OracleAttestationKey)
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := types.Attestation{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &att)
		// cb returns true to stop early
		if cb(iter.Key(), att) {
			return
		}
	}
}

// HasClaim returns true if a claim exists
func (k Keeper) HasClaim(ctx sdk.Context, details types.Claim) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetClaimKey(details))
}

func (k Keeper) GetClaim(ctx sdk.Context, valaddr string, epoch uint64) types.Claim {
	store := ctx.KVStore(k.storeKey)
	details := &types.MsgPriceClaim{
		Epoch:        epoch,
		Orchestrator: valaddr,
	}

	bz := store.Get(types.GetClaimKey(details))
	if len(bz) == 0 {
		return nil
	}

	var claim types.GenericClaim
	k.cdc.MustUnmarshalBinaryBare(bz, &claim)

	return &claim
}

// IterateClaimsByValidatorAndType takes a validator key and a claim type and then iterates over these claims
func (k Keeper) IterateClaimsByValidatorAndType(ctx sdk.Context, claimType types.ClaimType, validatorKey sdk.ValAddress, cb func([]byte, types.Claim) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleClaimKey)
	prefix := []byte(validatorKey)
	iter := prefixStore.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		genericClaim := types.GenericClaim{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &genericClaim)
		// cb returns true to stop early
		if cb(iter.Key(), &genericClaim) {
			break
		}
	}
}

// GetClaimsByValidatorAndType returns the list of claims a validator has signed for
func (k Keeper) GetClaimsByValidatorAndType(ctx sdk.Context, claimType types.ClaimType, val sdk.ValAddress) (out []types.Claim) {
	k.IterateClaimsByValidatorAndType(ctx, claimType, val, func(_ []byte, claim types.Claim) bool {
		out = append(out, claim)
		return false
	})
	return
}
