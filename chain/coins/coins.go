package coins

import "errors"

type Coins []Coin

type Coin struct {
	Denom      string
	EthAddress string
	MinterID   uint64
}

func GetDenomByEthereumAddress(address string) (string, error) {
	list := GetCoins()
	for _, coin := range list {
		if coin.EthAddress == address {
			return coin.Denom, nil
		}
	}

	return "", errors.New("coin not found")
}

func GetDenomByMinterId(id uint64) (string, error) {
	list := GetCoins()
	for _, coin := range list {
		if coin.MinterID == id {
			return coin.Denom, nil
		}
	}

	return "", errors.New("coin not found")
}

func GetMinterIdByDenom(denom string) (uint64, error) {
	list := GetCoins()
	for _, coin := range list {
		if coin.Denom == denom {
			return coin.MinterID, nil
		}
	}

	return 0, errors.New("coin not found")
}

func GetCoins() Coins {
	return Coins{
		//{
		//	Denom:      "lashin",
		//	EthAddress: "0x98C4408691165a7D892C2D9b5A2D9b9c9ac6FF19",
		//
		//	MinterID:   10,
		//},
		//{
		//	Denom:      "chain",
		//	EthAddress: "0x7186b91eB6EaeE563bc670d475A9E8555b755A57",
		//	MinterID:   1761,
		//},
		{
			Denom:      "hub",
			EthAddress: "0x8C2B6949590bEBE6BC1124B670e58DA85b081b2E",
			MinterID:   1,
		},
	}
}
