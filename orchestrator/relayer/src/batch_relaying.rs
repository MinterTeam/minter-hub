//! This module contains code for the batch update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use crate::find_latest_valset::find_latest_valset;
use clarity::address::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use cosmos_peggy::query::get_latest_transaction_batches;
use cosmos_peggy::query::get_transaction_batch_signatures;
use ethereum_peggy::submit_batch::send_eth_transaction_batch;
use ethereum_peggy::utils::get_tx_batch_nonce;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use std::ops::Add;
use std::time::Duration;
use tonic::transport::Channel;
use web30::client::Web3;

/// Check the last validator set on Ethereum, if it's lower than our latest validator
/// set then we should package and submit the update as an Ethereum transaction
pub async fn relay_batches(
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    mut grpc_client: &mut PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
    timeout: Duration,
) {
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();

    let latest_batches = get_latest_transaction_batches(grpc_client).await;
    trace!("Latest batches {:?}", latest_batches);
    if latest_batches.is_err() {
        return;
    }
    let mut latest_batches = latest_batches.unwrap();
    latest_batches.sort();
    //latest_batches.reverse();

    let nonce = web3.eth_get_transaction_count(our_ethereum_address).await;
    if nonce.is_err() {
        return;
    }
    let nonce = nonce.unwrap();

    let mut i = 0u32;

    for batch in latest_batches {
        let sigs =
            get_transaction_batch_signatures(grpc_client, batch.nonce, batch.token_contract).await;
        trace!("Got sigs {:?}", sigs);
        if let Ok(sigs) = sigs {
            // todo check that enough people have signed

            let erc20_contract = batch.token_contract;
            let latest_ethereum_batch = get_tx_batch_nonce(
                peggy_contract_address,
                erc20_contract,
                our_ethereum_address,
                web3,
            )
            .await;
            if latest_ethereum_batch.is_err() {
                error!(
                    "Failed to get latest Ethereum batch with {:?}",
                    latest_ethereum_batch
                );
            }
            let latest_ethereum_batch = latest_ethereum_batch.unwrap();

            if batch.clone().nonce > latest_ethereum_batch {
                info!(
                    "We have detected latest batch {} but latest on Ethereum is {} sending an update!",
                    batch.clone().nonce, latest_ethereum_batch
                );
                let current_valset = find_latest_valset(
                    &mut grpc_client,
                    our_ethereum_address,
                    peggy_contract_address,
                    web3,
                )
                .await;
                if let Ok(current_valset) = current_valset {
                    let current_nonce = nonce.clone().add(i.clone().into());
                    info!("Sending eth tx with nonce {}", current_nonce);

                    let _res = send_eth_transaction_batch(
                        current_valset,
                        batch,
                        &sigs,
                        web3,
                        timeout,
                        peggy_contract_address,
                        ethereum_key,
                        current_nonce,
                    )
                    .await;

                    i += 1;
                } else {
                    error!("Failed to find latest valset with {:?}", current_valset);
                }
            }
        } else {
            error!(
                "could not get signatures for {}:{} with {:?}",
                batch.token_contract, batch.nonce, sigs
            );
        }
    }
}
