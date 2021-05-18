package etherchain

import (
	"encoding/json"

	"github.com/MinterTeam/minter-hub-oracle/services/ethereum/gasprice/providers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/valyala/fasthttp"
)

type result struct {
	Fast int64 `json:"fast"`
}

type Provider struct {
}

func New() *Provider {
	return &Provider{}
}

func (p Provider) Name() string {
	return "etherchain"
}

func (p *Provider) GetGasPrice() (*providers.GasPrice, error) {
	_, body, err := fasthttp.Get(nil, "https://www.etherchain.org/api/gasPriceOracle")

	if err != nil {
		return nil, err
	}

	var result result

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &providers.GasPrice{
		Fast: sdk.NewInt(result.Fast).MulRaw(10),
	}, nil
}
