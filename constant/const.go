package constant

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ChainId uint64
type Role string

const (
	MethodOfGetBytes = "getBytes"
)

var (
	RoleOfMaintainer Role = "maintainer"
	RoleOfMessenger  Role = "messenger"
	RoleOfMonitor    Role = "monitor"
)

var (
	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

var (
	ConfirmsOfMatic       = big.NewInt(10)
	HeaderLengthOfEth2    = 20
	HeaderLengthOfConflux = 20
)
