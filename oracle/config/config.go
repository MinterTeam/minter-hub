package config

import (
	"flag"

	"github.com/MinterTeam/minter-hub-oracle/services/ethereum_gas_price"
	"github.com/MinterTeam/minter-hub-oracle/services/ethereum_gas_price/etherchain"
	"github.com/MinterTeam/minter-hub-oracle/services/ethereum_gas_price/ethgasstation"
	"github.com/tendermint/tendermint/libs/log"
)

func Get(logger log.Logger) *Config {
	cfg := &Config{}

	minterNodeUrl := flag.String("minter-node-url", "", "")

	cosmosMnemonic := flag.String("cosmos-mnemonic", "", "")
	cosmosNodeUrl := flag.String("cosmos-node-url", "", "")
	tendermintNodeUrl := flag.String("tm-node-url", "", "")
	ethGasPriceProviderName := flag.String("eth-gas-price-provider", "ethgasstation", "")

	flag.Parse()

	cfg.Cosmos = CosmosConfig{
		Mnemonic:    *cosmosMnemonic,
		NodeGrpcUrl: *cosmosNodeUrl,
		TmUrl:       *tendermintNodeUrl,
	}

	cfg.Minter = MinterConfig{
		NodeUrl: *minterNodeUrl,
	}

	var ethGasPriceProvider ethereum_gas_price.Service

	switch *ethGasPriceProviderName {
	case "ethgasstation":
		ethGasPriceProvider = ethgasstation.New(logger)
		break
	case "etherchain":
		ethGasPriceProvider = etherchain.New(logger)
		break
	}

	cfg.EthGasPriceProvider = EthGasPriceProvider{
		Service: ethGasPriceProvider,
	}

	return cfg
}

type MinterConfig struct {
	NodeUrl string
}

type CosmosConfig struct {
	Mnemonic    string
	NodeGrpcUrl string
	TmUrl       string
}

type EthGasPriceProvider struct {
	Service ethereum_gas_price.Service
}

type Config struct {
	Cosmos              CosmosConfig
	Minter              MinterConfig
	EthGasPriceProvider EthGasPriceProvider
}
