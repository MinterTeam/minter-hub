package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx client.Context, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/coins", storeName), coinsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/prices", storeName), pricesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/eth_fee", storeName), ethFeeHandler(cliCtx, storeName)).Methods("GET")
}
