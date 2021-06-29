package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-hub-oracle/config"
	"github.com/MinterTeam/minter-hub-oracle/cosmos"
	"github.com/MinterTeam/minter-hub-oracle/services/ethereum/gasprice"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

const multiplier = 1e10
const usdteCoinId = 1993

var pipInBip = sdk.NewInt(1000000000000000000)

func main() {
	logger := log.NewTMLogger(os.Stdout)

	cfg := config.Get()

	cosmos.Setup(cfg)

	orcAddress, orcPriv := cosmos.GetAccount(cfg.Cosmos.Mnemonic)

	logger.Info("Orc address", "address", orcAddress.String())

	minterClient, err := http_client.New(cfg.Minter.ApiAddr)
	if err != nil {
		panic(err)
	}

	cosmosConn, err := grpc.DialContext(context.Background(), cfg.Cosmos.GrpcAddr, grpc.WithInsecure(), grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: time.Second * 5,
	}))

	if err != nil {
		panic(err)
	}
	defer cosmosConn.Close()

	ethGasPrice, err := gasprice.NewService(cfg, logger)

	if err != nil {
		panic(err)
	}

	for {
		relayPrices(minterClient, ethGasPrice, cosmosConn, orcAddress, orcPriv, logger)

		time.Sleep(1 * time.Second)
	}
}

func relayPrices(
	minterClient *http_client.Client,
	ethGasPrice *gasprice.Service,
	cosmosConn *grpc.ClientConn,
	orcAddress sdk.AccAddress,
	orcPriv *secp256k1.PrivKey,
	logger log.Logger,
) {
	cosmosClient := types.NewQueryClient(cosmosConn)

	response, err := cosmosClient.CurrentEpoch(context.Background(), &types.QueryCurrentEpochRequest{})
	if err != nil {
		logger.Error("Error getting current epoch", "err", err.Error())
		time.Sleep(time.Second)
		return
	}

	// check if already voted
	for _, vote := range response.GetEpoch().GetVotes() {
		if vote.Oracle == orcAddress.String() {
			return
		}
	}

	coins, err := cosmosClient.Coins(context.Background(), &types.QueryCoinsRequest{})
	if err != nil {
		logger.Error("Error getting coins list", "err", err.Error())
		time.Sleep(time.Second)
		return
	}

	prices := &types.Prices{List: []*types.Price{}}

	basecoinPrice := getBasecoinPrice(logger, minterClient)
	ethPrice := getEthPrice(logger)

	for _, coin := range coins.GetCoins() {
		if coin.MinterId == 0 {
			prices.List = append(prices.List, &types.Price{
				Name:  fmt.Sprintf("minter/%d", coin.MinterId),
				Value: basecoinPrice,
			})

			continue
		}

		switch coin.Denom {
		case "usdt", "usdc", "busd", "dai", "ust", "pax", "tusd", "husd":
			prices.List = append(prices.List, &types.Price{
				Name:  fmt.Sprintf("minter/%d", coin.MinterId),
				Value: sdk.NewInt(10000000000),
			})
		case "wbtc":
			prices.List = append(prices.List, &types.Price{
				Name:  fmt.Sprintf("minter/%d", coin.MinterId),
				Value: getBitcoinPrice(logger),
			})
		case "weth":
			prices.List = append(prices.List, &types.Price{
				Name:  fmt.Sprintf("minter/%d", coin.MinterId),
				Value: ethPrice,
			})
		case "bnb":
			prices.List = append(prices.List, &types.Price{
				Name:  fmt.Sprintf("minter/%d", coin.MinterId),
				Value: getBnbPrice(logger),
			})
		default:
			var route []uint64
			if coin.Denom == "hubabuba" {
				const hubCoinId = 1902
				route = []uint64{hubCoinId}
			}

			response, err := minterClient.EstimateCoinIDSellExtended(0, uint64(coin.MinterId), pipInBip.String(), 0, "pool", route)
			if err != nil {
				_, payload, err := http_client.ErrorBody(err)
				if err != nil {
					logger.Error("Error estimating coin sell", "coin", coin.Denom, "err", err.Error())
				} else {
					logger.Error("Error estimating coin sell", "coin", coin.Denom, "err", payload.Error.Message)
				}

				time.Sleep(time.Second)
				return
			}

			priceInBasecoin, _ := sdk.NewIntFromString(response.WillGet)
			price := priceInBasecoin.Mul(basecoinPrice).Quo(pipInBip)

			prices.List = append(prices.List, &types.Price{
				Name:  fmt.Sprintf("minter/%d", coin.MinterId),
				Value: price,
			})
		}
	}

	prices.List = append(prices.List, &types.Price{
		Name:  "eth/0",
		Value: ethPrice,
	})

	prices.List = append(prices.List, &types.Price{
		Name:  "eth/gas",
		Value: ethGasPrice.GetGasPrice().Fast,
	})

	jsonPrices, _ := json.Marshal(prices.List)
	logger.Info("Prices", "val", jsonPrices)

	msg := &types.MsgPriceClaim{
		Epoch:        response.Epoch.Nonce,
		Prices:       prices,
		Orchestrator: orcAddress.String(),
	}

	cosmos.SendCosmosTx([]sdk.Msg{msg}, orcAddress, orcPriv, cosmosConn, logger)
}

func getBasecoinPrice(logger log.Logger, client *http_client.Client) sdk.Int {
	response, err := client.EstimateCoinIDSell(usdteCoinId, 0, pipInBip.String(), 0)
	if err != nil {
		_, payload, err := http_client.ErrorBody(err)
		if err != nil {
			logger.Error("Error estimating coin sell", "coin", "basecoin", "err", err.Error())
		} else {
			logger.Error("Error estimating coin sell", "coin", "basecoin", "err", payload.Error.Message)
		}

		time.Sleep(time.Second)

		return getBasecoinPrice(logger, client)
	}

	price, _ := sdk.NewIntFromString(response.WillGet)

	return price.Mul(sdk.NewInt(multiplier)).Quo(pipInBip)
}

func getEthPrice(logger log.Logger) sdk.Int {
	_, body, err := fasthttp.Get(nil, "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		logger.Error("Error getting eth price", "err", err.Error())
		time.Sleep(time.Second)
		return getEthPrice(logger)
	}
	var result CoingeckoResult
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("Error getting eth price", "err", err.Error())
		time.Sleep(time.Second)
		return getEthPrice(logger)
	}

	return sdk.NewInt(int64(result["ethereum"]["usd"] * multiplier)) // todo
}

func getBitcoinPrice(logger log.Logger) sdk.Int {
	_, body, err := fasthttp.Get(nil, "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	if err != nil {
		logger.Error("Error getting btc price", "err", err.Error())
		time.Sleep(time.Second)
		return getEthPrice(logger)
	}
	var result CoingeckoResult
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("Error getting btc price", "err", err.Error())
		time.Sleep(time.Second)
		return getEthPrice(logger)
	}

	return sdk.NewInt(int64(result["bitcoin"]["usd"] * multiplier)) // todo
}

func getBnbPrice(logger log.Logger) sdk.Int {
	_, body, err := fasthttp.Get(nil, "https://api.coingecko.com/api/v3/simple/price?ids=binancecoin&vs_currencies=usd")
	if err != nil {
		logger.Error("Error getting bnb price", "err", err.Error())
		time.Sleep(time.Second)
		return getEthPrice(logger)
	}
	var result CoingeckoResult
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("Error getting bnb price", "err", err.Error())
		time.Sleep(time.Second)
		return getEthPrice(logger)
	}

	return sdk.NewInt(int64(result["bitcoin"]["usd"] * multiplier)) // todo
}

type CoingeckoResult map[string]map[string]float64
