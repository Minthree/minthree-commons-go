package blockchain

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/provider"
	"github.com/klever-io/klever-go-sdk/provider/network"
	"github.com/klever-io/klever-go-sdk/provider/utils"
	"github.com/klever-io/minthree-commons-go/common/chains"
)

type KleverChain struct {
	provider       provider.KleverChain
	blockchainOpts IBlockchainOpts
	chain          int
}

func NewKleverChain(chain int, opts IBlockchainOpts) (*KleverChain, error) {
	var net network.NetworkConfig
	if chain == chains.KLV {
		net = network.NewNetworkConfig(network.MainNet)
	} else {
		net = network.NewNetworkConfig(network.TestNet)
	}

	httpClient := utils.NewHttpClient(opts.GetHTTPTimeout())
	kc, err := provider.NewKleverChain(net, httpClient)
	if err != nil {
		return nil, err
	}

	return &KleverChain{
		provider:       kc,
		blockchainOpts: opts,
		chain:          chain,
	}, nil
}

func (kc KleverChain) GetChain() int {
	return kc.chain
}

func (kc KleverChain) GetBlockchainOpts() IBlockchainOpts {
	return kc.blockchainOpts
}

func (kc KleverChain) VerifySignedMessage(message, signature, addr string) (bool, error) {
	data, err := base64.RawStdEncoding.DecodeString(signature)
	if err != nil {
		data, err = base64.StdEncoding.DecodeString(signature)
		if err != nil {
			return false, err
		}
	}

	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return false, err
	}

	sm := models.NewSM(addr, []byte(message))
	sm.SetSignature(decoded)

	var valid bool
	defer func() {
		if r := recover(); r != nil {
			valid = false
		}
	}()

	valid = sm.Verify()

	return valid, nil
}

func (kc KleverChain) HoldsNFT(address, assetID string) (bool, int, error) {
	account, err := kc.provider.GetAccount(address)
	if err != nil {
		return false, 0, err
	}

	asset := account.Assets[assetID]
	if asset == nil {
		return false, 0, nil
	}

	balance := asset.Balance

	return balance > 0, int(balance), nil
}
