package rest

import (
	"fmt"
	"github.com/MinterTeam/mhub/chain/x/minter/client/utils"
	"github.com/MinterTeam/mhub/chain/x/minter/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/gorilla/mux"
)

const (
	nonce                  = "nonce"
	txHash                 = "tx_hash"
	bech32ValidatorAddress = "bech32ValidatorAddress"
	claimType              = "claimType"
	signType               = "signType"
)

// Here are the routes that are actually queried by the rust
// "peggy/valset_request/{}"
// "peggy/pending_valset_requests/{}"
// "peggy/valset_requests"
// "peggy/valset_confirm/{}"
// "peggy/pending_batch_requests/{}"
// "peggy/transaction_batches/"
// "peggy/signed_batches"

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx client.Context, r *mux.Router, storeName string) {

	/// Valsets

	// This endpoint gets all of the validator set confirmations for a given nonce. In order to determine if a valset is complete
	// the relayer queries the latest valsets and then compares the number of members they show versus the length of this endpoints output
	// if they match every validator has submitted a signature and we can go forward with relaying that validator set update.
	r.HandleFunc(fmt.Sprintf("/%s/valset_confirm/{%s}", storeName, nonce), allValsetConfirmsHandler(cliCtx, storeName)).Methods("GET")
	// gets the latest 5 validator set requests, used heavily by the relayer. Which hits this endpoint before checking which
	// of these last 5 have sufficient signatures to relay
	r.HandleFunc(fmt.Sprintf("/%s/valset_requests", storeName), lastValsetRequestsHandler(cliCtx, storeName)).Methods("GET")
	// Returns the last 'pending' (unsigned) validator set for a given validator address.
	r.HandleFunc(fmt.Sprintf("/%s/pending_valset_requests/{%s}", storeName, bech32ValidatorAddress), lastValsetRequestsByAddressHandler(cliCtx, storeName)).Methods("GET")
	// gets valset request by nonce, used to look up a specific valset. This is needed to lookup data about the current validator set on the contract
	// and determine what can or can not be submitted as a relayer
	r.HandleFunc(fmt.Sprintf("/%s/valset_request/{%s}", storeName, nonce), getValsetRequestHandler(cliCtx, storeName)).Methods("GET")
	// Provides the current validator set with powers and eth addresses, useful to check the current validator state
	// used to deploy the contract by the contract deployer script
	r.HandleFunc(fmt.Sprintf("/%s/current_valset", storeName), currentValsetHandler(cliCtx, storeName)).Methods("GET")

	/// Batches

	// The Ethereum signer queries this endpoint and signs whatever it returns once per loop iteration
	r.HandleFunc(fmt.Sprintf("/%s/pending_batch_requests/{%s}", storeName, bech32ValidatorAddress), lastBatchesByAddressHandler(cliCtx, storeName)).Methods("GET")
	// Gets all outgoing batches in the batch queue, up to 100
	r.HandleFunc(fmt.Sprintf("/%s/transaction_batches", storeName), lastBatchesHandler(cliCtx, storeName)).Methods("GET")
	// Gets a specific batch request from the outgoing queue by denom
	r.HandleFunc(fmt.Sprintf("/%s/transaction_batch/{%s}", storeName, nonce), batchByNonceHandler(cliCtx, storeName)).Methods("GET")
	// This endpoint gets all of the batch confirmations for a given nonce and denom In order to determine if a batch is complete
	// the relayer will compare the valset power on the contract to the number of signatures
	r.HandleFunc(fmt.Sprintf("/%s/batch_confirm/{%s}", storeName, nonce), allBatchConfirmsHandler(cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/tx_status/{%s}", storeName, txHash), txStatusHandler(cliCtx, storeName)).Methods("GET")
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the param
// change REST handler with a given sub-route.
func ProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "minter_cold_storage_transfer",
		Handler:  postProposalHandlerFn(clientCtx),
	}
}

func postProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.ColdStorageTransferProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewColdStorageTransferProposal(req.Amount)

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
