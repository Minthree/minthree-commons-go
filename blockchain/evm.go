package blockchain

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

type Evm struct {
	chain          int
	blockchainOpts IBlockchainOpts
}

func NewEvm(chain int, opts IBlockchainOpts) (*Evm, error) {
	return &Evm{
		chain:          chain,
		blockchainOpts: opts,
	}, nil
}

func (e Evm) GetChain() int {
	return e.chain
}

func (e Evm) GetBlockchainOpts() IBlockchainOpts {
	return e.blockchainOpts
}

func (e Evm) VerifySignedMessage(message, signature, addr string) (bool, error) {
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	sig[64] -= 27

	msgDecoded, err := hexutil.Decode(message)
	if err != nil {
		msgDecoded = []byte(message)
	}

	msgBytes := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msgDecoded), msgDecoded)))

	sigPublicKeyECDSA, err := crypto.SigToPub(msgBytes, sig)
	if err != nil {
		return false, nil
	}

	address := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	if strings.ToLower(address.String()) != strings.ToLower(addr) {
		return false, nil
	}

	sigPublicKey, err := crypto.Ecrecover(msgBytes, sig)
	if err != nil {
		return false, nil
	}

	return crypto.VerifySignature(sigPublicKey, msgBytes, sig[:len(sig)-1]), nil
}

func (e Evm) HoldsNFT(address, assetID string) (bool, int, error) {
	panic("not implemented")
}
