package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os"

	"github.com/althea-net/peggy/module/cmd/mhub/cmd"
	"github.com/cosmos/cosmos-sdk/server"
)

func main() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("hub", "hubpub")
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
