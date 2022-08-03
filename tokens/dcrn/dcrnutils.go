package dcrn

import (
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrutil/v3"
)

// Inheritable interface
type Inheritable interface {
	GetChainParams() *chaincfg.Params
}

type dcrnAmountType dcrutil.Amount