package cli

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/MinterTeam/mhub/chain/x/peggy/types"
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
		CmdGetCurrentValset(storeKey),
		CmdGetValsetRequest(storeKey),
		CmdGetValsetHistory(storeKey),
		CmdGetAttestationHistory(storeKey),
		CmdGetValsetConfirm(storeKey),
		CmdGetPendingValsetRequest(storeKey),
		CmdGetPendingOutgoingTXBatchRequest(storeKey),
		// CmdGetAllOutgoingTXBatchRequest(storeKey),
		// CmdGetOutgoingTXBatchByNonceRequest(storeKey),
		// CmdGetAllAttestationsRequest(storeKey),
		// CmdGetAttestationRequest(storeKey),
		QueryObserved(storeKey),
		QueryApproved(storeKey),
	}...)

	return peggyQueryCmd
}

func QueryObserved(storeKey string) *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "observed",
		Short:                      "observed ETH events",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		// CmdGetLastObservedNonceRequest(storeKey, cdc),
		// CmdGetLastObservedNoncesRequest(storeKey, cdc),
		// CmdGetLastObservedMultiSigUpdateRequest(storeKey, cdc),
		// CmdGetAllBridgedDenominatorsRequest(storeKey, cdc),
	}...)

	return testingTxCmd
}
func QueryApproved(storeKey string) *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "approved",
		Short:                      "approved cosmos operation",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		// CmdGetLastApprovedNoncesRequest(storeKey, cdc),
		// CmdGetLastApprovedMultiSigUpdateRequest(storeKey, cdc),
		// CmdGetInflightBatchesRequest(storeKey, cdc),
	}...)

	return testingTxCmd
}

func CmdGetCurrentValset(storeKey string) *cobra.Command {
	return &cobra.Command{
		Use:   "current-valset",
		Short: "Query current valset",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/currentValset", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return errors.New("empty response")
			}

			var out types.Valset
			cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintProto(&out)
		},
	}
}

func CmdGetValsetRequest(storeKey string) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-request [nonce]",
		Short: "Get requested valset with a particular nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			nonce := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetRequest/%s", storeKey, nonce), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("no valset request found for nonce %s", nonce)
			}

			var out types.Valset
			cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintProto(&out)
		},
	}
}

func CmdGetValsetConfirm(storeKey string) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-confirm [nonce] [bech32 validator address]",
		Short: "Get valset confirmation with a particular nonce from a particular validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetConfirm/%s/%s", storeKey, args[0], args[1]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("no valset confirmation found for nonce %s and address %s", args[0], args[1])
			}

			var out types.MsgValsetConfirm
			cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintProto(&out)
		},
	}
}

func CmdGetPendingValsetRequest(storeKey string) *cobra.Command {
	return &cobra.Command{
		Use:   "pending-valset-request [bech32 validator address]",
		Short: "Get the latest valset request which has not been signed by a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastPendingValsetRequest/%s", storeKey, args[0]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}

			var out types.Valset
			cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintProto(&out)
		},
	}
}

func CmdGetPendingOutgoingTXBatchRequest(storeKey string) *cobra.Command {
	return &cobra.Command{
		Use:   "pending-batch-request [bech32 validator address]",
		Short: "Get the latest outgoing TX batch request which has not been signed by a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastPendingBatchRequest/%s", storeKey, args[0]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}

			var out types.OutgoingTxBatch
			cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintProto(&out)
		},
	}
}

func CmdGetValsetHistory(storeKey string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "valset-history",
		Short: "Get valset history",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetHistory", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}

			var history types.ValsetHistory
			cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &history)

			for i, item := range history.GetValsets() {
				fmt.Printf("Valset(nonce: %d, members: %d)\n", item.GetValset().GetNonce(), len(item.GetValset().GetMembers()))
				if len(history.GetValsets()) < i+2 {
					continue
				}

				next := history.GetValsets()[i+1]
				for _, val := range next.GetValset().GetMembers() {
					confirmed := false
					for _, confirm := range item.GetConfirms() {
						if confirm.GetEthAddress() == val.GetEthereumAddress() {
							confirmed = true
							break
						}
					}

					if confirmed {
						fmt.Printf("  %s - ok\n", val.GetEthereumAddress())
					} else {
						fmt.Printf("  %s - NOT CONFIRMED\n", val.GetEthereumAddress())
					}
				}

				fmt.Printf("\n")
			}

			return nil
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetAttestationHistory(storeKey string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attestation-history",
		Short: "Get attestation history",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/attestationHistory", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}

			var history types.AttestationHistory
			cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &history)

			for i, item := range history.GetAttestations() {
				fmt.Printf("Attestation(nonce: %d)\n", item.GetEventNonce())

				for j, addr := range history.GetSigners()[i].GetVals() {
					if history.GetSigners()[i].GetSigned()[j] {
						fmt.Printf("  %s - ok\n", addr)
					} else {
						fmt.Printf("  %s - NOT CONFIRMED\n", addr)
					}
				}

				fmt.Printf("\n")
			}

			return nil
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
