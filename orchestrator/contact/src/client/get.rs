use crate::client::Contact;
use crate::jsonrpc::error::JsonRpcError;
use crate::types::*;
use deep_space::{address::Address, coin::Coin};

impl Contact {
    pub async fn get_latest_block_number(&self) -> Result<u64, JsonRpcError> {
        let none: Option<bool> = None;
        let res: Result<LatestBlockEndpointResponse, JsonRpcError> = self
            .jsonrpc_client
            .request_method("blocks/latest", none, self.timeout, None)
            .await;

        match res {
            Ok(res) => Ok(res.block.last_commit.height),
            Err(e) => Err(e),
        }
    }

    pub async fn get_latest_block(&self) -> Result<LatestBlockEndpointResponse, JsonRpcError> {
        let none: Option<bool> = None;
        self.jsonrpc_client
            .request_method("blocks/latest", none, self.timeout, None)
            .await
    }

    /// Gets account info for the provided Cosmos account using the accounts endpoint
    /// accounts do not have any info if they have no tokens or are otherwise never seen
    /// before an Ok(None) result indicates this
    pub async fn get_account_info(
        &self,
        address: Address,
    ) -> Result<ResponseWrapper<TypeWrapper<Option<CosmosAccountInfo>>>, JsonRpcError> {
        let none: Option<bool> = None;
        let res = self
            .jsonrpc_client
            .request_method(
                &format!("auth/accounts/{}", address),
                none,
                self.timeout,
                None,
            )
            .await;
        if let Err(JsonRpcError::BadStruct(_)) = res {
            let res: Result<ResponseWrapper<TypeWrapper<Blank>>, JsonRpcError> = self
                .jsonrpc_client
                .request_method(
                    &format!("auth/accounts/{}", address),
                    none,
                    self.timeout,
                    None,
                )
                .await;
            let res = res?;
            Ok(ResponseWrapper {
                height: res.height,
                result: TypeWrapper {
                    struct_type: res.result.struct_type,
                    value: None,
                },
            })
        } else {
            res
        }
    }

    pub async fn get_tx_by_hash(&self, txhash: &str) -> Result<TXSendResponse, JsonRpcError> {
        let none: Option<bool> = None;
        self.jsonrpc_client
            .request_method(&format!("txs/{}", txhash), none, self.timeout, None)
            .await
    }

    pub async fn get_balances(
        &self,
        address: Address,
    ) -> Result<ResponseWrapper<Vec<Coin>>, JsonRpcError> {
        let none: Option<bool> = None;
        self.jsonrpc_client
            .request_method(
                &format!("bank/balances/{}", address),
                none,
                self.timeout,
                None,
            )
            .await
    }
}
