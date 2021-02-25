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
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"math"
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

	println("Syncing with Minter")
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
	}

	ctx = minter.GetLatestMinterBlockAndNonce(ctx, cosmos.GetLastMinterNonce(orcAddress.String(), cosmosConn))
	println("Starting with block", ctx.LastCheckedMinterBlock, "event nonce", ctx.LastEventNonce, "batch nonce", ctx.LastBatchNonce, "valset nonce", ctx.LastValsetNonce)

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
		}, orcAddress, orcPriv, cosmosConn)

		go cosmos.SendCosmosTx([]sdk.Msg{
			types.NewMsgValsetRequest(orcAddress),
		}, orcAddress, orcPriv, cosmosConn)
	}

	// main loop
	for {
		relayBatches(ctx)
		relayValsets(ctx)
		ctx = relayMinterEvents(ctx)

		println("lastCheckedMinterBlock", ctx.LastCheckedMinterBlock, "event nonce", ctx.LastEventNonce, "batch nonce", ctx.LastBatchNonce, "valset nonce", ctx.LastValsetNonce)
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
			println("ERROR: ", err.Error())
			return
		}

		if response.Batch != nil {
			println("Sending batch confirm for", response.Batch.BatchNonce)

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

			cosmos.SendCosmosTx([]sdk.Msg{msg}, ctx.OrcAddress, ctx.OrcPriv, ctx.CosmosConn)
		}
	}

	latestBatches, err := cosmosClient.OutgoingTxBatches(c.Background(), &types.QueryOutgoingTxBatchesRequest{})
	if err != nil {
		println(err.Error())
		return
	}

	var oldestSignedBatch *types.OutgoingTxBatch
	var oldestSignatures []*types.MsgConfirmBatch

	for _, batch := range latestBatches.Batches {
		sigs, err := cosmosClient.BatchConfirms(c.Background(), &types.QueryBatchConfirmsRequest{
			Nonce: batch.BatchNonce,
		})
		if err != nil {
			println("ERROR: ", err.Error())
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

	println("Sending batch to Minter")

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

	println(encodedTx)
	response, err := ctx.MinterClient.SendTransaction(encodedTx)
	if err != nil {
		code, body, err := http_client.ErrorBody(err)
		println(code, body.Error.Message, err)
	} else if response.Code != 0 {
		println(response.Log)
	}
}

func relayValsets(ctx context.Context) {
	cosmosClient := types.NewQueryClient(ctx.CosmosConn)

	{
		response, err := cosmosClient.LastPendingValsetRequestByAddr(c.Background(), &types.QueryLastPendingValsetRequestByAddrRequest{
			Address: ctx.OrcAddress.String(),
		})
		if err != nil {
			println("ERROR: ", err.Error())
			return
		}

		if response.Valset != nil {
			println("Sending valset confirm for", response.Valset.Nonce)

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

			cosmos.SendCosmosTx([]sdk.Msg{msg}, ctx.OrcAddress, ctx.OrcPriv, ctx.CosmosConn)
		}
	}

	latestValsets, err := cosmosClient.LastValsetRequests(c.Background(), &types.QueryLastValsetRequestsRequest{})
	if err != nil {
		println(err.Error())
		return
	}

	var oldestSignedValset *types.Valset
	var oldestSignatures []*types.MsgValsetConfirm

	for _, valset := range latestValsets.Valsets {
		sigs, err := cosmosClient.ValsetConfirmsByNonce(c.Background(), &types.QueryValsetConfirmsByNonceRequest{
			Nonce: valset.Nonce,
		})
		if err != nil {
			println("ERROR: ", err.Error())
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

	println("Sending valset to Minter")

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

	println(encodedTx)
	response, err := ctx.MinterClient.SendTransaction(encodedTx)
	if err != nil {
		code, body, err2 := http_client.ErrorBody(err)
		if err2 != nil {
			println(err, err2)
		} else {
			println(code, body.Error.Message, err)
		}
	} else if response.Code != 0 {
		println(response.Log)
	}
}

func relayMinterEvents(ctx context.Context) context.Context {
	latestBlock := minter.GetLatestMinterBlock(ctx.MinterClient)
	if latestBlock-ctx.LastCheckedMinterBlock > 100 {
		latestBlock = ctx.LastCheckedMinterBlock + 100
	}

	oracleClient := oracleTypes.NewQueryClient(ctx.CosmosConn)
	coinList, err := oracleClient.Coins(c.Background(), &oracleTypes.QueryCoinsRequest{})
	if err != nil {
		println("ERROR: ", err.Error())
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
			println("ERROR: ", err.Error())
			time.Sleep(time.Second)
			i--
			continue
		}

		for _, block := range blocks.Blocks {
			ctx.LastCheckedMinterBlock = block.Height

			println("Checking block at height", block.Height)
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
						println(err.Error())
						continue
					}

					for _, hubCoin := range coinList.GetCoins() {
						if sendData.Coin.ID == hubCoin.MinterId {
							println("Found new deposit from", tx.From, "to", string(tx.Payload), "amount", sendData.Value, "coin", sendData.Coin.ID)
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
					println("Found withdrawal")
					batches = append(batches, cosmos.Batch{
						BatchNonce: ctx.LastBatchNonce,
						EventNonce: ctx.LastEventNonce,
						TxHash:     tx.Hash,
					})

					ctx.LastEventNonce++
					ctx.LastBatchNonce++
				}

				if tx.Type == uint64(transaction.TypeEditMultisig) && tx.From == cfg.Minter.MultisigAddr {
					println("Found valset update")

					nonce, err := strconv.Atoi(string(tx.Payload))
					if err != nil {
						println("ERROR:", err.Error())
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
		cosmos.SendCosmosTx(cosmos.CreateClaims(ctx.CosmosConn, ctx.OrcAddress, deposits, batches, valsets), ctx.OrcAddress, ctx.OrcPriv, ctx.CosmosConn)
	}

	return ctx
}
