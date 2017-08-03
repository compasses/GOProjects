package anchor_test

import (
	"fmt"
	"strconv"
	"testing"
)

func TestAnchorETH_StrToInt(t *testing.T) {
	s1 := "0x15ac8b"

	h, e := strconv.ParseInt(s1, 0, 0)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("h ", h)
}
