package cli

import (
	"errors"
	"fmt"

	"github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string) *cobra.Command {
	peggyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the peggy module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	peggyQueryCmd.AddCommand([]*cobra.Command{
		CmdGetPrices(storeKey),
	}...)

	return peggyQueryCmd
}

func CmdGetPrices(storeKey string) *cobra.Command {
	return &cobra.Command{
		Use:   "prices",
		Short: "Query current prices",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/prices", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return errors.New("empty response")
			}

			var out types.Prices
			cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintProto(&out)
		},
	}
}
