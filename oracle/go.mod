module github.com/MinterTeam/minter-hub-oracle

go 1.13

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace github.com/MinterTeam/mhub/chain => /Users/daniillashin/Desktop/mhub/chain

require (
	github.com/MinterTeam/mhub/chain v0.0.0-20210211115446-c4ccc6c0254d
	github.com/MinterTeam/minter-go-sdk/v2 v2.1.0-rc2.0.20210209133819-011976d40e49
	github.com/cosmos/cosmos-sdk v0.40.1
	github.com/cosmos/go-bip39 v1.0.0
	github.com/tendermint/tendermint v0.34.3
	github.com/valyala/fasthttp v1.19.0
	google.golang.org/grpc v1.35.0
)
