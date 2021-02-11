package command

import (
	"errors"
	pTypes "github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeSendToEth = "send_to_eth"
const TypeSendToHub = "send_to_hub"

type Command struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
	Fee       string `json:"fee"`
}

func (cmd Command) Validate(amount sdk.Int) error {
	switch cmd.Type {
	case TypeSendToEth:
		if err := pTypes.ValidateEthAddress(cmd.Recipient); err != nil {
			return err
		}
	case TypeSendToHub:
		if _, err := sdk.AccAddressFromBech32(cmd.Recipient); err != nil {
			return err
		}
	default:
		return errors.New("wrong type")
	}

	fee, ok := sdk.NewIntFromString(cmd.Fee)
	if !ok {
		return errors.New("incorrect fee")
	}

	if amount.Sub(amount.QuoRaw(100)).LTE(fee) {
		return errors.New("incorrect fee")
	}

	return nil
}
