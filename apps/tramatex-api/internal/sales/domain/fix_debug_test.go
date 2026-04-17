package domain

import (
	"testing"
	"fmt"
)

func TestValidateReadyStatus(t *testing.T) {
	val := SalesOrderStatus("READY_FOR_PRODUCTION")
	err := val.IsValid()
	if err != nil {
		t.Errorf("Validation failed for READY_FOR_PRODUCTION: %v", err)
	} else {
		fmt.Println("Validation SUCCEEDED for READY_FOR_PRODUCTION")
	}

	fmt.Printf("Constant value: %q\n", string(SalesOrderStatusReadyForProduction))
	if SalesOrderStatusReadyForProduction != val {
		t.Errorf("Constant mismatch: %q != %q", string(SalesOrderStatusReadyForProduction), string(val))
	}
}
