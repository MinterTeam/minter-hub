package gasprice

import (
	"errors"
	"fmt"
	"time"

	"github.com/MinterTeam/minter-hub-oracle/config"
	"github.com/MinterTeam/minter-hub-oracle/services/ethereum/gasprice/providers"
	"github.com/MinterTeam/minter-hub-oracle/services/ethereum/gasprice/providers/etherchain"
	"github.com/MinterTeam/minter-hub-oracle/services/ethereum/gasprice/providers/ethgasstation"
	"github.com/tendermint/tendermint/libs/log"
)

type Service struct {
	providers []providers.Provider
	logger    log.Logger
}

func NewService(cfg *config.Config, logger log.Logger) (*Service, error) {
	if len(cfg.Ethereum.GasPriceProviders) == 0 {
		return nil, errors.New("define at least 1 ethereum gas price provider")
	}

	var pp []providers.Provider

	for _, name := range cfg.Ethereum.GasPriceProviders {
		var p providers.Provider

		switch name {
		case "ethgasstation":
			p = ethgasstation.New()
			break
		case "etherchain":
			p = etherchain.New()
			break
		default:
			return nil, errors.New(fmt.Sprintf(
				"unknown eth gas price provider: %s",
				name,
			))
		}

		pp = append(pp, p)
	}

	return &Service{
		providers: pp,
		logger:    logger,
	}, nil
}

func (s *Service) GetGasPrice() providers.GasPrice {
	for _, p := range s.providers {
		res, err := p.GetGasPrice()

		if err != nil {
			s.logger.Error(
				fmt.Sprintf("Error getting eth gas price from %s", p.Name()),
				"err",
				err.Error(),
			)

			continue
		}

		return *res
	}

	s.logger.Error("Failed to get eth gas price using all providers", "err")
	time.Sleep(time.Second)

	return s.GetGasPrice()
}
