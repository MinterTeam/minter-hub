package types

const (
	EventTypeObservation              = "minter_observation"
	EventTypeOutgoingBatch            = "minter_outgoing_batch"
	EventTypeMultisigBootstrap        = "minter_multisig_bootstrap"
	EventTypeMultisigUpdateRequest    = "minter_multisig_update_request"
	EventTypeOutgoingBatchCanceled    = "minter_outgoing_batch_canceled"
	EventTypeBridgeWithdrawalReceived = "minter_withdrawal_received"
	EventTypeBridgeDepositReceived    = "minter_deposit_received"
	EventTypeRefund                   = "minter_refund"
	EventTypeWithdrawRequest          = "minter_withdraw"
	EventTypeOutgoingBatchExecuted    = "minter_batch_executed"
	EventTypeProcessAttestation       = "minter_process_attestation"

	AttributeKeyAttestationID   = "attestation_id"
	AttributeKeyMultisigID      = "multisig_id"
	AttributeKeyOutgoingBatchID = "batch_id"
	AttributeKeyOutgoingTXID    = "outgoing_tx_id"
	AttributeKeyAttestationType = "attestation_type"
	AttributeKeyContract        = "bridge_contract"
	AttributeKeyNonce           = "nonce"
	AttributeKeyBridgeChainID   = "bridge_chain_id"
	AttributeKeyTxHash          = "tx_hash"
)
