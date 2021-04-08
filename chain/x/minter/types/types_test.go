package types

import (
	"bytes"
	mrand "math/rand"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestValsetSort(t *testing.T) {
	specs := map[string]struct {
		src BridgeValidators
		exp BridgeValidators
	}{
		"by power desc": {
			src: BridgeValidators{
				{Power: 1, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()},
				{Power: 2, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()},
				{Power: 3, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()},
			},
			exp: BridgeValidators{
				{Power: 3, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()},
				{Power: 2, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()},
				{Power: 1, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()},
			},
		},
		"by eth addr on same power": {
			src: BridgeValidators{
				{Power: 1, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()},
				{Power: 1, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()},
				{Power: 1, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()},
			},
			exp: BridgeValidators{
				{Power: 1, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()},
				{Power: 1, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()},
				{Power: 1, MinterAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()},
			},
		},
		// if you're thinking about changing this due to a change in the sorting algorithm
		// you MUST go change this in peggy_utils/types.rs as well. You will also break all
		// bridges in production when they try to migrate so use extreme caution!
		"real world": {
			src: BridgeValidators{
				{Power: 678509841, MinterAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, MinterAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 685294939, MinterAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 671724742, MinterAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, MinterAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 617443955, MinterAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 6785098, MinterAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
				{Power: 291759231, MinterAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			exp: BridgeValidators{
				{Power: 685294939, MinterAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 678509841, MinterAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, MinterAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, MinterAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 671724742, MinterAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 617443955, MinterAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 291759231, MinterAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
				{Power: 6785098, MinterAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
			},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			spec.src.Sort()
			assert.Equal(t, spec.src, spec.exp)
			shuffled := shuffled(spec.src)
			shuffled.Sort()
			assert.Equal(t, shuffled, spec.exp)
		})
	}
}

func shuffled(v BridgeValidators) BridgeValidators {
	mrand.Shuffle(len(v), func(i, j int) {
		v[i], v[j] = v[j], v[i]
	})
	return v
}

func TestBridgeValidators_PowerDiff(t *testing.T) {
	type args struct {
		c BridgeValidators
	}
	tests := []struct {
		name string
		b    BridgeValidators
		args args
		want float64
	}{
		{
			name: "case 1",
			b:    BridgeValidators{
				{
					Power:         4286394505+8572789,
					MinterAddress: "Mx000",
				},
			},
			args: args{
				c: BridgeValidators{
					{
						Power:         4286394505,
						MinterAddress: "Mx000",
					},
					{
						Power:         8572789,
						MinterAddress: "Mx001",
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.PowerDiff(tt.args.c); got != tt.want {
				t.Errorf("PowerDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}