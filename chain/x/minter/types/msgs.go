package types

import (
	"encoding/hex"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	_ sdk.Msg = &MsgValsetConfirm{}
	_ sdk.Msg = &MsgValsetRequest{}
	_ sdk.Msg = &MsgSetMinterAddress{}
	_ sdk.Msg = &MsgSendToMinter{}
	_ sdk.Msg = &MsgRequestBatch{}
	_ sdk.Msg = &MsgConfirmBatch{}
)

// NewMsgValsetConfirm returns a new msgValsetConfirm
func NewMsgValsetConfirm(nonce uint64, minterAddress string, validator sdk.AccAddress, signature string) *MsgValsetConfirm {
	return &MsgValsetConfirm{
		Nonce:         nonce,
		Validator:     validator.String(),
		MinterAddress: minterAddress,
		Signature:     signature,
	}
}

// Route should return the name of the module
func (msg *MsgValsetConfirm) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgValsetConfirm) Type() string { return "valset_confirm" }

// ValidateBasic performs stateless checks
func (msg *MsgValsetConfirm) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if err := ValidateMinterAddress(msg.MinterAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgValsetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgValsetConfirm) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgValsetRequest returns a new msgValsetRequest
func NewMsgValsetRequest(requester sdk.AccAddress) *MsgValsetRequest {
	return &MsgValsetRequest{
		Requester: requester.String(),
	}
}

// Route should return the name of the module
func (msg MsgValsetRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgValsetRequest) Type() string { return "valset_request" }

// ValidateBasic performs stateless checks
func (msg MsgValsetRequest) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Requester); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Requester)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgValsetRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgValsetRequest) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSetMinterAddress return a new MsgSetMinterAddress
// TODO: figure out if we need sdk.ValAddress here
func NewMsgSetMinterAddress(address string, validator sdk.AccAddress, signature string) *MsgSetMinterAddress {
	return &MsgSetMinterAddress{
		Address:   address,
		Validator: validator.String(),
		Signature: signature,
	}
}

// Route should return the name of the module
func (msg MsgSetMinterAddress) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetMinterAddress) Type() string { return "set_minter_address" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid, and whether the Eth address has signed the validator address
// (proving control of the Eth address)
func (msg MsgSetMinterAddress) ValidateBasic() error {
	val, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if err := ValidateMinterAddress(msg.Address); err != nil {
		return sdkerrors.Wrap(err, "minter address")
	}
	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	err = ValidateMinterSignature(crypto.Keccak256(val.Bytes()), sigBytes, msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(err, "digest: %x\nsig: %x\naddress %s\nerror: %s\n", crypto.Keccak256(val.Bytes()), msg.Signature, msg.Address, err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetMinterAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetMinterAddress) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSendToMinter returns a new MsgSendToMinter
func NewMsgSendToMinter(sender sdk.AccAddress, destAddress string, send sdk.Coin) *MsgSendToMinter {
	return &MsgSendToMinter{
		Sender:     sender.String(),
		MinterDest: destAddress,
		Amount:     send,
	}
}

// Route should return the name of the module
func (msg MsgSendToMinter) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSendToMinter) Type() string { return "send_to_minter" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgSendToMinter) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if err := ValidateMinterAddress(msg.MinterDest); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO validate fee is sufficient, fixed fee to start
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSendToMinter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSendToMinter) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgRequestBatch returns a new msgRequestBatch
func NewMsgRequestBatch(requester sdk.AccAddress) *MsgRequestBatch {
	return &MsgRequestBatch{
		Requester: requester.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatch) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Requester); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Requester)
	}
	// TODO ensure that Demon matches hardcoded allowed value
	// TODO later make sure that Demon matches a list of tokens already
	// in the bridge to send
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgConfirmBatch) Type() string { return "confirm_batch" }

// ValidateBasic performs stateless checks
func (msg MsgConfirmBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if err := ValidateMinterAddress(msg.MinterSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	_, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MinterClaim represents a claim on ethereum state
type MinterClaim interface {
	GetEventNonce() uint64
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ MinterClaim = &MsgDepositClaim{}
	_ MinterClaim = &MsgWithdrawClaim{}
	_ MinterClaim = &MsgValsetClaim{}
	_ MinterClaim = &MsgSendToEthClaim{}
)

// GetType returns the type of the claim
func (msg *MsgValsetClaim) GetType() ClaimType {
	return CLAIM_TYPE_VALSET
}

// ValidateBasic performs stateless checks
func (msg *MsgValsetClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if msg.ValsetNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgValsetClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgValsetClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgValsetClaim) Type() string { return "valset_claim" }

// Route should return the name of the module
func (msg MsgValsetClaim) Route() string { return RouterKey }

const (
	TypeMsgValsetClaim = "valset_claim"
)

// Hash implements BridgeDeposit.Hash
func (msg *MsgValsetClaim) ClaimHash() []byte {
	path := fmt.Sprintf("valset/%d/%d", msg.EventNonce, msg.ValsetNonce)
	return tmhash.Sum([]byte(path))
}

// GetType returns the type of the claim
func (msg *MsgDepositClaim) GetType() ClaimType {
	return CLAIM_TYPE_DEPOSIT
}

// ValidateBasic performs stateless checks
func (msg *MsgDepositClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosReceiver)
	}
	if err := ValidateMinterAddress(msg.MinterSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if msg.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDepositClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDepositClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgDepositClaim) Type() string { return "deposit_claim" }

// Route should return the name of the module
func (msg MsgDepositClaim) Route() string { return RouterKey }

const (
	TypeMsgWithdrawClaim = "withdraw_claim"
)

// Hash implements BridgeDeposit.Hash
func (msg *MsgDepositClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%s/%s/", msg.CoinId, msg.MinterSender, msg.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}

// GetType returns the type of the claim
func (msg *MsgSendToEthClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_ETH
}

// ValidateBasic performs stateless checks
func (msg *MsgSendToEthClaim) ValidateBasic() error {
	if err := ValidateEthAddress(msg.EthReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.EthReceiver)
	}
	if err := ValidateMinterAddress(msg.MinterSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if msg.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSendToEthClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSendToEthClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgSendToEthClaim) Type() string { return "send_to_eth_claim" }

// Route should return the name of the module
func (msg MsgSendToEthClaim) Route() string { return RouterKey }

const (
	TypeMsgSendToEthClaim = "send_to_eth_claim"
)

// Hash implements BridgeDeposit.Hash
func (msg *MsgSendToEthClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%s/%s/", msg.CoinId, msg.MinterSender, msg.EthReceiver)
	return tmhash.Sum([]byte(path))
}

// GetType returns the claim type
func (msg *MsgWithdrawClaim) GetType() ClaimType {
	return CLAIM_TYPE_WITHDRAW
}

// ValidateBasic performs stateless checks
func (msg *MsgWithdrawClaim) ValidateBasic() error {
	if msg.EventNonce == 0 {
		return fmt.Errorf("event_nonce == 0")
	}
	if msg.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	return nil
}

// Hash implements WithdrawBatch.Hash
func (msg *MsgWithdrawClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d", msg.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (msg MsgWithdrawClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgWithdrawClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgWithdrawClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgWithdrawClaim) Type() string { return "withdraw_claim" }

const (
	TypeMsgDepositClaim = "deposit_claim"
)

func ValidateEthAddress(a string) error {
	if a == "" {
		return fmt.Errorf("empty")
	}
	if !regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(a) {
		return fmt.Errorf("address(%s) doesn't pass regex", a)
	}
	return nil
}
