package cli

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"log"

	"github.com/cosmos/cosmos-sdk/types/errors"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/MinterTeam/mhub/chain/x/minter/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func GetTxCmd(storeKey string) *cobra.Command {
	minterTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "minter",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	minterTxCmd.AddCommand([]*cobra.Command{
		CmdWithdrawToMinter(),
		CmdRequestBatch(),
		CmdUpdateMinterAddress(),
		CmdValsetRequest(),
		GetUnsafeTestingCmd(),
	}...)

	return minterTxCmd
}

func GetUnsafeTestingCmd() *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "unsafe_testing",
		Short:                      "helpers for testing. not going into production",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		CmdUnsafeMinterPrivKey(),
		CmdUnsafeMinterAddr(),
	}...)

	return testingTxCmd
}

// GetCmdUpdateEthAddress updates the network about the eth address that you have on record.
func CmdUpdateMinterAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-minter-addr [minter_private_key]",
		Short: "Update your Minter address which will be used for signing executables for the `multisig set`",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cosmosAddr := cliCtx.GetFromAddress()

			privKeyString := args[0][2:]

			// Make Eth Signature over validator address
			privateKey, err := ethCrypto.HexToECDSA(privKeyString)
			if err != nil {
				return err
			}

			hash := ethCrypto.Keccak256(cosmosAddr.Bytes())
			signature, err := types.NewMinterSignature(hash, privateKey)
			if err != nil {
				return sdkerrors.Wrap(err, "signing cosmos address with Minter key")
			}
			// You've got to do all this to get an Eth address from the private key
			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				return sdkerrors.Wrap(err, "casting public key to ECDSA")
			}
			minterAddress := ethCrypto.PubkeyToAddress(*publicKeyECDSA)

			msg := types.NewMsgSetMinterAddress("Mx"+minterAddress.String()[2:], cosmosAddr, hex.EncodeToString(signature))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdValsetRequest() *cobra.Command {
	return &cobra.Command{
		Use:   "valset-request",
		Short: "Trigger a new `multisig set` update request on the cosmos side",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			cosmosAddr := cliCtx.GetFromAddress()

			// Make the message
			msg := types.NewMsgValsetRequest(cosmosAddr)

			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
}

func CmdUnsafeMinterPrivKey() *cobra.Command {
	return &cobra.Command{
		Use:   "gen-minter-key",
		Short: "Generate and print a new ecdsa key",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := ethCrypto.GenerateKey()
			if err != nil {
				return errors.Wrap(err, "can not generate key")
			}
			k := "Mx" + hex.EncodeToString(ethCrypto.FromECDSA(key))
			println(k)
			return nil
		},
	}
}

func CmdUnsafeMinterAddr() *cobra.Command {
	return &cobra.Command{
		Use:   "minter-address",
		Short: "Print address for an ECDSA eth key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			privKeyString := args[0][2:]
			privateKey, err := ethCrypto.HexToECDSA(privKeyString)
			if err != nil {
				log.Fatal(err)
			}
			// You've got to do all this to get an Eth address from the private key
			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("error casting public key to ECDSA")
			}
			minterAddress := ethCrypto.PubkeyToAddress(*publicKeyECDSA).Hex()
			println(minterAddress)
			return nil
		},
	}
}

func CmdWithdrawToMinter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [from_key_or_cosmos_address] [to_minter_address] [amount]",
		Short: "Adds a new entry to the transaction pool to withdraw an amount from the Minter multisig",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			cosmosAddr := cliCtx.GetFromAddress()

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "amount")
			}

			// Make the message
			msg := types.MsgSendToMinter{
				Sender:     cosmosAddr.String(),
				MinterDest: args[1],
				Amount:     amount,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRequestBatch() *cobra.Command {
	return &cobra.Command{
		Use:   "build-batch [token_contract_address]",
		Short: "Build a new batch on the cosmos side for pooled withdrawal transactions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			cosmosAddr := cliCtx.GetFromAddress()

			msg := types.MsgRequestBatch{
				Requester: cosmosAddr.String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
}
