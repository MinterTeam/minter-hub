use deep_space::address::Address;
use deep_space::public_key::PublicKey;
use serde::de::Deserializer;
use serde::{de, Deserialize};
use serde_json::Value;
use std::{fmt::Display, str::FromStr};

/// A generic wrapper for Cosmos REST server responses which always
/// include the height
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ResponseWrapper<T> {
    #[serde(deserialize_with = "parse_val")]
    pub height: u64,
    pub result: T,
}

/// A generic wrapper for Cosmos REST server responses which always
/// include the struct type
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct TypeWrapper<T> {
    #[serde(rename = "type")]
    pub struct_type: String,
    pub value: T,
}

#[derive(Serialize, Deserialize, Debug, Clone, Default)]
pub struct PubKeyWrapper {
    #[serde(rename = "type")]
    key_type: String,
    #[serde(deserialize_with = "parse_val")]
    value: PublicKey,
}

fn default_account_number() -> u64 {
    0
}
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct CosmosAccountInfo {
    #[serde(deserialize_with = "parse_val")]
    pub address: Address,
    pub public_key: Option<PubKeyWrapper>,
    #[serde(deserialize_with = "parse_val", default = "default_account_number")]
    pub sequence: u64,
    #[serde(default = "default_account_number", deserialize_with = "parse_val")]
    pub account_number: u64,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BlockId {
    pub hash: String,
    pub parts: BlockParts,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BlockParts {
    pub total: u64,
    pub hash: String,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BlockHeader {
    pub version: BlockVersion,
    pub chain_id: String,
    pub time: String,
    pub last_block_id: BlockId,
    pub last_commit_hash: String,
    pub data_hash: String,
    pub validators_hash: String,
    pub next_validators_hash: String,
    pub consensus_hash: String,
    pub app_hash: String,
    pub last_results_hash: String,
    pub evidence_hash: String,
    #[serde(deserialize_with = "parse_val")]
    pub proposer_address: Address,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BlockVersion {
    #[serde(deserialize_with = "parse_val")]
    pub block: u64,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct LatestBlockEndpointResponse {
    pub block_id: BlockId,
    pub block: Block,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct Block {
    pub header: BlockHeader,
    pub data: BlockData,
    pub evidence: BlockEvidence,
    pub last_commit: LastCommit,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BlockData {
    pub txs: Option<Vec<String>>,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BlockEvidence {
    pub evidence: Option<Vec<String>>,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct LastCommit {
    #[serde(deserialize_with = "parse_val")]
    pub height: u64,
    pub round: u64,
    pub block_id: BlockId,
    pub signatures: Vec<BlockSignature>,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BlockSignature {
    pub block_id_flag: u64,
    #[serde(deserialize_with = "parse_val")]
    pub validator_address: Address,
    pub timestamp: String,
    pub signature: String,
}

#[derive(Debug, Clone)]
pub struct OptionalTXInfo {
    pub chain_id: String,
    pub account_number: u64,
    pub sequence: u64,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct TXSendResponse {
    pub logs: Option<Value>,
    pub txhash: String,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct TxSendErrorResponse {
    pub code: u64,
    pub codespace: String,
    #[serde(deserialize_with = "parse_val_option", default)]
    pub gas_used: Option<u64>,
    pub logs: Option<String>,
    pub raw_log: String,
}

/// Adapter that lets us parse any val that implements from_str into
/// the type we want, this helps solve type problems like sigs or addresses
/// being presented as strings and requiring a parse. For our own types like
/// Address we just implement deserialize such that the string representation
/// is accepted implicitly. But for native types like u128 this is the only
/// way to go
pub fn parse_val<'de, T, D>(deserializer: D) -> Result<T, D::Error>
where
    T: FromStr,
    T::Err: Display,
    D: Deserializer<'de>,
{
    let s: String = String::deserialize(deserializer)?;
    T::from_str(&s).map_err(de::Error::custom)
}

fn parse_val_option<'de, T, D>(deserializer: D) -> Result<Option<T>, D::Error>
where
    T: FromStr,
    T::Err: Display,
    D: Deserializer<'de>,
{
    let s: String = String::deserialize(deserializer)?;
    match T::from_str(&s) {
        Ok(val) => Ok(Some(val)),
        Err(_e) => Ok(None),
    }
}

/// A blank struct, used to parse blank responses
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct Blank {}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs::read_to_string;

    #[test]
    fn decode_account_info() {
        let file = read_to_string("test_files/account_info_active.json")
            .expect("Failed to read test files!");

        let _decoded: ResponseWrapper<TypeWrapper<CosmosAccountInfo>> =
            serde_json::from_str(&file).unwrap();

        let file = read_to_string("test_files/account_info_has_tokens.json")
            .expect("Failed to read test files!");

        let _decoded: ResponseWrapper<TypeWrapper<CosmosAccountInfo>> =
            serde_json::from_str(&file).unwrap();
    }
}
