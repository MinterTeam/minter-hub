use crate::client::Contact;
use crate::jsonrpc::error::JsonRpcError;
use crate::types::OptionalTXInfo;
use deep_space::address::Address;

/// retrieves 'optional' components of a transaction if not provided.
/// These are things like the chain ID, account_number, or sequence
/// that can be auto retrieved but may be troublesome to do so for
/// the caller.
pub async fn maybe_get_optional_tx_info(
    our_address: Address,
    chain_id: Option<String>,
    account_number: Option<u64>,
    sequence: Option<u64>,
    client: &Contact,
) -> Result<OptionalTXInfo, JsonRpcError> {
    // if the user provides values use those, otherwise fallback to retrieving them
    let (account_number, sequence) = if account_number.is_none() || sequence.is_none() {
        let info = client.get_account_info(our_address).await?;
        if info.result.value.is_none() {
            return Err(JsonRpcError::NoToken);
        }
        (
            info.result.value.clone().unwrap().account_number,
            info.result.value.unwrap().sequence,
        )
    } else {
        (account_number.unwrap(), sequence.unwrap())
    };

    // likewise with the chain id, if there's a user provided value
    // we can avoid the request

    let chain_id = if let Some(chain_id) = chain_id {
        chain_id
    } else {
        let block = client.get_latest_block().await?;
        block.block.header.chain_id
    };

    Ok(OptionalTXInfo {
        account_number,
        sequence,
        chain_id,
    })
}
