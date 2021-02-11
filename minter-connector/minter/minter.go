package minter

import (
	"encoding/json"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client/models"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-hub-connector/command"
	"github.com/althea-net/peggy/module/coins"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func GetLatestMinterBlockAndNonce(startMinterBlock uint64, startEventNonce uint64, startBatchNonce uint64, startValsetNonce uint64, multisigAddr string, currentNonce uint64, client *http_client.Client) (block, eventNonce, batchNonce, valsetNonce uint64) {
	latestBlock := GetLatestMinterBlock(client)

	eventNonce = startEventNonce
	batchNonce = startBatchNonce
	valsetNonce = startValsetNonce

	coinList := coins.GetCoins()

	for i := startMinterBlock; i <= latestBlock; i++ {
		println(i, "of", latestBlock)
		block, err := client.Block(i)
		if err != nil {
			println("ERROR: ", err.Error())
			time.Sleep(time.Second)
			i--
			continue
		}

		for _, tx := range block.Transactions {
			if tx.Type == uint64(transaction.TypeSend) {
				data, _ := tx.Data.UnmarshalNew()
				sendData := data.(*models.SendData)
				cmd := command.Command{}
				json.Unmarshal(tx.Payload, &cmd)

				value, _ := sdk.NewIntFromString(sendData.Value)
				if sendData.To == multisigAddr && cmd.Validate(value) == nil {
					for _, c := range coinList {
						if sendData.Coin.ID == c.MinterID {
							if currentNonce < eventNonce {
								return i - 1, eventNonce, batchNonce, valsetNonce
							}

							eventNonce++
						}
					}
				}
			}

			if tx.Type == uint64(transaction.TypeMultisend) && tx.From == multisigAddr {
				if currentNonce < eventNonce {
					return i - 1, eventNonce, batchNonce, valsetNonce
				}

				eventNonce++
				batchNonce++
			}

			if tx.Type == uint64(transaction.TypeEditMultisig) && tx.From == multisigAddr {
				nonce, err := strconv.Atoi(string(tx.Payload))
				if err != nil {
					println("ERROR:", err.Error())
				} else {
					if currentNonce < eventNonce {
						return i - 1, eventNonce, batchNonce, valsetNonce
					}

					valsetNonce = uint64(nonce)
					eventNonce++
				}
			}
		}
	}

	return latestBlock, eventNonce, batchNonce, valsetNonce
}
