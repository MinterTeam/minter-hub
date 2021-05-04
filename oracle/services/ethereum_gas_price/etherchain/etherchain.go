package etherchain

import (
	"encoding/json"
	"time"

	"github.com/MinterTeam/minter-hub-oracle/services/ethereum_gas_price"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/valyala/fasthttp"
)

type result struct {
	Fast int64 `json:"fast"`
}

type Service struct {
	logger log.Logger
}

func New(logger log.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) GetGasPrice() ethereum_gas_price.GasPrice {
	_, body, err := fasthttp.Get(nil, "https://www.etherchain.org/api/gasPriceOracle")

	if err != nil {
		s.logger.Error("Error getting eth gas price", "err", err.Error())
		time.Sleep(time.Second)
		return s.GetGasPrice()
	}

	var result result

	if err := json.Unmarshal(body, &result); err != nil {
		s.logger.Error("Error getting eth gas price", "err", err.Error())
		time.Sleep(time.Second)
		return s.GetGasPrice()
	}

	return ethereum_gas_price.GasPrice{
		Fast: sdk.NewInt(result.Fast),
	}
}
