package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx client.Context, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/prices", storeName), pricesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/min_eth_fee", storeName), minEthFeeHandler(cliCtx, storeName)).Methods("GET")
}
