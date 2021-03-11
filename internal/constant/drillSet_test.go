package constant

import (
	"fmt"
	"testing"
)

func TestGetBoundary(t *testing.T) {
	x, y := GetBoundary()
	if len(x)!= len(y){
		t.Error("error")
	}
	fmt.Println(x)
	fmt.Println(y)
}
