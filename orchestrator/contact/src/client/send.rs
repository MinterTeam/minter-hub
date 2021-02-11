use crate::client::Contact;
use crate::jsonrpc::error::JsonRpcError;
use crate::types::*;
use crate::utils::maybe_get_optional_tx_info;
use actix_web::client::ConnectError;
use actix_web::client::SendRequestError;
use deep_space::address::Address;
use deep_space::coin::Coin;
use deep_space::msg::{Msg, SendMsg};
use deep_space::private_key::PrivateKey;
use deep_space::stdfee::StdFee;
use deep_space::stdsignmsg::StdSignMsg;
use deep_space::transaction::Transaction;
use deep_space::transaction::TransactionSendType;
use serde::Deserialize;
use serde::Serialize;
use std::fmt::Debug;
use std::time::Instant;
use std::{clone::Clone, time::Duration};
use tokio::time::delay_for;

impl Contact {
    /// The advanced version of create_and_send transaction that expects you to
    /// perform your own signing and prep first.
    pub async fn send_transaction<M: Clone + Serialize>(
        &self,
        msg: Transaction<M>,
    ) -> Result<TXSendResponse, JsonRpcError> {
        self.jsonrpc_client
            .request_method("txs", Some(msg), self.timeout, None)
            .await
    }

    /// When a transaction is in 'block' mode it actually asynchronously waits to go into the blockchain
    /// before returning. This is very useful in many contexts but is somewhat limited by the fact that
    /// nodes by default are configured to time out after 10 seconds. The caller of Contact of course
    /// expects the timeout they provide to be honored. This routine allows us to do that, retrying
    /// as needed until we reach the specific timeout allowed.
    pub async fn retry_on_block<
        M: Clone + Serialize,
        T: 'static + for<'de> Deserialize<'de> + Debug,
    >(
        &self,
        tx: Transaction<M>,
    ) -> Result<T, JsonRpcError> {
        if let Transaction::Block(..) = tx {
            let start = Instant::now();
            let mut res = self
                .jsonrpc_client
                .request_method("txs", Some(tx.clone()), self.timeout, None)
                .await;
            trace!("Sending tx got {:?}", res);
            while let Err(JsonRpcError::FailedToSend(SendRequestError::Connect(
                ConnectError::Disconnected,
            )))
            | Err(JsonRpcError::BadResponse(_))
            | Err(JsonRpcError::BadStruct(_)) = res
            {
                // since we can't combine logical statements and destructuring with let bindings
                // this will have to do
                if Instant::now() - start > self.timeout {
                    break;
                }
                // subtract two durations to get how much time we have left until
                // the actual user provided timeout. This will be passed as the call timeout
                // we must consider the case where the remote server does not have a short timeout
                // but our call fails for some other reason and we then get stuck waiting beyond
                // the expected timeout duration.
                let time_left = self.timeout - (Instant::now() - start);
                delay_for(Duration::from_secs(1)).await;
                res = self
                    .jsonrpc_client
                    .request_method("txs", Some(tx.clone()), time_left, None)
                    .await;
            }
            res
        } else {
            self.jsonrpc_client
                .request_method("txs", Some(tx.clone()), self.timeout, None)
                .await
        }
    }

    /// The hand holding version of send transaction that does it all for you
    #[allow(clippy::too_many_arguments)]
    pub async fn create_and_send_transaction(
        &self,
        coin: Coin,
        fee: Coin,
        destination: Address,
        private_key: PrivateKey,
        chain_id: Option<String>,
        account_number: Option<u64>,
        sequence: Option<u64>,
    ) -> Result<TXSendResponse, JsonRpcError> {
        trace!("Creating transaction");
        let our_address = private_key
            .to_public_key()
            .expect("Invalid private key!")
            .to_address();

        let tx_info =
            maybe_get_optional_tx_info(our_address, chain_id, account_number, sequence, &self)
                .await?;

        let std_sign_msg = StdSignMsg {
            chain_id: tx_info.chain_id,
            // this is not actually used in signing so we make
            // a best effort to provide it
            account_number: tx_info.account_number,
            sequence: tx_info.sequence,
            fee: StdFee {
                amount: vec![fee],
                gas: 500_000u64.into(),
            },
            msgs: vec![Msg::SendMsg(SendMsg {
                from_address: our_address,
                to_address: destination,
                amount: vec![coin],
            })],
            memo: String::new(),
        };

        let tx = private_key
            .sign_std_msg(std_sign_msg, TransactionSendType::Block)
            .unwrap();
        trace!("{}", json!(tx));

        self.retry_on_block(tx).await
    }
}
