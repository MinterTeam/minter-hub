package providers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Provider interface {
	Name() string
	GetGasPrice() (*GasPrice, error)
}

type GasPrice struct {
	Fast sdk.Int
}
