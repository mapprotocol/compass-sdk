package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type Hash string

// ToCommonHash converts hash to common.Hash
func (hash Hash) ToCommonHash() *common.Hash {
	newHash := common.HexToHash(string(hash))
	return &newHash
}

// String implements the interface stringer
func (hash Hash) String() string {
	return string(hash)
}

func (hash *Hash) UnmarshalJSON(input []byte) error {
	var h common.Hash
	if err := json.Unmarshal(input, &h); err != nil {
		return err
	}
	*hash = Hash(h.String())
	return nil
}

// Bloom is a hash type with 256 bytes.
type Bloom string

type NonceType int

const (
	NONCE_TYPE_AUTO NonceType = iota
	NONCE_TYPE_NONCE
	NONCE_TYPE_PENDING_NONCE
)

// NewBigInt creates a big number with specified uint64 value.
func NewBigInt(x uint64) *hexutil.Big {
	n1 := new(big.Int).SetUint64(x)
	n2 := hexutil.Big(*n1)
	return &n2
}

// NewBigIntByRaw creates a hexutil.big with specified big.int value.
func NewBigIntByRaw(x *big.Int) *hexutil.Big {
	if x == nil {
		return nil
	}
	v := hexutil.Big(*x)
	return &v
}

// NewUint64 creates a hexutil.Uint64 with specified uint64 value.
func NewUint64(x uint64) *hexutil.Uint64 {
	n1 := hexutil.Uint64(x)
	return &n1
}

// NewUint creates a hexutil.Uint with specified uint value.
func NewUint(x uint) *hexutil.Uint {
	n1 := hexutil.Uint(x)
	return &n1
}

// NewBytes creates a hexutil.Bytes with specified input value.
func NewBytes(input []byte) hexutil.Bytes {
	return hexutil.Bytes(input)
}
