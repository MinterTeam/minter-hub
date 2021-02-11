package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.NewLegacyAmino()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterInterfaces regiesteres the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPriceClaim{},
	)

	registry.RegisterInterface(
		"oracle.v1beta1.Claim",
		(*Claim)(nil),
		&MsgPriceClaim{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Claim)(nil), nil)
	cdc.RegisterConcrete(&MsgPriceClaim{}, "oracle/MsgPriceClaim", nil)
	cdc.RegisterConcrete(&Attestation{}, "oracle/Attestation", nil)
}
