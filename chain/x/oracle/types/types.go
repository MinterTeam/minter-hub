package types

import (
	"encoding/binary"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// UInt64FromBytes create uint from binary big endian representation
func UInt64FromBytes(s []byte) uint64 {
	return binary.BigEndian.Uint64(s)
}

// UInt64Bytes uses the SDK byte marshaling to encode a uint64
func UInt64Bytes(n uint64) []byte {
	return sdk.Uint64ToBigEndian(n)
}

// UInt64FromString to parse out a uint64 for a nonce
func UInt64FromString(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

type Coins struct {
	list []*Coin
}

func NewCoins(list []*Coin) Coins {
	return Coins{list: list}
}

func (c Coins) List() []*Coin {
	return c.list
}

func (c Coins) GetByEthereumAddress(address string) (*Coin, error) {
	for _, coin := range c.list {
		if coin.EthAddr == address {
			return coin, nil
		}
	}

	return nil, errors.New("coin not found")
}

func (c Coins) GetDenomByEthereumAddress(address string) (string, error) {
	coin, err := c.GetByEthereumAddress(address)
	if err != nil {
		return "", errors.New("coin not found")
	}

	return coin.Denom, nil
}

func (c Coins) GetDenomByMinterId(id uint64) (string, error) {
	for _, coin := range c.list {
		if coin.MinterId == id {
			return coin.Denom, nil
		}
	}

	return "", errors.New("coin not found")
}

func (c Coins) GetMinterIdByDenom(denom string) (uint64, error) {
	for _, coin := range c.list {
		if coin.Denom == denom {
			return coin.MinterId, nil
		}
	}

	return 0, errors.New("coin not found")
}

func (c Coins) GetEthereumAddressByDenom(denom string) (string, error) {
	for _, coin := range c.list {
		if denom == coin.Denom {
			return coin.EthAddr, nil
		}
	}

	return "", errors.New("coin not found")
}
