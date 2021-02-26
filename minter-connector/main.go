package main

import (
	c "context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/MinterTeam/mhub/chain/x/minter/types"
	oracleTypes "github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client/models"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
	"github.com/MinterTeam/minter-hub-connector/command"
	"github.com/MinterTeam/minter-hub-connector/config"
	"github.com/MinterTeam/minter-hub-connector/context"
	"github.com/MinterTeam/minter-hub-connector/cosmos"
	"github.com/MinterTeam/minter-hub-connector/minter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"math"
	"os"
	"strconv"
	"time"
)

const threshold = 667

var cfg = config.Get()

func main() {
	cosmos.Setup()
	orcAddress, orcPriv := cosmos.GetAccount(cfg.Cosmos.Mnemonic)
	minterWallet, err := wallet.Create(cfg.Minter.Mnemonic, "")
	if err != nil {
		panic(err)
	}

	minterClient, err := http_client.New(cfg.Minter.NodeUrl)
	if err != nil {
		panic(err)
	}

	cosmosConn, err := grpc.DialContext(c.Background(), cfg.Cosmos.NodeGrpcUrl, grpc.WithInsecure(), grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: time.Second * 5,
	}))

	ctx := context.Context{
		LastCheckedMinterBlock: cfg.Minter.StartBlock,
		LastEventNonce:         cfg.Minter.StartEventNonce,
		LastBatchNonce:         cfg.Minter.StartBatchNonce,
		LastValsetNonce:        cfg.Minter.StartValsetNonce,
		MinterMultisigAddr:     cfg.Minter.MultisigAddr,
		CosmosConn:             cosmosConn,
		MinterClient:           minterClient,
		OrcAddress:             orcAddress,
		OrcPriv:                orcPriv,
		MinterWallet:           minterWallet,
		Logger:                 log.NewTMLogger(os.Stdout),
	}

	ctx.Logger.Info("Syncing with Minter")

	ctx = minter.GetLatestMinterBlockAndNonce(ctx, cosmos.GetLastMinterNonce(orcAddress.String(), cosmosConn))

	ctx.Logger.Info("Starting with block", ctx.LastCheckedMinterBlock, "event nonce", ctx.LastEventNonce, "batch nonce", ctx.LastBatchNonce, "valset nonce", ctx.LastValsetNonce)

	if true { // todo: check if we have address
		privateKey, err := ethCrypto.HexToECDSA(minterWallet.PrivateKey)
		if err != nil {
			panic(err)
		}

		hash := ethCrypto.Keccak256(orcAddress.Bytes())
		signature, err := types.NewMinterSignature(hash, privateKey)
		if err != nil {
			panic("signing cosmos address with Minter key")
		}
		// You've got to do all this to get an Eth address from the private key
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			panic("casting public key to ECDSA")
		}
		minterAddress := ethCrypto.PubkeyToAddress(*publicKeyECDSA)

		cosmos.SendCosmosTx([]sdk.Msg{
			types.NewMsgSetMinterAddress("Mx"+minterAddress.String()[2:], orcAddress, hex.EncodeToString(signature)),
		}, orcAddress, orcPriv, cosmosConn, ctx.Logger)

		go cosmos.SendCosmosTx([]sdk.Msg{
			types.NewMsgValsetRequest(orcAddress),
		}, orcAddress, orcPriv, cosmosConn, ctx.Logger)
	}

	// main loop
	for {
		relayBatches(ctx)
		relayValsets(ctx)
		ctx = relayMinterEvents(ctx)

		ctx.Logger.Info("Last checked minter block", "height", ctx.LastCheckedMinterBlock, "event nonce", ctx.LastEventNonce, "batch nonce", ctx.LastBatchNonce, "valset nonce", ctx.LastValsetNonce)
		time.Sleep(2 * time.Second)
	}
}

func relayBatches(ctx context.Context) {
	cosmosClient := types.NewQueryClient(ctx.CosmosConn)

	{
		response, err := cosmosClient.LastPendingBatchRequestByAddr(c.Background(), &types.QueryLastPendingBatchRequestByAddrRequest{
			Address: ctx.OrcAddress.String(),
		})
		if err != nil {
			ctx.Logger.Error("Error while getting last pending batch", "err", err.Error())
			return
		}

		if response.Batch != nil {
			ctx.Logger.Info("Sending batch confirm", "batch nonce", response.Batch.BatchNonce)

			txData := transaction.NewMultisendData()
			for _, out := range response.Batch.Transactions {
				txData.AddItem(transaction.NewSendData().SetCoin(out.MinterToken.CoinId).MustSetTo(out.DestAddress).SetValue(out.MinterToken.Amount.BigInt()))
			}

			tx, _ := transaction.NewBuilder(cfg.Minter.ChainID).NewTransaction(txData)
			signedTx, _ := tx.SetNonce(response.Batch.MinterNonce).SetGasPrice(1).SetGasCoin(0).SetSignatureType(transaction.SignatureTypeMulti).Sign(
				cfg.Minter.MultisigAddr,
				ctx.MinterWallet.PrivateKey,
			)

			sigData, err := signedTx.SingleSignatureData()
			if err != nil {
				panic(err)
			}

			msg := &types.MsgConfirmBatch{
				Nonce:        response.Batch.BatchNonce,
				MinterSigner: ctx.MinterWallet.Address,
				Validator:    ctx.OrcAddress.String(),
				Signature:    sigData,
			}

			cosmos.SendCosmosTx([]sdk.Msg{msg}, ctx.OrcAddress, ctx.OrcPriv, ctx.CosmosConn, ctx.Logger)
		}
	}

	latestBatches, err := cosmosClient.OutgoingTxBatches(c.Background(), &types.QueryOutgoingTxBatchesRequest{})
	if err != nil {
		ctx.Logger.Error("Error getting last batches", "err", err.Error())
		return
	}

	var oldestSignedBatch *types.OutgoingTxBatch
	var oldestSignatures []*types.MsgConfirmBatch

	for _, batch := range latestBatches.Batches {
		sigs, err := cosmosClient.BatchConfirms(c.Background(), &types.QueryBatchConfirmsRequest{
			Nonce: batch.BatchNonce,
		})
		if err != nil {
			ctx.Logger.Error("Error while getting batch confirms", "err", err.Error())
			return
		}

		if sigs.Size() > 0 { // todo: check if we have enough votes
			oldestSignedBatch = batch
			oldestSignatures = sigs.Confirms
		}
	}

	if oldestSignedBatch == nil {
		return
	}

	if oldestSignedBatch.BatchNonce < ctx.LastBatchNonce {
		return
	}

	ctx.Logger.Info("Sending batch to Minter")

	txData := transaction.NewMultisendData()
	for _, out := range oldestSignedBatch.Transactions {
		txData.AddItem(transaction.NewSendData().SetCoin(out.MinterToken.CoinId).MustSetTo(out.DestAddress).SetValue(out.MinterToken.Amount.BigInt()))
	}

	tx, _ := transaction.NewBuilder(cfg.Minter.ChainID).NewTransaction(txData)
	tx.SetNonce(oldestSignedBatch.MinterNonce).SetGasPrice(1).SetGasCoin(0).SetSignatureType(transaction.SignatureTypeMulti)
	signedTx, err := tx.Sign(cfg.Minter.MultisigAddr)
	if err != nil {
		panic(err)
	}

	for _, sig := range oldestSignatures {
		signedTx, err = signedTx.AddSignature(sig.Signature)
		if err != nil {
			panic(err)
		}
	}

	encodedTx, err := signedTx.Encode()
	if err != nil {
		panic(err)
	}

	ctx.Logger.Debug("Batch tx", "tx", encodedTx)
	response, err := ctx.MinterClient.SendTransaction(encodedTx)
	if err != nil {
		code, body, err := http_client.ErrorBody(err)
		if err != nil {
			ctx.Logger.Error("Error on sending Minter Tx", "err", err.Error())
		} else {
			ctx.Logger.Error("Error on sending Minter Tx", "code", code, "err", body.Error.Message)
		}
	} else if response.Code != 0 {
		ctx.Logger.Error("Error on sending Minter Tx", "err", response.Log)
	}
}

func relayValsets(ctx context.Context) {
	cosmosClient := types.NewQueryClient(ctx.CosmosConn)

	{
		response, err := cosmosClient.LastPendingValsetRequestByAddr(c.Background(), &types.QueryLastPendingValsetRequestByAddrRequest{
			Address: ctx.OrcAddress.String(),
		})
		if err != nil {
			ctx.Logger.Error("Error while getting last pending valset", "err", err.Error())
			return
		}

		if response.Valset != nil {
			ctx.Logger.Info("Sending valset confirm", "valset nonce", response.Valset.Nonce)

			txData := transaction.NewEditMultisigData()
			txData.Threshold = threshold

			totalPower := uint64(0)
			for _, val := range response.Valset.Members {
				totalPower += val.Power
			}

			for _, val := range response.Valset.Members {
				var addr transaction.Address
				bytes, _ := wallet.AddressToHex(val.MinterAddress)
				copy(addr[:], bytes)

				weight := uint32(sdk.NewUint(val.Power).MulUint64(1000).QuoUint64(totalPower).Uint64())

				txData.Addresses = append(txData.Addresses, addr)
				txData.Weights = append(txData.Weights, weight)
			}

			tx, _ := transaction.NewBuilder(cfg.Minter.ChainID).NewTransaction(txData)
			tx.SetPayload([]byte(strconv.Itoa(int(response.Valset.Nonce))))
			signedTx, _ := tx.SetNonce(response.Valset.MinterNonce).SetGasPrice(1).SetGasCoin(0).SetSignatureType(transaction.SignatureTypeMulti).Sign(
				cfg.Minter.MultisigAddr,
				ctx.MinterWallet.PrivateKey,
			)

			sigData, err := signedTx.SingleSignatureData()
			if err != nil {
				panic(err)
			}

			msg := &types.MsgValsetConfirm{
				Nonce:         response.Valset.Nonce,
				Validator:     ctx.OrcAddress.String(),
				MinterAddress: ctx.MinterWallet.Address,
				Signature:     sigData,
			}

			cosmos.SendCosmosTx([]sdk.Msg{msg}, ctx.OrcAddress, ctx.OrcPriv, ctx.CosmosConn, ctx.Logger)
		}
	}

	latestValsets, err := cosmosClient.LastValsetRequests(c.Background(), &types.QueryLastValsetRequestsRequest{})
	if err != nil {
		ctx.Logger.Error("Error on getting last valset requests", "err", err.Error())
		return
	}

	var oldestSignedValset *types.Valset
	var oldestSignatures []*types.MsgValsetConfirm

	for _, valset := range latestValsets.Valsets {
		sigs, err := cosmosClient.ValsetConfirmsByNonce(c.Background(), &types.QueryValsetConfirmsByNonceRequest{
			Nonce: valset.Nonce,
		})
		if err != nil {
			ctx.Logger.Error("Error while getting valset confirms", "err", err.Error())
			return
		}

		if sigs.Size() > 0 { // todo: check if we have enough votes
			oldestSignedValset = valset
			oldestSignatures = sigs.Confirms
		}
	}

	if oldestSignedValset == nil {
		return
	}

	if oldestSignedValset.Nonce <= ctx.LastValsetNonce {
		return
	}

	ctx.Logger.Info("Sending valset to Minter")

	txData := transaction.NewEditMultisigData()
	txData.Threshold = threshold

	totalPower := uint64(0)
	for _, val := range oldestSignedValset.Members {
		totalPower += val.Power
	}

	for _, val := range oldestSignedValset.Members {
		var addr transaction.Address
		bytes, _ := wallet.AddressToHex(val.MinterAddress)
		copy(addr[:], bytes)

		weight := uint32(sdk.NewUint(val.Power).MulUint64(1000).QuoUint64(totalPower).Uint64())

		txData.Addresses = append(txData.Addresses, addr)
		txData.Weights = append(txData.Weights, weight)
	}

	tx, _ := transaction.NewBuilder(cfg.Minter.ChainID).NewTransaction(txData)
	tx.SetNonce(oldestSignedValset.MinterNonce).SetGasPrice(1).SetGasCoin(0).SetSignatureType(transaction.SignatureTypeMulti)
	tx.SetPayload([]byte(strconv.Itoa(int(oldestSignedValset.Nonce))))
	signedTx, err := tx.Sign(cfg.Minter.MultisigAddr)
	if err != nil {
		panic(err)
	}

	for _, sig := range oldestSignatures {
		signedTx, err = signedTx.AddSignature(sig.Signature)
		if err != nil {
			panic(err)
		}
	}

	encodedTx, err := signedTx.Encode()
	if err != nil {
		panic(err)
	}

	ctx.Logger.Debug("Valset update tx", "tx", encodedTx)
	response, err := ctx.MinterClient.SendTransaction(encodedTx)
	if err != nil {
		code, body, err := http_client.ErrorBody(err)
		if err != nil {
			ctx.Logger.Error("Error on sending Minter Tx", "err", err.Error())
		} else {
			ctx.Logger.Error("Error on sending Minter Tx", "code", code, "err", body.Error.Message)
		}
	} else if response.Code != 0 {
		ctx.Logger.Error("Error on sending Minter Tx", "err", response.Log)
	}
}

func relayMinterEvents(ctx context.Context) context.Context {
	latestBlock := minter.GetLatestMinterBlock(ctx.MinterClient, ctx.Logger)
	if latestBlock-ctx.LastCheckedMinterBlock > 100 {
		latestBlock = ctx.LastCheckedMinterBlock + 100
	}

	oracleClient := oracleTypes.NewQueryClient(ctx.CosmosConn)
	coinList, err := oracleClient.Coins(c.Background(), &oracleTypes.QueryCoinsRequest{})
	if err != nil {
		ctx.Logger.Info("Error getting coins from hub", "err", err.Error())
		time.Sleep(time.Second)
		return ctx
	}

	var deposits []cosmos.Deposit
	var batches []cosmos.Batch
	var valsets []cosmos.Valset

	const blocksPerBatch = 100
	for i := uint64(0); i <= uint64(math.Ceil(float64(latestBlock-ctx.LastCheckedMinterBlock)/blocksPerBatch)); i++ {
		from := ctx.LastCheckedMinterBlock + 1 + i*blocksPerBatch
		to := ctx.LastCheckedMinterBlock + (i+1)*blocksPerBatch

		if to > latestBlock {
			to = latestBlock
		}

		blocks, err := ctx.MinterClient.Blocks(from, to, false)
		if err != nil {
			ctx.Logger.Info("Error getting minter blocks", "err", err.Error())
			time.Sleep(time.Second)
			i--
			continue
		}

		for _, block := range blocks.Blocks {
			ctx.LastCheckedMinterBlock = block.Height

			ctx.Logger.Debug("Checking block", "height", block.Height)
			for _, tx := range block.Transactions {
				if tx.Type == uint64(transaction.TypeSend) {
					data, _ := tx.Data.UnmarshalNew()
					sendData := data.(*models.SendData)
					if sendData.To != cfg.Minter.MultisigAddr {
						continue
					}

					cmd := command.Command{}
					json.Unmarshal(tx.Payload, &cmd)

					value, _ := sdk.NewIntFromString(sendData.Value)

					if err := cmd.Validate(value); err != nil {
						ctx.Logger.Error("Cannot validate incoming tx", "err", err.Error())
						continue
					}

					for _, hubCoin := range coinList.GetCoins() {
						if sendData.Coin.ID == hubCoin.MinterId {
							ctx.Logger.Info("Found new deposit", "from", tx.From, "to", string(tx.Payload), "amount", sendData.Value, "coin", sendData.Coin.ID)
							deposits = append(deposits, cosmos.Deposit{
								Sender:     tx.From,
								Type:       cmd.Type,
								Recipient:  cmd.Recipient,
								Amount:     sendData.Value,
								Fee:        cmd.Fee,
								EventNonce: ctx.LastEventNonce,
								CoinID:     sendData.Coin.ID,
								TxHash:     tx.Hash,
							})

							ctx.LastEventNonce++
						}
					}
				}

				if tx.Type == uint64(transaction.TypeMultisend) && tx.From == cfg.Minter.MultisigAddr {
					ctx.Logger.Info("Found withdrawal")
					batches = append(batches, cosmos.Batch{
						BatchNonce: ctx.LastBatchNonce,
						EventNonce: ctx.LastEventNonce,
						TxHash:     tx.Hash,
					})

					ctx.LastEventNonce++
					ctx.LastBatchNonce++
				}

				if tx.Type == uint64(transaction.TypeEditMultisig) && tx.From == cfg.Minter.MultisigAddr {
					ctx.Logger.Info("Found valset update")

					nonce, err := strconv.Atoi(string(tx.Payload))
					if err != nil {
						ctx.Logger.Error("Error while decoding valset update nonce", "err", err.Error())
					} else {
						valsets = append(valsets, cosmos.Valset{
							ValsetNonce: uint64(nonce),
							EventNonce:  ctx.LastEventNonce,
						})

						ctx.LastEventNonce++
						ctx.LastValsetNonce = uint64(nonce)
					}
				}
			}
		}
	}

	if len(deposits) > 0 || len(batches) > 0 || len(valsets) > 0 {
		cosmos.SendCosmosTx(cosmos.CreateClaims(ctx.CosmosConn, ctx.OrcAddress, deposits, batches, valsets, ctx.Logger), ctx.OrcAddress, ctx.OrcPriv, ctx.CosmosConn, ctx.Logger)
	}

	return ctx
}
