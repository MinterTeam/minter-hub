package cosmos

import (
	"context"
	"github.com/MinterTeam/mhub/chain/app"
	"github.com/MinterTeam/mhub/chain/coins"
	mhub "github.com/MinterTeam/mhub/chain/x/minter/types"
	phub "github.com/MinterTeam/mhub/chain/x/peggy/types"
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
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"
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

func CreateClaims(orcAddress sdk.AccAddress, deposits []Deposit, batches []Batch, valsets []Valset) []sdk.Msg {
	var msgs []sdk.Msg
	for _, deposit := range deposits {
		amount, _ := sdk.NewIntFromString(deposit.Amount)
		fee, _ := sdk.NewIntFromString(deposit.Fee)

		denom, err := coins.GetDenomByMinterId(deposit.CoinID)
		if err != nil {
			panic(err)
		}

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
			}, &phub.MsgRequestBatch{
				Orchestrator: orcAddress.String(),
				Denom:        denom,
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

	return msgs
}

func SendCosmosTx(msgs []sdk.Msg, address sdk.AccAddress, priv crypto.PrivKey, cosmosConn *grpc.ClientConn) {
	number, sequence := getAccount(address.String(), cosmosConn)

	fee := sdk.NewCoins(sdk.NewCoin("hub", sdk.NewInt(1)))

	tx := encoding.TxConfig.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
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

	signBytes, err := encoding.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signing2.SignerData{
		ChainID:       "mhub-test",
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

	client, err := tmClient.New("http://"+cfg.Cosmos.TmUrl, "")
	if err != nil {
		println(err.Error())
		time.Sleep(1 * time.Second)
		SendCosmosTx(msgs, address, priv, cosmosConn)
		return
	}

	result, err := client.BroadcastTxCommit(context.Background(), txBytes)
	if err != nil {
		println(err.Error())
		time.Sleep(1 * time.Second)
		SendCosmosTx(msgs, address, priv, cosmosConn)
		return
	}

	cj, _ := result.CheckTx.MarshalJSON()
	println(string(cj))

	println(result.DeliverTx.GetCode(), result.DeliverTx.GetLog(), result.DeliverTx.GetInfo())

	if result.DeliverTx.GetCode() != 0 || result.CheckTx.GetCode() != 0 {
		time.Sleep(1 * time.Second)
		SendCosmosTx(msgs, address, priv, cosmosConn)
	}
}

func GetLastMinterNonce(address string, conn *grpc.ClientConn) uint64 {
	client := mhub.NewQueryClient(conn)

	result, err := client.LastEventNonceByAddr(context.Background(), &mhub.QueryLastEventNonceByAddrRequest{Address: address})
	if err != nil {
		panic(err)
	}

	return result.EventNonce
}

func getAccount(address string, conn *grpc.ClientConn) (number, sequence uint64) {
	authClient := types.NewQueryClient(conn)

	response, err := authClient.Account(context.Background(), &types.QueryAccountRequest{Address: address})
	if err != nil {
		println(err.Error())
		time.Sleep(1 * time.Second)
		return getAccount(address, conn)
	}

	var account types.AccountI
	if err := encoding.Marshaler.UnpackAny(response.Account, &account); err != nil {
		println(err.Error())
		time.Sleep(1 * time.Second)
		return getAccount(address, conn)
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
