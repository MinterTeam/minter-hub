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
		&MsgValsetConfirm{},
		&MsgValsetRequest{},
		&MsgSetMinterAddress{},
		&MsgSendToMinter{},
		&MsgRequestBatch{},
		&MsgConfirmBatch{},
		&MsgDepositClaim{},
		&MsgValsetClaim{},
		&MsgWithdrawClaim{},
		&MsgSendToEthClaim{},
	)

	registry.RegisterInterface(
		"minter.v1beta1.MinterClaim",
		(*MinterClaim)(nil),
		&MsgDepositClaim{},
		&MsgWithdrawClaim{},
		&MsgValsetClaim{},
		&MsgSendToEthClaim{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*MinterClaim)(nil), nil)

	cdc.RegisterConcrete(&MsgSetMinterAddress{}, "minter/MsgSetMinterAddress", nil)
	cdc.RegisterConcrete(&MsgValsetRequest{}, "minter/MsgValsetRequest", nil)
	cdc.RegisterConcrete(&MsgValsetConfirm{}, "minter/MsgValsetConfirm", nil)
	cdc.RegisterConcrete(&MsgSendToMinter{}, "minter/MsgSendToMinter", nil)
	cdc.RegisterConcrete(&MsgRequestBatch{}, "minter/MsgRequestBatch", nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, "minter/MsgConfirmBatch", nil)
	cdc.RegisterConcrete(&Valset{}, "minter/Valset", nil)
	cdc.RegisterConcrete(&MsgDepositClaim{}, "minter/MsgDepositClaim", nil)
	cdc.RegisterConcrete(&MsgSendToEthClaim{}, "minter/MsgSendToEthClaim", nil)
	cdc.RegisterConcrete(&MsgValsetClaim{}, "minter/MsgValsetClaim", nil)
	cdc.RegisterConcrete(&MsgWithdrawClaim{}, "minter/MsgWithdrawClaim", nil)
	cdc.RegisterConcrete(&OutgoingTxBatch{}, "minter/OutgoingTxBatch", nil)
	cdc.RegisterConcrete(&OutgoingTransferTx{}, "minter/OutgoingTransferTx", nil)
	cdc.RegisterConcrete(&MinterCoin{}, "minter/MinterCoin", nil)
	cdc.RegisterConcrete(&IDSet{}, "minter/IDSet", nil)
	cdc.RegisterConcrete(&Attestation{}, "minter/Attestation", nil)
}
