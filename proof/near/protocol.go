package near

import (
	"github.com/mapprotocol/near-api-go/pkg/types"
)

var (
	NewFunctionCallGas types.Gas = 30 * 10000000000000
	Deposit                      = "0.3"
)

type Result struct {
	BlockHash   string        `json:"block_hash"`
	BlockHeight int           `json:"block_height"`
	Logs        []interface{} `json:"logs"`
	Result      []byte        `json:"result"`
}

type TransferOut struct {
	FromChain string `json:"from_chain"`
	ToChain   string `json:"to_chain"`
	OrderId   string `json:"order_id"`
}
