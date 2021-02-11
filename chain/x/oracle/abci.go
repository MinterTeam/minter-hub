package oracle

import (
	"github.com/althea-net/peggy/module/x/oracle/keeper"
	"github.com/althea-net/peggy/module/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

	// process claims
	if ctx.BlockHeight()%5 == 0 {
		k.ProcessCurrentEpoch(ctx)
	}

	params := k.GetParams(ctx)
	currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)

	// #3 condition
	attmap := k.GetAttestationMapping(ctx)
	for _, atts := range attmap {
		// slash conflicting votes
		if len(atts) > 1 {
			var unObs []types.Attestation
			oneObserved := false
			for _, att := range atts {
				if att.Observed {
					oneObserved = true
					continue
				}
				unObs = append(unObs, att)
			}
			// if one is observed delete the *other* attestations, do not delete
			// the original as we will need it later.
			if oneObserved {
				for _, att := range unObs {
					for _, valaddr := range att.Votes {
						validator, _ := sdk.ValAddressFromBech32(valaddr)
						val := k.StakingKeeper.Validator(ctx, validator)
						cons, _ := val.GetConsAddr()
						k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), k.StakingKeeper.GetLastValidatorPower(ctx, validator), params.SlashFractionConflictingClaim)
						k.StakingKeeper.Jail(ctx, cons)
					}
					k.DeleteAttestation(ctx, att)
				}
			}
		}

		if len(atts) == 1 {
			att := atts[0]
			windowPassed := uint64(ctx.BlockHeight()) > params.SignedClaimsWindow && uint64(ctx.BlockHeight())-params.SignedClaimsWindow > att.Height
			// if the signing window has passed and the attestation is still unobserved wait.
			if windowPassed && att.Observed {
				for _, bv := range currentBondedSet {
					found := false
					for _, val := range att.Votes {
						confVal, _ := sdk.ValAddressFromBech32(val)
						if confVal.Equals(bv.GetOperator()) {
							found = true
							break
						}
					}
					if !found {
						cons, _ := bv.GetConsAddr()
						k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), k.StakingKeeper.GetLastValidatorPower(ctx, bv.GetOperator()), params.SlashFractionClaim)
						k.StakingKeeper.Jail(ctx, cons)
					}
				}
				k.DeleteAttestation(ctx, att)
			}
		}
	}

	// #4 condition (stretch goal)
	// TODO: lost eth key or delegate key
	// 1. submit a message signed by the priv key to the chain and it slashes the validator who delegated to that key
	// return

	// TODO: prune outgoing tx batches while looping over them above, older than 15h and confirmed
	// TODO: prune claims, attestations
}
