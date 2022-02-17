package utils

import (
	"fmt"
	"testing"
)

func TestRandNum(t *testing.T) {
	for i := 0; i < 100; i ++ {
		fmt.Println(RandNum())
	}
}
