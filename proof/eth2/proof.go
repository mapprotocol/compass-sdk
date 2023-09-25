package eth2

import (
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

func GenerateByApi(slot []string) [][32]byte {
	ret := make([][32]byte, 0, len(slot))
	for _, op := range slot {
		ret = append(ret, common.HexToHash(op))
	}

	return ret
}

type ReceiptProof struct {
	Header    BlockHeader
	TxReceipt mapprotocol.TxReceipt
	KeyIndex  []byte
	Proof     [][]byte
}

func GetProof(client *ethclient.Client, endPoint string, latestBlock *big.Int, log *types.Log, method string, fId constant.ChainId) ([]byte, error) {
	header, err := client.EthLatestHeaderByNumber(endPoint, latestBlock)
	if err != nil {
		return nil, err
	}
	// when syncToMap we need to assemble a tx proof
	txsHash, err := mapprotocol.GetTxsHashByBlockNumber(client, latestBlock)
	if err != nil {
		return nil, fmt.Errorf("unable to get tx hashes Logs: %w", err)
	}
	receipts, err := mapprotocol.GetReceiptsByTxsHash(client, txsHash)
	if err != nil {
		return nil, fmt.Errorf("unable to get receipts hashes Logs: %w", err)
	}
	return AssembleProof(*ConvertHeader(header), *log, receipts, method, fId)
}

func AssembleProof(header BlockHeader, log types.Log, receipts []*types.Receipt, method string, fId constant.ChainId) ([]byte, error) {
	txIndex := log.TxIndex
	receipt, err := mapprotocol.GetTxReceipt(receipts[txIndex])
	if err != nil {
		return nil, err
	}

	prf, err := proof.Get(receipts, txIndex)
	if err != nil {
		return nil, err
	}

	var key []byte
	key = rlp.AppendUint64(key[:0], uint64(txIndex))
	ek := utils.Key2Hex(key, len(prf))

	pd := ReceiptProof{
		Header:    header,
		TxReceipt: *receipt,
		KeyIndex:  ek,
		Proof:     prf,
	}

	input, err := constant.Eth2.Methods[constant.MethodOfGetBytes].Inputs.Pack(pd)
	if err != nil {
		return nil, err
	}

	pack, err := utils.PackInput(constant.Mcs, method, new(big.Int).SetUint64(uint64(fId)), input)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
