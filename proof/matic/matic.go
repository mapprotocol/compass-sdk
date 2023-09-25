package matic

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/mapprotocol/compass-sdk/constant"
	"github.com/mapprotocol/compass-sdk/mapprotocol"
	"github.com/mapprotocol/compass-sdk/pkg/ethclient"
	"github.com/mapprotocol/compass-sdk/proof"
	"github.com/mapprotocol/compass-sdk/utils"
	"math/big"
)

type BlockHeader struct {
	ParentHash       []byte         `json:"parentHash"`
	Sha3Uncles       []byte         `json:"sha3Uncles"`
	Miner            common.Address `json:"miner"`
	StateRoot        []byte         `json:"stateRoot"`
	TransactionsRoot []byte         `json:"transactionsRoot"`
	ReceiptsRoot     []byte         `json:"receiptsRoot"`
	LogsBloom        []byte         `json:"logsBloom"`
	Difficulty       *big.Int       `json:"difficulty"`
	Number           *big.Int       `json:"number"`
	GasLimit         *big.Int       `json:"gasLimit"`
	GasUsed          *big.Int       `json:"gasUsed"`
	Timestamp        *big.Int       `json:"timestamp"`
	ExtraData        []byte         `json:"extraData"`
	MixHash          []byte         `json:"mixHash"`
	Nonce            []byte         `json:"nonce"`
	BaseFeePerGas    *big.Int       `json:"baseFeePerGas"`
}

func ConvertHeader(header *types.Header) BlockHeader {
	bloom := make([]byte, 0, len(header.Bloom))
	for _, b := range header.Bloom {
		bloom = append(bloom, b)
	}
	nonce := make([]byte, 0, len(header.Nonce))
	for _, b := range header.Nonce {
		nonce = append(nonce, b)
	}
	return BlockHeader{
		ParentHash:       hashToByte(header.ParentHash),
		Sha3Uncles:       hashToByte(header.UncleHash),
		Miner:            constant.ZeroAddress,
		StateRoot:        hashToByte(header.Root),
		TransactionsRoot: hashToByte(header.TxHash),
		ReceiptsRoot:     hashToByte(header.ReceiptHash),
		LogsBloom:        bloom,
		Difficulty:       header.Difficulty,
		Number:           header.Number,
		GasLimit:         new(big.Int).SetUint64(header.GasLimit),
		GasUsed:          new(big.Int).SetUint64(header.GasUsed),
		Timestamp:        new(big.Int).SetUint64(header.Time),
		ExtraData:        header.Extra,
		MixHash:          hashToByte(header.MixDigest),
		Nonce:            nonce,
		BaseFeePerGas:    header.BaseFee,
	}
}

func hashToByte(h common.Hash) []byte {
	ret := make([]byte, 0, len(h))
	for _, b := range h {
		ret = append(ret, b)
	}
	return ret
}

type ProofData struct {
	Headers      []BlockHeader
	ReceiptProof ReceiptProof
}

type ReceiptProof struct {
	TxReceipt mapprotocol.TxReceipt
	KeyIndex  []byte
	Proof     [][]byte
}

func GetProof(client *ethclient.Client, latestBlock *big.Int, log *types.Log, method string, fId constant.ChainId) ([]byte, error) {
	txsHash, err := mapprotocol.GetTxsHashByBlockNumber(client, latestBlock)
	if err != nil {
		return nil, fmt.Errorf("unable to get tx hashes Logs: %w", err)
	}
	receipts, err := mapprotocol.GetReceiptsByTxsHash(client, txsHash)
	if err != nil {
		return nil, fmt.Errorf("unable to get receipts hashes Logs: %w", err)
	}

	headers := make([]*types.Header, constant.ConfirmsOfMatic.Int64())
	for i := 0; i < int(constant.ConfirmsOfMatic.Int64()); i++ {
		headerHeight := new(big.Int).Add(latestBlock, new(big.Int).SetInt64(int64(i)))
		tmp, err := client.HeaderByNumber(context.Background(), headerHeight)
		if err != nil {
			return nil, fmt.Errorf("getHeader failed, err is %v", err)
		}
		headers[i] = tmp
	}

	mHeaders := make([]BlockHeader, 0, len(headers))
	for _, h := range headers {
		mHeaders = append(mHeaders, ConvertHeader(h))
	}
	return AssembleProof(mHeaders, *log, fId, receipts, method)
}

func AssembleProof(headers []BlockHeader, log types.Log, fId constant.ChainId, receipts []*types.Receipt, method string) ([]byte, error) {
	txIndex := log.TxIndex
	receipt, err := mapprotocol.GetTxReceipt(receipts[txIndex])
	if err != nil {
		return nil, err
	}

	prf, err := proof.GetByReceipt(receipts, txIndex)
	if err != nil {
		return nil, err
	}

	var key []byte
	key = rlp.AppendUint64(key[:0], uint64(txIndex))
	ek := utils.Key2Hex(key, len(prf))

	pd := ProofData{
		Headers: headers,
		ReceiptProof: ReceiptProof{
			TxReceipt: *receipt,
			KeyIndex:  ek,
			Proof:     prf,
		},
	}

	input, err := constant.Matic.Methods[constant.MethodOfGetBytes].Inputs.Pack(pd)
	if err != nil {
		return nil, err
	}

	pack, err := utils.PackInput(constant.Mcs, method, new(big.Int).SetUint64(uint64(fId)), input)
	if err != nil {
		return nil, err
	}

	return pack, nil
}
