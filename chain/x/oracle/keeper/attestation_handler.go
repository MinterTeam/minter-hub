package keeper

import (
	"github.com/MinterTeam/mhub/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"sort"
)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper        Keeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
}

// Handle is the entry point for Attestation processing.
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.Claim) error {
	switch claim := claim.(type) {
	case *types.MsgPriceClaim:
		votes := att.GetVotes()
		pricesSum := map[string][]sdk.Int{}
		totalPower := int64(0)

		for _, valaddr := range votes {
			validator, _ := sdk.ValAddressFromBech32(valaddr)
			power := a.stakingKeeper.GetLastValidatorPower(ctx, validator)
			totalPower += power

			priceClaim := a.keeper.GetClaim(ctx, sdk.AccAddress(validator).String(), claim.Epoch).(*types.GenericClaim).GetPriceClaim()
			prices := priceClaim.GetPrices()
			for _, item := range prices.List {
				for i := int64(0); i < power; i++ {
					pricesSum[item.Name] = append(pricesSum[item.Name], item.Value)
				}
			}
		}

		prices := types.Prices{}
		for name, price := range pricesSum {
			sort.Slice(price, func(i, j int) bool {
				return price[i].LT(price[j])
			})

			prices.List = append(prices.List, &types.Price{
				Name:  name,
				Value: price[len(price)/2],
			})
		}

		a.keeper.storePrices(ctx, &prices)

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
