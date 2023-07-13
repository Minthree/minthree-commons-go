package blockchain

import "time"

type IBlockchain interface {
	GetChain() int
	GetBlockchainOpts() IBlockchainOpts
	VerifySignedMessage(message, signature, address string) (bool, error)

	//  given an address and an assetID, returns whether the address holds the asset and the amount of the asset
	HoldsNFT(address, assetID string) (holds bool, balance int, err error)
}

type IBlockchainOpts interface {
	GetHTTPTimeout() time.Duration
}
