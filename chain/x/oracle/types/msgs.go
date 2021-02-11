package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// Claim represents a claim
type Claim interface {
	GetEpoch() uint64
	GetClaimer() sdk.AccAddress
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

func toClaimType(input int32) ClaimType {
	if input == 1 {
		return CLAIM_TYPE_PRICE
	} else {
		return CLAIM_TYPE_UNKNOWN
	}
}

func fromClaimType(input ClaimType) int32 {
	if input == CLAIM_TYPE_PRICE {
		return 1
	} else {
		return 0
	}
}

func (e *GenericClaim) GetType() ClaimType {
	return toClaimType(e.ClaimType)
}

func (e *GenericClaim) ClaimHash() []byte {
	return e.Hash
}

// by the time anything is turned into a generic
// claim it has already been validated
func (e *GenericClaim) ValidateBasic() error {
	return nil
}

func (e *GenericClaim) GetClaimer() sdk.AccAddress {
	val, _ := sdk.AccAddressFromBech32(e.EventClaimer)
	return val
}

func GenericClaimfromInterface(claim Claim) (*GenericClaim, error) {
	err := claim.ValidateBasic()
	if err != nil {
		return nil, err
	}
	gc := &GenericClaim{
		Epoch:     claim.GetEpoch(),
		ClaimType: fromClaimType(claim.GetType()),
		Hash:      claim.ClaimHash(),
	}

	switch claim := claim.(type) {
	case *MsgPriceClaim:
		gc.Claim = &GenericClaim_PriceClaim{
			PriceClaim: claim,
		}
	}

	return gc, nil
}

var (
	_ Claim = &MsgPriceClaim{}
	_ Claim = &GenericClaim{}
)

// GetType returns the type of the claim
func (e *MsgPriceClaim) GetType() ClaimType {
	return CLAIM_TYPE_PRICE
}

// ValidateBasic performs stateless checks
func (e *MsgPriceClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.Epoch == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgPriceClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgPriceClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic(fmt.Sprintf("MsgDepositClaim failed ValidateBasic! Should have been handled earlier %d %s", msg.Epoch, msg.Orchestrator))
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgPriceClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgPriceClaim) Type() string { return "price_claim" }

// Route should return the name of the module
func (msg MsgPriceClaim) Route() string { return RouterKey }

func (b *MsgPriceClaim) ClaimHash() []byte {
	return tmhash.Sum([]byte("claim"))
}
