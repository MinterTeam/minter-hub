package minter

import (
	"context"
	"encoding/json"
	oracleTypes "github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client/models"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-hub-connector/command"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
	"math"
	"sort"
	"strconv"
	"time"
)

func GetLatestMinterBlock(client *http_client.Client) uint64 {
	status, err := client.Status()
	if err != nil {
		println(err.Error())
		time.Sleep(1 * time.Second)
		return GetLatestMinterBlock(client)
	}

	return status.LatestBlockHeight
}

func GetLatestMinterBlockAndNonce(cosmosConn *grpc.ClientConn, startMinterBlock uint64, startEventNonce uint64, startBatchNonce uint64, startValsetNonce uint64, multisigAddr string, currentNonce uint64, client *http_client.Client) (block, eventNonce, batchNonce, valsetNonce uint64) {
	println("Current nonce @ hub", currentNonce)

	latestBlock := GetLatestMinterBlock(client)

	eventNonce = startEventNonce
	batchNonce = startBatchNonce
	valsetNonce = startValsetNonce

	oracleClient := oracleTypes.NewQueryClient(cosmosConn)
	coinList, err := oracleClient.Coins(context.Background(), &oracleTypes.QueryCoinsRequest{})
	if err != nil {
		panic(err)
	}

	const blocksPerBatch = 100
	for i := uint64(0); i <= uint64(math.Ceil(float64(latestBlock-startMinterBlock)/blocksPerBatch)); i++ {
		from := startMinterBlock + 1 + i * blocksPerBatch
		to := startMinterBlock + (i + 1) * blocksPerBatch

		if to > latestBlock {
			to = latestBlock
		}

		println("\rScanning from", from, "to", to)

		blocks, err := client.Blocks(from, to, false)
		if err != nil {
			println("ERROR: ", err.Error())
			time.Sleep(time.Second)
			i--
			continue
		}

		sort.Slice(blocks.Blocks, func(i, j int) bool {
			return blocks.Blocks[i].Height < blocks.Blocks[j].Height
		})

		for _, block := range blocks.Blocks {
			println(block.Height, "of", latestBlock)
			for _, tx := range block.Transactions {
				if tx.Type == uint64(transaction.TypeSend) {
					data, _ := tx.Data.UnmarshalNew()
					sendData := data.(*models.SendData)
					cmd := command.Command{}
					json.Unmarshal(tx.Payload, &cmd)

					value, _ := sdk.NewIntFromString(sendData.Value)
					if sendData.To == multisigAddr && cmd.Validate(value) == nil {
						for _, c := range coinList.GetCoins() {
							if sendData.Coin.ID == c.MinterId {
								println("deposit")
								if currentNonce < eventNonce {
									return block.Height - 1, eventNonce, batchNonce, valsetNonce
								}

								eventNonce++
							}
						}
					}
				}

				if tx.Type == uint64(transaction.TypeMultisend) && tx.From == multisigAddr {
					println("batch")
					if currentNonce < eventNonce {
						return block.Height - 1, eventNonce, batchNonce, valsetNonce
					}

					eventNonce++
					batchNonce++
				}

				if tx.Type == uint64(transaction.TypeEditMultisig) && tx.From == multisigAddr {
					nonce, err := strconv.Atoi(string(tx.Payload))
					if err != nil {
						println("ERROR:", err.Error())
					} else {
						println("valset update")
						if currentNonce < eventNonce {
							return block.Height - 1, eventNonce, batchNonce, valsetNonce
						}

						valsetNonce = uint64(nonce)
						eventNonce++
					}
				}
			}
		}
	}

	return latestBlock, eventNonce, batchNonce, valsetNonce
}
