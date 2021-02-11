use actix_web::client::SendRequestError as ActixError;
use std::error::Error;
use std::fmt::Display;
use std::fmt::Formatter;
use std::fmt::Result;

#[derive(Debug)]
pub enum JsonRpcError {
    NoToken,
    BadResponse(String),
    BadStruct(String),
    FailedToSend(ActixError),
    ResponseError {
        code: i64,
        message: String,
        data: String,
    },
    BadInput(String),
}

impl Display for JsonRpcError {
    fn fmt(&self, f: &mut Formatter) -> Result {
        match self {
            JsonRpcError::NoToken => {
                write!(f, "Account has no tokens! No details!")
            }
            JsonRpcError::BadResponse(val) => write!(f, "JsonRPC bad response {}", val),
            JsonRpcError::BadStruct(val) => write!(f, "JsonRPC unexpected json returned {}", val),
            JsonRpcError::BadInput(val) => write!(f, "JsonRPC bad input {}", val),
            JsonRpcError::FailedToSend(val) => write!(f, "JsonRPC Failed to send {}", val),
            JsonRpcError::ResponseError {
                code,
                message,
                data,
            } => write!(
                f,
                "JsonRPC Response error code {} message {} data {:?}",
                code, message, data
            ),
        }
    }
}

impl Error for JsonRpcError {}
