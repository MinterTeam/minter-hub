package rest

import (
	"fmt"
	"net/http"

	"github.com/MinterTeam/mhub/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func pricesHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/prices", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "prices not found")
			return
		}

		var out types.Prices
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func coinsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/coins", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "coins not found")
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func ethFeeHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/eth_fee", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "prices not found")
			return
		}

		var out types.QueryEthFeeResponse
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}
