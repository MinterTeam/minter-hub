package helpers

import (
	pTypes "github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func IsValidHubRecipient(addr string) bool {
	if _, err := sdk.AccAddressFromBech32(addr); err == nil {
		return true
	}

	if err := pTypes.ValidateEthAddress(addr); err == nil {
		return true
	}

	return false
}
