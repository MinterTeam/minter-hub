package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	privateKey,_ := crypto.GenerateKey()
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Printf("Ethereum private key: %s\n",hexutil.Encode(privateKeyBytes)[2:])
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("Ethereum address: %s\n", address)

	println()

	minterWallet, _ := wallet.New()
	fmt.Printf("Minter mnemonic: %s\n", minterWallet.Mnemonic)
	fmt.Printf("Minter address: %s\n", minterWallet.Address)
}
