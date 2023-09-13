package constant

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ChainId uint64
type Role string
type BlockIdOfEth2 string

const (
	FinalBlockIdOfEth2   BlockIdOfEth2 = "finalized"
	HeadBlockIdOfEth2    BlockIdOfEth2 = "head"
	GenesisBlockIdOfEth2 BlockIdOfEth2 = "genesis"
)

const (
	MethodOfGetBytes           = "getBytes"
	MethodOfGetFinalBytes      = "getFinalBytes"
	MethodOfVerifyReceiptProof = "verifyReceiptProof"
)

const (
	EpochOfMap          = 50000
	EpochOfBsc          = 200
	HeaderCountOfBsc    = 12
	HeaderCountOfPlaton = 430
	EpochOfKlaytn       = 3600
	HeaderOneCount      = 1
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
