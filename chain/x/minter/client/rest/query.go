package rest

import (
	"fmt"
	"net/http"

	"github.com/MinterTeam/mhub/chain/x/minter/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func getValsetRequestHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetRequest/%s", storeName, nonce))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		var out types.Valset
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// USED BY RUST
func batchByNonceHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/batch/%s", storeName, nonce))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		var out types.OutgoingTxBatch
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// USED BY RUST
func lastBatchesHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastBatches", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// gets all the confirm messages for a given validator set nonce
func allValsetConfirmsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetConfirms/%s", storeName, nonce))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset confirms not found")
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// gets all the confirm messages for a given transaction batch
func allBatchConfirmsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/batchConfirms/%s", storeName, nonce))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset confirms not found")
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastValsetRequestsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastValsetRequests", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset requests not found")
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastValsetRequestsByAddressHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingValsetRequest/%s", storeName, operatorAddr))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no pending valset requests found")
			return
		}

		var out types.Valset
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastBatchesByAddressHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingBatchRequest/%s", storeName, operatorAddr))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no pending valset requests found")
			return
		}

		var out types.OutgoingTxBatch
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func currentValsetHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/currentValset", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var out types.Valset
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func txStatusHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		node, err := cliCtx.GetNode()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		status, err := node.Status(r.Context())
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		height := status.SyncInfo.LatestBlockHeight

		vars := mux.Vars(r)
		tx_hash := vars[txHash]

		statuses := []string{
			"eth_outgoing_batch_executed",
			"eth_outgoing_batch",
			"minter_refund",
			"eth_refund",
			"minter_deposit_received",
		}

		for _, status := range statuses {
			result, err := node.TxSearch(r.Context(), fmt.Sprintf("%s.tx_hash='%s'", status, tx_hash), false, nil, nil, "")
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			if result.TotalCount > 0 {
				if status == "eth_outgoing_batch_executed" {
					var hash string
					for _, event := range result.Txs[0].TxResult.Events {
						for _, att := range event.Attributes {
							if string(att.GetKey()) == "batch_tx_hash" {
								hash = string(att.GetValue())
							}
						}
					}
					rest.PostProcessResponse(w, cliCtx.WithHeight(height), cliCtx.LegacyAmino.MustMarshalJSON(struct {
						Status string `json:"status"`
						TxHash string `json:"tx_hash"`
					}{
						Status: status,
						TxHash: hash,
					}))
					return
				}

				rest.PostProcessResponse(w, cliCtx.WithHeight(height), cliCtx.LegacyAmino.MustMarshalJSON(struct {
					Status string `json:"status"`
				}{
					Status: status,
				}))

				return
			}
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), "not found")
	}
}
