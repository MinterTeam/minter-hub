package types

import (
	"crypto/ecdsa"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewMinterSignature creates a new signuature over a given byte array
func NewMinterSignature(hash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, sdkerrors.Wrap(ErrEmpty, "private key")
	}
	protectedHash := crypto.Keccak256Hash(hash)
	return crypto.Sign(protectedHash.Bytes(), privateKey)
}

// ValidateMinterSignature takes a message, an associated signature and public key and
// returns an error if the signature isn't valid
func ValidateMinterSignature(hash []byte, signature []byte, minterAddress string) error {
	if len(signature) < 65 {
		return sdkerrors.Wrap(ErrInvalid, "signature too short")
	}

	protectedHash := crypto.Keccak256Hash(hash)

	pubkey, err := crypto.SigToPub(protectedHash.Bytes(), signature)
	if err != nil {
		return sdkerrors.Wrap(err, "signature to public key")
	}

	addr := crypto.PubkeyToAddress(*pubkey)
	if "Mx"+addr.Hex()[2:] != minterAddress { // todo
		return sdkerrors.Wrap(ErrInvalid, "signature not matching")
	}

	return nil
}
