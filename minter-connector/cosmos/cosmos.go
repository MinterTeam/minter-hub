package cosmos

import (
	"context"
	"github.com/MinterTeam/mhub/chain/app"
	mhub "github.com/MinterTeam/mhub/chain/x/minter/types"
	"github.com/MinterTeam/minter-hub-connector/command"
	"github.com/MinterTeam/minter-hub-connector/config"
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
	tmTypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"sort"
	"strings"
	"time"
)

var encoding = app.MakeEncodingConfig()
var cfg = config.Get()

type Batch struct {
	BatchNonce uint64
	EventNonce uint64
	TxHash     string
}

type Valset struct {
	ValsetNonce uint64
	EventNonce  uint64
}

type Deposit struct {
	Recipient  string
	Amount     string
	Fee        string
	EventNonce uint64
	Sender     string
	CoinID     uint64
	Type       string
	TxHash     string
}

func Setup() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("hub", "hubpub")
	config.Seal()
}

func CreateClaims(orcAddress sdk.AccAddress, deposits []Deposit, batches []Batch, valsets []Valset, logger log.Logger) []sdk.Msg {
	var msgs []sdk.Msg
	for _, deposit := range deposits {
		amount, _ := sdk.NewIntFromString(deposit.Amount)
		fee, _ := sdk.NewIntFromString(deposit.Fee)

		if deposit.Type == command.TypeSendToEth {
			msgs = append(msgs, &mhub.MsgSendToEthClaim{
				EventNonce:   deposit.EventNonce,
				CoinId:       deposit.CoinID,
				Amount:       amount,
				Fee:          fee,
				MinterSender: deposit.Sender,
				EthReceiver:  deposit.Recipient,
				Orchestrator: orcAddress.String(),
				TxHash:       deposit.TxHash,
			})
		} else {
			msgs = append(msgs, &mhub.MsgDepositClaim{
				EventNonce:     deposit.EventNonce,
				CoinId:         deposit.CoinID,
				Amount:         amount,
				MinterSender:   deposit.Sender,
				CosmosReceiver: deposit.Recipient,
				Orchestrator:   orcAddress.String(),
				TxHash:         deposit.TxHash,
			})
		}
	}

	for _, batch := range batches {
		msgs = append(msgs, &mhub.MsgWithdrawClaim{
			EventNonce:   batch.EventNonce,
			BatchNonce:   batch.BatchNonce,
			Orchestrator: orcAddress.String(),
			TxHash:       batch.TxHash,
		})
	}

	for _, valset := range valsets {
		msgs = append(msgs, &mhub.MsgValsetClaim{
			EventNonce:   valset.EventNonce,
			ValsetNonce:  valset.ValsetNonce,
			Orchestrator: orcAddress.String(),
		})
	}

	sort.Slice(msgs, func(i, j int) bool {
		return getEventNonceFromMsg(msgs[i]) < getEventNonceFromMsg(msgs[j])
	})

	return msgs
}

func getEventNonceFromMsg(msg sdk.Msg) uint64 {
	switch m := msg.(type) {
	case *mhub.MsgValsetClaim:
		return m.EventNonce
	case *mhub.MsgDepositClaim:
		return m.EventNonce
	case *mhub.MsgSendToEthClaim:
		return m.EventNonce
	case *mhub.MsgWithdrawClaim:
		return m.EventNonce
	}

	return 999999999
}

func SendCosmosTx(msgs []sdk.Msg, address sdk.AccAddress, priv crypto.PrivKey, cosmosConn *grpc.ClientConn, logger log.Logger) {
	if len(msgs) > 10 {
		SendCosmosTx(msgs[:10], address, priv, cosmosConn, logger)
		SendCosmosTx(msgs[10:], address, priv, cosmosConn, logger)
		return
	}

	number, sequence := getAccount(address.String(), cosmosConn, logger)

	fee := sdk.NewCoins(sdk.NewCoin("hub", sdk.NewInt(1)))

	tx := encoding.TxConfig.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}

	tx.SetMemo("")
	tx.SetFeeAmount(fee)
	tx.SetGasLimit(100000000)

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
			logger.Error("Failed broadcast tx", "err", err.Error())
		}

		time.Sleep(5 * time.Second)
		txResponse, err := client.Tx(context.Background(), tmTypes.Tx(txBytes).Hash(), false)
		if err != nil || txResponse.TxResult.IsErr() {
			SendCosmosTx(msgs, address, priv, cosmosConn, logger)
		}

		return
	}

	if result.DeliverTx.GetCode() != 0 || result.CheckTx.GetCode() != 0 {
		logger.Error("Error on sending cosmos tx with", "code", result.CheckTx.GetCode(), "log", result.DeliverTx.GetLog())
		time.Sleep(1 * time.Second)
		SendCosmosTx(msgs, address, priv, cosmosConn, logger)
	}

	logger.Info("Sending cosmos tx", "code", result.DeliverTx.GetCode(), "log", result.DeliverTx.GetLog(), "info", result.DeliverTx.GetInfo())
}

func GetLastMinterNonce(address string, conn *grpc.ClientConn) uint64 {
	client := mhub.NewQueryClient(conn)

	result, err := client.LastEventNonceByAddr(context.Background(), &mhub.QueryLastEventNonceByAddrRequest{Address: address})
	if err != nil {
		panic(err)
	}

	return result.EventNonce
}

func getAccount(address string, conn *grpc.ClientConn, logger log.Logger) (number, sequence uint64) {
	authClient := types.NewQueryClient(conn)

	response, err := authClient.Account(context.Background(), &types.QueryAccountRequest{Address: address})
	if err != nil {
		logger.Error("Error getting cosmos account", "err", err.Error())
		time.Sleep(1 * time.Second)
		return getAccount(address, conn, logger)
	}

	var account types.AccountI
	if err := encoding.Marshaler.UnpackAny(response.Account, &account); err != nil {
		logger.Error("Error unpacking cosmos account", "err", err.Error())
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
