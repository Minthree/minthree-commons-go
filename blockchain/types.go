package blockchain

import "time"

type BlockchainOpts struct {
	HTTPTimeout time.Duration
}

func (opts BlockchainOpts) GetHTTPTimeout() time.Duration {
	return opts.HTTPTimeout
}
