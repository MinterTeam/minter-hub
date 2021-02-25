use crate::error::PeggyError;

/// This represents an individual transaction being bridged over to Ethereum
/// parallel is the OutgoingTransferTx in x/peggy/types/batch.go
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct Coin {
    pub denom: String,
    pub minter_id: u64,
    pub eth_addr: String,
}

impl Coin {
    pub fn from_proto(input: peggy_proto::oracle::Coin) -> Result<Self, PeggyError> {
        Ok(Coin {
            denom: input.denom,
            minter_id: input.minter_id,
            eth_addr: input.eth_addr,
        })
    }
}
