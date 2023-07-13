package blockchain

import (
	"github.com/Minthree/minthree-commons-go/common/chains"
	"github.com/Minthree/minthree-commons-go/common/constants"
)

func NewBlockchain(chain int, opts ...IBlockchainOpts) (IBlockchain, error) {
	if len(opts) == 0 || opts[0] == nil {
		opts = make([]IBlockchainOpts, 1)
		opts[0] = &BlockchainOpts{
			HTTPTimeout: constants.DefaultHTTPTimeout,
		}
	}

	var blockchain IBlockchain
	var err error

	switch chain {
	case chains.KLV, chains.TKLV:
		blockchain, err = NewKleverChain(chain, opts[0])

	case chains.EVM:
		blockchain, err = NewEvm(chain, opts[0])

	default:
		blockchain, err = NewKleverChain(chains.TKLV, opts[0])
	}

	return blockchain, err
}
