package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisStateValidate(t *testing.T) {
	specs := map[string]struct {
		src    *GenesisState
		expErr bool
	}{
		"default params": {src: DefaultGenesisState(), expErr: false},
		"empty params":   {src: &GenesisState{Params: &Params{}}, expErr: false},
		"invalid params": {src: &GenesisState{
			Params: &Params{
				StartThreshold:                0,
				MinterAddress:                 "",
				BridgeChainId:                 3279089,
				SignedValsetsWindow:           0,
				SignedBatchesWindow:           0,
				SignedClaimsWindow:            0,
				SlashFractionValset:           sdk.Dec{},
				SlashFractionBatch:            sdk.Dec{},
				SlashFractionClaim:            sdk.Dec{},
				SlashFractionConflictingClaim: sdk.Dec{},
				Stopped:                       false,
			},
		}, expErr: true},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestStringToByteArray(t *testing.T) {
	specs := map[string]struct {
		testString string
		expErr     bool
	}{
		"16 bytes": {"lakjsdflaksdjfds", false},
		"32 bytes": {"lakjsdflaksdjfdslakjsdflaksdjfds", false},
		"33 bytes": {"€€€€€€€€€€€", true},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			_, err := strToFixByteArray(spec.testString)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
