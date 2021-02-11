#![warn(clippy::all)]
#![allow(clippy::pedantic)]

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate serde_json;
#[macro_use]
extern crate log;

pub mod client;
pub mod jsonrpc;
pub mod types;
pub mod utils;
