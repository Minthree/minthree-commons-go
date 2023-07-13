package blockchain

import (
	"github.com/Minthree/minthree-commons-go/common/chains"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockchainOptsGetHTTPTimeout(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		opts BlockchainOpts
		want time.Duration
	}{
		{"should return the HTTP timeout", BlockchainOpts{HTTPTimeout: 15 * time.Second}, 15 * time.Second},
		{"should return 0 since it is not set", BlockchainOpts{}, 0 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.want, tt.opts.GetHTTPTimeout())
		})
	}
}

func TestNewBlockchain(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name        string
		chain       int
		opts        *BlockchainOpts
		wantedChain int
	}{
		{"should return a KleverChain with chain 33 (TKLV) as default", 0, nil, chains.TKLV},
		{"should return a KleverChain with chain 38 (KLV)", chains.KLV, nil, chains.KLV},
		{"should return a KleverChain with chain 33 (TKLV)", chains.KLV, nil, chains.KLV},
		{"should return a Ethereum with chain 3 (ETH)", chains.EVM, nil, chains.EVM},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got IBlockchain
			var err error
			if tt.opts == nil {
				got, err = NewBlockchain(tt.chain, nil)
			} else {
				got, err = NewBlockchain(tt.chain, tt.opts)
			}

			assert.NoError(err)
			assert.Equal(tt.wantedChain, got.GetChain(), "test %s failed: wanted chain %d, got %d", tt.name, tt.wantedChain, got.GetChain())
		})
	}
}

func TestVerifySignedMessage(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		name      string
		chain     int
		message   string
		signature string
		address   string
		want      bool
	}{
		{"should return true for KLV (33)", chains.KLV, "hello world", "0x123456789abcdef", "my-klv-address", true},
		{"should return false for KLV (33)", chains.KLV, "hello world", "0x123456789abcdeA", "my-klv-address", false},
		{"should return true for ETH (3)", chains.EVM, "0xfaf6c184972e31ee5728839d0be7900c744ac4562ab946b7a1f147728a0afcd8", "0xbb0d378207c1dd182405b61740b3f22e0ed89763f38e0a8cf506e012465ebb012ce07f0bf9a2255a72232aba72e9339f7a0929ad687f05a6faa5afc1b659959f1c", "0x9DAcbE9bD5D177e92DC1E0f380f1A53912d7F7f4", true},
		{"should return false for invalid address ETH (3)", chains.EVM, "0x329932e925e8a98eebe11d80a0d1fc77b1e08b085afcb9660ccddba8f026896a", "0xda7d15442ebac39068b61b3a85cef73ae91313cdcb7b2b3db523b79677fcf7470307049dcc0e209dbcc60ce5dea6ac1883fb029298df2a2d302ab3c8d5d893a31c", "0x70997970C51812dc3A010C7d01b50e0d17dc79C1", false},
		{"should return false for invalid signature ETH (3)", chains.EVM, "0x329932e925e8a98eebe11d80a0d1fc77b1e08b085afcb9660ccddba8f026896a", "0xda7d15442ebac39068b61b3a85cef73ae91313cdcb7b2b3db523b79677fcf7470307049dcc0e209dbcc60ce5dea6ac1883fb029298df2a2d302ab3c8d5d893a31a", "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", false},
		{"should return false for invalid message ETH (3)", chains.EVM, "0x329932e925e8a98eebe11d80a0d1fc77b1e08b085afcb9660ccddba8f026896b", "0xda7d15442ebac39068b61b3a85cef73ae91313cdcb7b2b3db523b79677fcf7470307049dcc0e209dbcc60ce5dea6ac1883fb029298df2a2d302ab3c8d5d893a31c", "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", false},
		{"should return true for valid messa ge ETH (3)", chains.EVM, "Sign in to Minthree", "0xbb711ef43c21187edb66dda6353163dd5bf829b6a66f6cf349725bae2abeeadb455bf8871206d6a02b2cf928eb877217ab4df73a21df6ab7e5285caf1d11034d1c", "0x7de33fb9d06df18b000c8af247c906ef7e5251f0", true},
		{"should return true for TKLV (38)", chains.TKLV, "hello world", "0x123456789abcdef", "my-tklv-address", true},
		{"should return false for TKLV (38)", chains.TKLV, "hello world", "0x123456789abcdeA", "my-tklv-address", false},
		{"should return true for test", chains.KLV, "{\"message\":\"minthree-authorization\"}", "ZjkxZjkxZjNmNGE4MDY3ZjJjYzllNzNlN2Q0MjFkNjc2MzJjOTA5NmM0MjMzYzNkZDIzY2E1NTRlOWRmYjgwYjEwZGUxYjYyODU2YWIwM2IxNWRjYzQ3OTJiMDEzNjYxMTJlMTkyZTMwMDk2ZmEwYmEzZGQ4ZmYxNmFhNDFjMDc=", "klv1szjj2fyfxjtgxw040umd7x37e3e3fj0mpsd5gsfq4vsdj809u9qsqg9w9r", true},
		{"should return true for test", chains.KLV, "{\"message\":\"minthree-authorization\"}", "MzE3ZmZmMTlkYmJhYmE3ZGZkOTA1NDhkMTdiYWRmN2U5MWYwOWU5NjBlYzVhNDRiZWYyMjYxMmNlNGUyYmI4ODYzMjdmNGIwNmRhMTYyZDdjNzFhZDQ4ZDY4ZmUzNzJlM2E1MGEzNjAwY2RjZjc5NDRjMGQ0MzlkMTQ5NmI5MGY", "klv1p5s05q9rzmf9c8n8n5glsufxvwqdh7glw25y30hnsvjk3f80jyyq85437l", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockchain, err := NewBlockchain(tt.chain)
			require.NoError(err, "test %s failed: could not create blockchain with error %v", tt.name, err)

			got, err := blockchain.VerifySignedMessage(tt.message, tt.signature, tt.address)
			require.NoError(err, "test %s failed: could not verify signed message with error %v", tt.name, err)
			assert.Equal(tt.want, got, "test %s failed: wanted %t, got %t", tt.name, tt.want, got)
		})
	}
}
