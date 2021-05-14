package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeColdStorageTransfer defines the type for a ColdStorageTransferProposal
	ProposalTypeColdStorageTransfer = "MinterColdStorageTransfer"
)

// Assert ColdStorageTransferProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &ColdStorageTransferProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeColdStorageTransfer)
	govtypes.RegisterProposalTypeCodec(&ColdStorageTransferProposal{}, "minter/ColdStorageTransferProposal")
}

// NewCommunityPoolSpendProposal creates a new community pool spned proposal.
//nolint:interfacer
func NewCommunityPoolSpendProposal(amount sdk.Coins) *ColdStorageTransferProposal {
	return &ColdStorageTransferProposal{amount}
}

// GetTitle returns the title of a community pool spend proposal.
func (csp *ColdStorageTransferProposal) GetTitle() string { return "ColdStorageTransferProposal" }

// GetDescription returns the description of a community pool spend proposal.
func (csp *ColdStorageTransferProposal) GetDescription() string { return "ColdStorageTransferProposal" }

// GetDescription returns the routing key of a community pool spend proposal.
func (csp *ColdStorageTransferProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (csp *ColdStorageTransferProposal) ProposalType() string { return ProposalTypeColdStorageTransfer }

// ValidateBasic runs basic stateless validity checks
func (csp *ColdStorageTransferProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(csp)
	if err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (csp ColdStorageTransferProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Cold Storage Transfer Proposal:
  Amount:      %s`, csp.Amount))
	return b.String()
}
