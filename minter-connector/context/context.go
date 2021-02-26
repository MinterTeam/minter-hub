package context

import (
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"
)

type Context struct {
	LastCheckedMinterBlock uint64
	LastEventNonce         uint64
	LastBatchNonce         uint64
	LastValsetNonce        uint64

	MinterMultisigAddr string

	CosmosConn   *grpc.ClientConn
	MinterClient *http_client.Client

	OrcAddress   sdk.AccAddress
	OrcPriv      *secp256k1.PrivKey
	MinterWallet *wallet.Wallet
	Logger       log.Logger
}
