package keeper

import (
	"context"

	"github.com/MinterTeam/mhub/chain/x/minter/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// ValsetConfirm
func (k msgServer) ValsetConfirm(c context.Context, msg *types.MsgValsetConfirm) (*types.MsgValsetConfirmResponse, error) {
	return nil, nil
}

// ValsetRequest
func (k msgServer) ValsetRequest(c context.Context, msg *types.MsgValsetRequest) (*types.MsgValsetRequestResponse, error) {
	return nil, nil
}

// SetEthAddress
func (k msgServer) SetMinterAddress(c context.Context, msg *types.MsgSetMinterAddress) (*types.MsgSetMinterAddressResponse, error) {
	return nil, nil
}

// SendToEth
func (k msgServer) SendToEth(c context.Context, msg *types.MsgSendToMinter) (*types.MsgSendToMinterResponse, error) {
	return nil, nil
}

// RequestBatch
func (k msgServer) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	return nil, nil
}

// ConfirmBatch
func (k msgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	return nil, nil
}

func (k msgServer) DepositClaim(c context.Context, msg *types.MsgDepositClaim) (*types.MsgDepositClaimResponse, error) {
	return nil, nil
}

func (k msgServer) WithdrawClaim(c context.Context, msg *types.MsgWithdrawClaim) (*types.MsgWithdrawClaimResponse, error) {
	return nil, nil
}

func (k msgServer) SendToEthClaim(c context.Context, msg *types.MsgSendToEthClaim) (*types.MsgSendToEthClaimResponse, error) {
	return nil, nil
}

func (k msgServer) ValsetClaim(c context.Context, msg *types.MsgValsetClaim) (*types.MsgValsetClaimResponse, error) {
	return nil, nil
}

func (k msgServer) SwapEthClaim(c context.Context, msg *types.MsgSwapEthClaim) (*types.MsgSwapEthClaimResponse, error) {
	return nil, nil
}
