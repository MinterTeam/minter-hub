package client

import (
	"github.com/MinterTeam/mhub/chain/x/peggy/client/cli"
	"github.com/MinterTeam/mhub/chain/x/peggy/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitColdStorageTransferProposalTxCmd, rest.ProposalRESTHandler)
