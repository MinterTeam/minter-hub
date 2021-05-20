package cosmos

import (
	"context"
	"strings"
	"time"

	"github.com/MinterTeam/mhub/chain/app"
	"github.com/MinterTeam/minter-hub-oracle/config"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	signing2 "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"
)

var encoding = app.MakeEncodingConfig()
var cfg *config.Config

func Setup(conf *config.Config) {
	cfg = conf

	cosmosCfg := sdk.GetConfig()
	cosmosCfg.SetBech32PrefixForAccount("hub", "hubpub")
	cosmosCfg.Seal()
}

func SendCosmosTx(msgs []sdk.Msg, address sdk.AccAddress, priv crypto.PrivKey, cosmosConn *grpc.ClientConn, logger log.Logger) {
	number, sequence := getAccount(address.String(), cosmosConn, logger)

	fee := sdk.NewCoins(sdk.NewCoin("hub", sdk.NewInt(1)))

	tx := encoding.TxConfig.NewTxBuilder()
	if err := tx.SetMsgs(msgs...); err != nil {
		panic(err)
	}

	tx.SetMemo("")
	tx.SetFeeAmount(fee)
	tx.SetGasLimit(500000)

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &sigData,
		Sequence: sequence,
	}

	if err := tx.SetSignatures(sig); err != nil {
		panic(err)
	}

	client, err := tmClient.New(cfg.Cosmos.RpcAddr, "")
	if err != nil {
		panic(err)
	}

	status, err := client.Status(context.TODO())
	if err != nil {
		panic(err)
	}

	signBytes, err := encoding.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signing2.SignerData{
		ChainID:       status.NodeInfo.Network,
		AccountNumber: number,
		Sequence:      sequence,
	}, tx.GetTx())
	if err != nil {
		panic(err)
	}

	// Sign those bytes
	sigBytes, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	// Construct the SignatureV2 struct
	sigData = signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &sigData,
		Sequence: sequence,
	}

	if err := tx.SetSignatures(sig); err != nil {
		panic(err)
	}

	txBytes, err := encoding.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		panic(err)
	}

	result, err := client.BroadcastTxCommit(context.Background(), txBytes)
	if err != nil {
		if !strings.Contains(err.Error(), "incorrect account sequence") {
			logger.Error("Error sending tx", "err", err.Error())
		}

		time.Sleep(1 * time.Second)
		SendCosmosTx(msgs, address, priv, cosmosConn, logger)
		return
	}

	cj, _ := result.DeliverTx.MarshalJSON()
	logger.Info("Sent tx", "result", string(cj))
}

func getAccount(address string, conn *grpc.ClientConn, logger log.Logger) (number, sequence uint64) {
	authClient := types.NewQueryClient(conn)

	response, err := authClient.Account(context.Background(), &types.QueryAccountRequest{Address: address})
	if err != nil {
		logger.Error("Error getting account", "err", err.Error())
		time.Sleep(1 * time.Second)
		return getAccount(address, conn, logger)
	}

	var account types.AccountI
	if err := encoding.Marshaler.UnpackAny(response.Account, &account); err != nil {
		logger.Error("Error getting account", "err", err.Error())
		time.Sleep(1 * time.Second)
		return getAccount(address, conn, logger)
	}

	return account.GetAccountNumber(), account.GetSequence()
}

func GetAccount(mnemonic string) (sdk.AccAddress, *secp256k1.PrivKey) {
	var orcPriv secp256k1.PrivKey
	seed := bip39.NewSeed(mnemonic, "")
	master, ch := hd.ComputeMastersFromSeed(seed)
	orcPriv.Key, _ = hd.DerivePrivateKeyForPath(master, ch, sdk.FullFundraiserPath)
	orcAddress := sdk.AccAddress(orcPriv.PubKey().Address())

	return orcAddress, &orcPriv
}
