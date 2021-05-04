package ethereum_gas_price

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Service interface {
	GetGasPrice() GasPrice
}

type GasPrice struct {
	Fast sdk.Int
}
