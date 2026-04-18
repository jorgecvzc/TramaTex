package domain

import (
	"testing"
	"fmt"
)

func TestDebugStatus(t *testing.T) {
	s := SalesOrderStatusReadyForProduction
	err := s.IsValid()
	if err != nil {
		t.Errorf("SalesOrderStatusReadyForProduction.IsValid() failed: %v", err)
		fmt.Printf("Hex of SalesOrderStatusReadyForProduction: %x\n", string(s))
	}

	raw := SalesOrderStatus("READY_FOR_PRODUCTION")
	errRaw := raw.IsValid()
	if errRaw != nil {
		t.Errorf("SalesOrderStatus(\"READY_FOR_PRODUCTION\").IsValid() failed: %v", errRaw)
		fmt.Printf("Hex of raw string: %x\n", string(raw))
	}
	
	if s != raw {
		t.Errorf("Mismatch between constant and raw string: %q != %q", s, raw)
	}
}
