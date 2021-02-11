package config

import (
	"flag"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
)

var cfg *Config

func Get() Config {
	if cfg == nil {
		cfg = &Config{}

		minterMsigAddr := flag.String("minter-multisig", "", "")
		minterChain := flag.String("minter-chain", "", "")
		minterMnemonic := flag.String("minter-mnemonic", "", "")
		minterNodeUrl := flag.String("minter-node-url", "", "")
		minterStartBlock := flag.Int("minter-start-block", 1, "")
		minterStartEventNonce := flag.Int("minter-start-event-nonce", 1, "")
		minterStartBatchNonce := flag.Int("minter-start-batch-nonce", 1, "")
		minterStartValsetNonce := flag.Int("minter-start-valset-nonce", 1, "")

		cosmosMnemonic := flag.String("cosmos-mnemonic", "", "")
		cosmosNodeUrl := flag.String("cosmos-node-url", "", "")
		tendermintNodeUrl := flag.String("tm-node-url", "", "")

		flag.Parse()

		var minterChainId transaction.ChainID
		switch *minterChain {
		case "mainnet":
			minterChainId = transaction.MainNetChainID
		case "testnet":
			minterChainId = transaction.TestNetChainID
		default:
			panic("unknown minter chain id")
		}

		cfg.Minter = MinterConfig{
			MultisigAddr:     *minterMsigAddr,
			ChainID:          minterChainId,
			Mnemonic:         *minterMnemonic,
			StartBlock:       uint64(*minterStartBlock),
			StartEventNonce:  uint64(*minterStartEventNonce),
			StartBatchNonce:  uint64(*minterStartBatchNonce),
			StartValsetNonce: uint64(*minterStartValsetNonce),
			NodeUrl:          *minterNodeUrl,
		}

		cfg.Cosmos = CosmosConfig{
			Mnemonic:    *cosmosMnemonic,
			NodeGrpcUrl: *cosmosNodeUrl,
			TmUrl:       *tendermintNodeUrl,
		}
	}

	return *cfg
}

type MinterConfig struct {
	MultisigAddr     string
	ChainID          transaction.ChainID
	Mnemonic         string
	StartBlock       uint64
	StartEventNonce  uint64
	StartBatchNonce  uint64
	StartValsetNonce uint64
	NodeUrl          string
}

type CosmosConfig struct {
	Mnemonic    string
	NodeGrpcUrl string
	TmUrl       string
}

type Config struct {
	Minter MinterConfig
	Cosmos CosmosConfig
}
