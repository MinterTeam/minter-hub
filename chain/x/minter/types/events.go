package types

const (
	EventTypeObservation              = "observation"
	EventTypeOutgoingBatch            = "outgoing_batch"
	EventTypeMultisigBootstrap        = "multisig_bootstrap"
	EventTypeMultisigUpdateRequest    = "multisig_update_request"
	EventTypeOutgoingBatchCanceled    = "outgoing_batch_canceled"
	EventTypeBridgeWithdrawalReceived = "withdrawal_received"
	EventTypeBridgeDepositReceived    = "deposit_received"
	EventTypeRefund                   = "refund"
	EventTypeWithdrawRequest          = "withdraw"
	EventTypeOutgoingBatchExecuted    = "batch_executed"
	EventTypeProcessAttestation       = "process_attestation"

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
