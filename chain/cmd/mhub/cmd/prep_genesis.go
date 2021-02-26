package cmd

import (
	"encoding/json"
	"fmt"
	oracletypes "github.com/MinterTeam/mhub/chain/x/oracle/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func AddPrepGenesisCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prep-genesis [hub-coin-id]",
		Short: "Prepare genesis values",
		Long: ``,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// staking

			stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)
			stakingGenState.Params.BondDenom = "hub"
			stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[stakingtypes.ModuleName] = stakingGenStateBz

			// crisis

			var crisisGenState crisistypes.GenesisState
			cdc.MustUnmarshalJSON(appState[crisistypes.ModuleName], &crisisGenState)
			crisisGenState.ConstantFee.Denom = "hub"
			crisisGenStateBz, err := cdc.MarshalJSON(&crisisGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[crisistypes.ModuleName] = crisisGenStateBz

			// gov

			var govGenState govtypes.GenesisState
			cdc.MustUnmarshalJSON(appState[govtypes.ModuleName], &govGenState)
			govGenState.DepositParams.MinDeposit = sdk.Coins{sdk.NewInt64Coin("hub", 10000)}
			govGenState.VotingParams.VotingPeriod = time.Minute * 20
			govGenStateBz, err := cdc.MarshalJSON(&govGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[govtypes.ModuleName] = govGenStateBz


			// mint

			var mintGenState minttypes.GenesisState
			cdc.MustUnmarshalJSON(appState[minttypes.ModuleName], &mintGenState)
			mintGenState.Params.MintDenom = "hub"
			mintGenState.Params.InflationMax = sdk.NewDec(0)
			mintGenState.Params.InflationMin = sdk.NewDec(0)
			mintGenState.Params.InflationRateChange = sdk.NewDec(0)
			mintGenState.Minter.Inflation = sdk.NewDec(0)
			mintGenStateBz, err := cdc.MarshalJSON(&mintGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[minttypes.ModuleName] = mintGenStateBz

			// oracle

			var oracleGenState oracletypes.GenesisState
			cdc.MustUnmarshalJSON(appState[oracletypes.ModuleName], &oracleGenState)
			minterId, _ := strconv.Atoi(args[0])
			oracleGenState.Params.Coins[0].MinterId = uint64(minterId)
			oracleGenStateBz, err := cdc.MarshalJSON(&oracleGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[oracletypes.ModuleName] = oracleGenStateBz

			appStateJSON, err := json.MarshalIndent(appState, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}
