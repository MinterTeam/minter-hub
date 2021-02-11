package rest

import (
	"fmt"
	"net/http"

	"github.com/althea-net/peggy/module/x/oracle/types"
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

func minEthFeeHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/min_eth_fee", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "prices not found")
			return
		}

		var out types.QueryMinEthFeeResponse
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}
