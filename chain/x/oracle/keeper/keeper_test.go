package keeper

import (
	"fmt"
	"testing"
)

func TestPrefixRange(t *testing.T) {
	p, e := prefixRange([]byte{1, 3, 4})

	fmt.Println(p)
	fmt.Println(e)

	p, e = prefixRange([]byte{15, 42, 255, 255})

	fmt.Println(p)
	fmt.Println(e)

	p, e = prefixRange([]byte{255, 255, 255, 255})

	fmt.Println(p)
	fmt.Println(e)

	p, e = prefixRange([]byte{2})

	fmt.Println(p)
	fmt.Println(e)
}