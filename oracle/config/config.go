package config

import (
	"flag"
)

var cfg *Config

func Get() Config {
	if cfg == nil {
		cfg = &Config{}

		minterNodeUrl := flag.String("minter-node-url", "", "")

		cosmosMnemonic := flag.String("cosmos-mnemonic", "", "")
		cosmosNodeUrl := flag.String("cosmos-node-url", "", "")
		tendermintNodeUrl := flag.String("tm-node-url", "", "")

		flag.Parse()

		cfg.Cosmos = CosmosConfig{
			Mnemonic:    *cosmosMnemonic,
			NodeGrpcUrl: *cosmosNodeUrl,
			TmUrl:       *tendermintNodeUrl,
		}

		cfg.Minter = MinterConfig{
			NodeUrl: *minterNodeUrl,
		}
	}

	return *cfg
}

type MinterConfig struct {
	NodeUrl string
}

type CosmosConfig struct {
	Mnemonic    string
	NodeGrpcUrl string
	TmUrl       string
}

type Config struct {
	Cosmos CosmosConfig
	Minter MinterConfig
}
