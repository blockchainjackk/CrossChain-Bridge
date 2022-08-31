package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDecodeMemo(t *testing.T) {

	memo := "5357415054583a307865393436663430653364343735653063373164623835316164383236333432363765623465613930613530613164363635343137663233366139396466323638"
	decodeString, _ := hex.DecodeString(memo)
	fmt.Println(string(decodeString))
}
