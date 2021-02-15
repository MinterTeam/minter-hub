package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/mhub/chain/coins"
	"github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-hub-oracle/config"
	"github.com/MinterTeam/minter-hub-oracle/cosmos"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"time"
)

var cfg = config.Get()

const multiplier = 1e10

var pipInBip = sdk.NewInt(1000000000000000000)

func main() {
	cosmos.Setup()
	orcAddress, orcPriv := cosmos.GetAccount(cfg.Cosmos.Mnemonic)

	println(orcAddress.String())

	minterClient, err := http_client.New(cfg.Minter.NodeUrl)
	if err != nil {
		panic(err)
	}

	cosmosConn, err := grpc.DialContext(context.Background(), cfg.Cosmos.NodeGrpcUrl, grpc.WithInsecure(), grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: time.Second * 5,
	}))

	if err != nil {
		panic(err)
	}
	defer cosmosConn.Close()

	for {
		relayPrices(minterClient, cosmosConn, orcAddress, orcPriv)

		time.Sleep(1 * time.Second)
	}
}

func relayPrices(minterClient *http_client.Client, cosmosConn *grpc.ClientConn, orcAddress sdk.AccAddress, orcPriv *secp256k1.PrivKey) {
	cosmosClient := types.NewQueryClient(cosmosConn)

	response, err := cosmosClient.CurrentEpoch(context.Background(), &types.QueryCurrentEpochRequest{})
	if err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		return
	}

	// check if already voted
	for _, vote := range response.GetEpoch().GetVotes() {
		if vote.Oracle == orcAddress.String() {
			return
		}
	}

	prices := &types.Prices{List: []*types.Price{}}

	basecoinPrice := getBasecoinPrice()
	for _, coin := range coins.GetCoins() {
		response, err := minterClient.EstimateCoinIDSell(0, coin.MinterID, pipInBip.String())
		if err != nil {
			code, payload, err := http_client.ErrorBody(err)
			if err != nil {
				println(err.Error())
			} else {
				println(code, payload.Error.Message)
			}

			time.Sleep(time.Second)
			return
		}

		priceInBasecoin, _ := sdk.NewIntFromString(response.WillGet)
		price := priceInBasecoin.Mul(basecoinPrice).Quo(pipInBip)

		prices.List = append(prices.List, &types.Price{
			Name:  fmt.Sprintf("minter/%d", coin.MinterID),
			Value: price,
		})
	}

	prices.List = append(prices.List, &types.Price{
		Name:  "eth/0",
		Value: getEthPrice(),
	})

	prices.List = append(prices.List, &types.Price{
		Name:  "eth/gas",
		Value: getEthGasPrice(),
	})

	msg := &types.MsgPriceClaim{
		Epoch:        response.Epoch.Nonce,
		Prices:       prices,
		Orchestrator: orcAddress.String(),
	}

	cosmos.SendCosmosTx([]sdk.Msg{msg}, orcAddress, orcPriv, cosmosConn)
}

func getBasecoinPrice() sdk.Int {
	_, body, err := fasthttp.Get(nil, "https://api.coingecko.com/api/v3/simple/price?ids=bip&vs_currencies=usd")
	if err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		return getBasecoinPrice()
	}
	var result CoingeckoResult
	if err := json.Unmarshal(body, &result); err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		return getBasecoinPrice()
	}

	bipPrice := result["bip"]["usd"]

	return sdk.NewInt(int64(bipPrice * multiplier)) // todo
}

func getEthPrice() sdk.Int {
	_, body, err := fasthttp.Get(nil, "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		return getEthPrice()
	}
	var result CoingeckoResult
	if err := json.Unmarshal(body, &result); err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		return getEthPrice()
	}

	return sdk.NewInt(int64(result["ethereum"]["usd"] * multiplier)) // todo
}

func getEthGasPrice() sdk.Int {
	_, body, err := fasthttp.Get(nil, "https://ethgasstation.info/api/ethgasAPI.json")
	if err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		return getEthGasPrice()
	}
	var result EthGasResult
	if err := json.Unmarshal(body, &result); err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		return getEthGasPrice()
	}

	return sdk.NewInt(result.Fast) // todo
}

type EthGasResult struct {
	Fast int64 `json:"fast"`
}

type CoingeckoResult map[string]map[string]float64
