module github.com/MinterTeam/minter-hub-connector

go 1.13

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

// replace github.com/MinterTeam/mhub/chain => /Users/daniillashin/Desktop/mhub/chain

require (
	github.com/MinterTeam/mhub/chain v0.0.0-20210226110917-2ec89dd74ac7
	github.com/MinterTeam/minter-go-sdk/v2 v2.2.0-alpha1.0.20210312100821-af23f9256cc4
	github.com/cosmos/cosmos-sdk v0.41.2
	github.com/cosmos/go-bip39 v1.0.0
	github.com/ethereum/go-ethereum v1.9.22
	github.com/tendermint/tendermint v0.34.7
	google.golang.org/grpc v1.35.0
)
