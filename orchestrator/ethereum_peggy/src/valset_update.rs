use crate::utils::get_valset_nonce;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use peggy_utils::error::PeggyError;
use peggy_utils::types::*;
use std::time::Duration;
use web30::client::Web3;
use web30::types::{SendTxOption, TransactionRequest};

/// this function generates an appropriate Ethereum transaction
/// to submit the provided validator set and signatures.
pub async fn send_eth_valset_update(
    new_valset: Valset,
    old_valset: Valset,
    confirms: &[ValsetConfirmResponse],
    web3: &Web3,
    timeout: Duration,
    peggy_contract_address: EthAddress,
    our_eth_key: EthPrivateKey,
) -> Result<(), PeggyError> {
    let (old_addresses, old_powers) = old_valset.filter_empty_addresses();
    let (new_addresses, new_powers) = new_valset.filter_empty_addresses();
    let old_nonce = old_valset.nonce;
    let new_nonce = new_valset.nonce;
    assert!(new_nonce > old_nonce);
    let eth_address = our_eth_key.to_public_key().unwrap();
    info!(
        "Ordering signatures and submitting validator set {} -> {} update to Ethereum",
        old_nonce, new_nonce
    );

    // we need to use the old valset here because our signatures need to match the current
    // members of the validator set in the contract.
    let sig_data = old_valset.order_valset_sigs(confirms)?;
    let sig_arrays = to_arrays(sig_data);

    // Solidity function signature
    // function updateValset(
    // // The new version of the validator set
    // address[] memory _newValidators,
    // uint256[] memory _newPowers,
    // uint256 _newValsetNonce,
    // // The current validators that approve the change
    // address[] memory _currentValidators,
    // uint256[] memory _currentPowers,
    // uint256 _currentValsetNonce,
    // // These are arrays of the parts of the current validator's signatures
    // uint8[] memory _v,
    // bytes32[] memory _r,
    // bytes32[] memory _s
    let tokens = &[
        new_addresses.into(),
        new_powers.into(),
        new_nonce.into(),
        old_addresses.into(),
        old_powers.into(),
        old_nonce.into(),
        sig_arrays.v,
        sig_arrays.r,
        sig_arrays.s,
    ];
    let payload = clarity::abi::encode_call("updateValset(address[],uint256[],uint256,address[],uint256[],uint256,uint8[],bytes32[],bytes32[])",
    tokens).unwrap();

    let before_nonce = get_valset_nonce(peggy_contract_address, eth_address, web3).await?;
    if before_nonce != old_nonce {
        info!(
            "Someone else updated the valset to {}, exiting early",
            before_nonce
        );
        return Ok(());
    }

    let estimate_result = web3.eth_estimate_gas(TransactionRequest {
        from: Some(eth_address),
        to: peggy_contract_address,
        nonce: None,
        gas_price: None,
        gas: None,
        value: Some(0u64.into()),
        data: Some(payload.clone().into()),
    }).await;

    let gas = match estimate_result {
        Ok(gas) => {
            if gas.gt(&1_000_000u64.into()) {
                error!("Error while sending tx: gas limit is too high, possibly trying to send failing tx {}", gas);
            }

            gas
        }
        Err(e) => {
            error!("Error while sending tx: {}", e);

            return Err(PeggyError::EthereumRestError(e));
        }
    };

    let tx = web3
        .send_transaction(
            peggy_contract_address,
            payload,
            0u32.into(),
            eth_address,
            our_eth_key,
            vec![SendTxOption::GasLimit(gas.into())],
        )
        .await?;
    info!("Sent valset update with txid {:#066x}", tx);

    // TODO this segment of code works around the race condition for submitting valsets mostly
    // by not caring if our own submission reverts and only checking if the valset has been updated
    // period not if our update succeeded in particular. This will require some further consideration
    // in the future as many independent relayers racing to update the same thing will hopefully
    // be the common case.
    web3.wait_for_transaction(tx, timeout, None).await?;

    let last_nonce = get_valset_nonce(peggy_contract_address, eth_address, web3).await?;
    if last_nonce != new_nonce {
        error!(
            "Current nonce is {} expected to update to nonce {}",
            last_nonce, new_nonce
        );
    } else {
        info!(
            "Successfully updated Valset with new Nonce {:?}",
            last_nonce
        );
    }
    Ok(())
}
