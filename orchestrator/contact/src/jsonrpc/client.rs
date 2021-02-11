use crate::{jsonrpc::error::JsonRpcError, types::TxSendErrorResponse};
use actix_web::client::Client;
use actix_web::http::header;
use serde::{Deserialize, Serialize};
use serde_json::from_value;
use serde_json::Value;
use std::str;
use std::time::Duration;

pub struct HTTPClient {
    url: String,
    client: Client,
}

impl HTTPClient {
    pub fn new(url: &str) -> Self {
        Self {
            url: url.to_string(),
            client: Client::default(),
        }
    }

    pub async fn request_method<T: Serialize, R: 'static>(
        &self,
        method: &str,
        params: Option<T>,
        timeout: Duration,
        request_size_limit: Option<usize>,
    ) -> Result<R, JsonRpcError>
    where
        for<'de> R: Deserialize<'de>,
        // T: std::fmt::Debug,
        R: std::fmt::Debug,
    {
        trace!(
            "About to make contact request to {} with payload {}",
            method,
            json!(params)
        );
        // the response payload size limit for this request, almost everything
        // will set this to None, and get the default 64k, but some requests
        // need bigger buffers (like full block requests)
        let limit = request_size_limit.unwrap_or(65536);
        let url_with_method = format!("{}/{}", self.url, method);
        // if we don't have a payload this is a get request
        let res = if let Some(params) = params {
            self.client
                .post(&url_with_method)
                .header(header::CONTENT_TYPE, "application/json")
                .timeout(timeout)
                .send_json(&params)
                .await
        } else {
            self.client
                .get(&url_with_method)
                .header(header::CONTENT_TYPE, "application/json")
                .timeout(timeout)
                .send()
                .await
        };
        let mut res = match res {
            Ok(val) => val,
            Err(e) => return Err(JsonRpcError::FailedToSend(e)),
        };
        let status = res.status();
        if !status.is_success() {
            return Err(JsonRpcError::BadResponse(format!(
                "Server Error {}",
                status
            )));
        }

        // this layer of error handling is not technically required, you could
        // replace this layer with a direct parse into Result<R, Error> but that's
        // much harder to debug since there's no way to actually display serde value
        // you're looking for.
        let json_value: Result<Value, _> = res.json().limit(limit).await;
        trace!("got Cosmos JSONRPC response {:#?}", json_value);
        let json: Value = match json_value {
            Ok(val) => val,
            Err(e) => return Err(JsonRpcError::BadResponse(e.to_string())),
        };
        let data: R = match from_value(json.clone()) {
            Ok(val) => val,
            Err(e) => {
                if let Ok(bad_tx_response) = from_value(json) {
                    let bad_tx_response: TxSendErrorResponse = bad_tx_response;
                    return Err(JsonRpcError::BadStruct(bad_tx_response.raw_log));
                }
                return Err(JsonRpcError::BadStruct(e.to_string()));
            }
        };

        Ok(data)
    }
}
