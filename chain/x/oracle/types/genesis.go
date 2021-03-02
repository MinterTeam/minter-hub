package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)

	// ParamsStoreKeySignedClaimsWindow stores the signed blocks window
	ParamsStoreKeySignedClaimsWindow = []byte("SignedClaimsWindow")

	// ParamsStoreSlashFractionClaim stores the slash fraction Claim
	ParamsStoreSlashFractionClaim = []byte("SlashFractionClaim")

	// ParamsStoreSlashFractionConflictingClaim stores the slash fraction ConflictingClaim
	ParamsStoreSlashFractionConflictingClaim = []byte("SlashFractionConflictingClaim")

	ParamsStoreCoins = []byte("Coins")

	ParamsMinBatchGas          = []byte("MinBatchGas")
	ParamsMinSingleWithdrawGas = []byte("MinSingleWithdrawGas")

	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (s GenesisState) ValidateBasic() error {
	if err := s.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	return nil
}

// DefaultGenesisState returns empty genesis state
// TODO: set some better defaults here
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// DefaultParams returns a copy of the default params
func DefaultParams() *Params {
	return &Params{
		SignedClaimsWindow:            10000,
		SlashFractionClaim:            sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionConflictingClaim: sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		Coins: []*Coin{
			{
				Denom:       "hub",
				EthAddr:     "0x8C2B6949590bEBE6BC1124B670e58DA85b081b2E",
				MinterId:    1,
				EthDecimals: 18,
			},
			{
				Denom:       "usdc",
				EthAddr:     "0x4d153722A1b75204c52CD8681eaED174b90fD1A8",
				MinterId:    207,
				EthDecimals: 6,
			},
		},
		MinBatchGas:          100000,
		MinSingleWithdrawGas: 50000,
	}
}

// ValidateBasic checks that the parameters have valid values.
func (p Params) ValidateBasic() error {
	if err := validateSignedClaimsWindow(p.SignedClaimsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSlashFractionClaim(p.SlashFractionClaim); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionConflictingClaim(p.SlashFractionConflictingClaim); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	return nil
}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeySignedClaimsWindow, &p.SignedClaimsWindow, validateSignedClaimsWindow),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionClaim, &p.SlashFractionClaim, validateSlashFractionClaim),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionConflictingClaim, &p.SlashFractionConflictingClaim, validateSlashFractionConflictingClaim),
		paramtypes.NewParamSetPair(ParamsStoreCoins, &p.Coins, validateCoins),
		paramtypes.NewParamSetPair(ParamsMinBatchGas, &p.MinBatchGas, validateMinBatchGas),
		paramtypes.NewParamSetPair(ParamsMinSingleWithdrawGas, &p.MinSingleWithdrawGas, validateMinSingleWithdrawGas),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func validateSignedClaimsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionClaim(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionConflictingClaim(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateCoins(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.([]*Coin); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, coin := range i.([]*Coin) {
		if coin.Denom == "" {
			return fmt.Errorf("empty denom is not allowed")
		}

		if coin.EthDecimals == 0 {
			return fmt.Errorf("incorrect eth decimals")
		}

		if coin.EthAddr == "" { // todo: check
			return fmt.Errorf("incorrect eth addr")
		}

		// todo: check duplicates
	}

	return nil
}

func validateMinSingleWithdrawGas(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if i.(uint64) < 1 {
		return fmt.Errorf("invalid parameter value: min single withdraw gas")
	}

	return nil
}

func validateMinBatchGas(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if i.(uint64) < 1 {
		return fmt.Errorf("invalid parameter value: min batch gas")
	}

	return nil
}

func strToFixByteArray(s string) ([32]byte, error) {
	var out [32]byte
	if len([]byte(s)) > 32 {
		return out, fmt.Errorf("string too long")
	}
	copy(out[:], s)
	return out, nil
}
