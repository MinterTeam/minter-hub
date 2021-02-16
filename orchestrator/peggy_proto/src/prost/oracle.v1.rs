/// Attestation is an aggregate of `claims` that eventually becomes `observed` by
/// all orchestrators
/// EVENT_NONCE:
/// EventNonce a nonce provided by the peggy contract that is unique per event fired
/// These event nonces must be relayed in order. This is a correctness issue,
/// if relaying out of order transaction replay attacks become possible
/// OBSERVED:
/// Observed indicates that >67% of validators have attested to the event,
/// and that the event should be executed by the peggy state machine
///
/// The actual content of the claims is passed in with the transaction making the claim
/// and then passed through the call stack alongside the attestation while it is processed
/// the key in which the attestation is stored is keyed on the exact details of the claim
/// but there is no reason to store those exact details becuause the next message sender
/// will kindly provide you with them.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Attestation {
    #[prost(uint64, tag="1")]
    pub epoch: u64,
    #[prost(bool, tag="2")]
    pub observed: bool,
    #[prost(string, repeated, tag="3")]
    pub votes: ::std::vec::Vec<std::string::String>,
    #[prost(bytes, tag="4")]
    pub claim_hash: std::vec::Vec<u8>,
    #[prost(uint64, tag="5")]
    pub height: u64,
}
/// ClaimType is the cosmos type of an event from the counterpart chain that can
/// be handled
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum ClaimType {
    Unknown = 0,
    Deposit = 1,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Prices {
    #[prost(message, repeated, tag="2")]
    pub list: ::std::vec::Vec<Price>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Price {
    #[prost(string, tag="1")]
    pub name: std::string::String,
    #[prost(string, tag="2")]
    pub value: std::string::String,
}
/// WithdrawClaim claims that a batch of withdrawal
/// operations on the bridge contract was executed.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgPriceClaim {
    #[prost(uint64, tag="1")]
    pub epoch: u64,
    #[prost(message, optional, tag="2")]
    pub prices: ::std::option::Option<Prices>,
    #[prost(string, tag="4")]
    pub orchestrator: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgPriceClaimResponse {
}
# [doc = r" Generated client implementations."] pub mod msg_client { # ! [allow (unused_variables , dead_code , missing_docs)] use tonic :: codegen :: * ; pub struct MsgClient < T > { inner : tonic :: client :: Grpc < T > , } impl MsgClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > MsgClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + HttpBody + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as HttpBody > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor (inner : T , interceptor : impl Into < tonic :: Interceptor >) -> Self { let inner = tonic :: client :: Grpc :: with_interceptor (inner , interceptor) ; Self { inner } } pub async fn price_claim (& mut self , request : impl tonic :: IntoRequest < super :: MsgPriceClaim > ,) -> Result < tonic :: Response < super :: MsgPriceClaimResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/oracle.v1.Msg/PriceClaim") ; self . inner . unary (request . into_request () , path , codec) . await } } impl < T : Clone > Clone for MsgClient < T > { fn clone (& self) -> Self { Self { inner : self . inner . clone () , } } } impl < T > std :: fmt :: Debug for MsgClient < T > { fn fmt (& self , f : & mut std :: fmt :: Formatter < '_ >) -> std :: fmt :: Result { write ! (f , "MsgClient {{ ... }}") } } }/// It's difficult to serialize and deserialize
/// interfaces, instead we can make this struct
/// that stores all the data the interface requires
/// and use it to store and then re-create a interface
/// object with all the same properties.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenericClaim {
    #[prost(uint64, tag="1")]
    pub epoch: u64,
    #[prost(int32, tag="2")]
    pub claim_type: i32,
    #[prost(bytes, tag="3")]
    pub hash: std::vec::Vec<u8>,
    #[prost(string, tag="4")]
    pub event_claimer: std::string::String,
    #[prost(oneof="generic_claim::Claim", tags="5")]
    pub claim: ::std::option::Option<generic_claim::Claim>,
}
pub mod generic_claim {
    #[derive(Clone, PartialEq, ::prost::Oneof)]
    pub enum Claim {
        #[prost(message, tag="5")]
        PriceClaim(super::MsgPriceClaim),
    }
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Epoch {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(message, repeated, tag="2")]
    pub votes: ::std::vec::Vec<Vote>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Vote {
    #[prost(string, tag="1")]
    pub oracle: std::string::String,
    #[prost(message, optional, tag="2")]
    pub claim: ::std::option::Option<MsgPriceClaim>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Coin {
    #[prost(string, tag="1")]
    pub denom: std::string::String,
    #[prost(string, tag="2")]
    pub eth_addr: std::string::String,
    #[prost(uint64, tag="3")]
    pub minter_id: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCurrentEpochRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCurrentEpochResponse {
    #[prost(message, optional, tag="1")]
    pub epoch: ::std::option::Option<Epoch>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCoinsRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCoinsResponse {
    #[prost(message, repeated, tag="1")]
    pub coins: ::std::vec::Vec<Coin>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryEthFeeRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryEthFeeResponse {
    #[prost(string, tag="1")]
    pub min: std::string::String,
    #[prost(string, tag="2")]
    pub fast: std::string::String,
}
# [doc = r" Generated client implementations."] pub mod query_client { # ! [allow (unused_variables , dead_code , missing_docs)] use tonic :: codegen :: * ; # [doc = " Query defines the gRPC querier service"] pub struct QueryClient < T > { inner : tonic :: client :: Grpc < T > , } impl QueryClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > QueryClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + HttpBody + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as HttpBody > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor (inner : T , interceptor : impl Into < tonic :: Interceptor >) -> Self { let inner = tonic :: client :: Grpc :: with_interceptor (inner , interceptor) ; Self { inner } } pub async fn current_epoch (& mut self , request : impl tonic :: IntoRequest < super :: QueryCurrentEpochRequest > ,) -> Result < tonic :: Response < super :: QueryCurrentEpochResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/oracle.v1.Query/CurrentEpoch") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn eth_fee (& mut self , request : impl tonic :: IntoRequest < super :: QueryEthFeeRequest > ,) -> Result < tonic :: Response < super :: QueryEthFeeResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/oracle.v1.Query/EthFee") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn coins (& mut self , request : impl tonic :: IntoRequest < super :: QueryCoinsRequest > ,) -> Result < tonic :: Response < super :: QueryCoinsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/oracle.v1.Query/Coins") ; self . inner . unary (request . into_request () , path , codec) . await } } impl < T : Clone > Clone for QueryClient < T > { fn clone (& self) -> Self { Self { inner : self . inner . clone () , } } } impl < T > std :: fmt :: Debug for QueryClient < T > { fn fmt (& self , f : & mut std :: fmt :: Formatter < '_ >) -> std :: fmt :: Result { write ! (f , "QueryClient {{ ... }}") } } }#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    #[prost(uint64, tag="1")]
    pub signed_claims_window: u64,
    #[prost(bytes, tag="2")]
    pub slash_fraction_claim: std::vec::Vec<u8>,
    #[prost(bytes, tag="3")]
    pub slash_fraction_conflicting_claim: std::vec::Vec<u8>,
    #[prost(message, repeated, tag="4")]
    pub coins: ::std::vec::Vec<Coin>,
}
/// GenesisState struct
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag="1")]
    pub params: ::std::option::Option<Params>,
    #[prost(message, repeated, tag="2")]
    pub attestations: ::std::vec::Vec<Attestation>,
}
