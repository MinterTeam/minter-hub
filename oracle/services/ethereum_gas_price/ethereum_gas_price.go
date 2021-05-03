package ethereum_gas_price

type Service interface {
	GetGasPrice() GasPrice
}

type GasPrice struct {
	Fast int64
}
