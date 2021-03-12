package goa

import (
	"testing"
)

func TestCheckHeader(t *testing.T) {
	var tb = []struct {
		val    []byte
		expect error
	}{
		{[]byte{}, ErrEmptyByte},
		{[]byte{0x81, 0x80 | byte(32), 0x7d}, nil},
		{[]byte{0x81, 0x8, 0x7d}, ErrMaskInvalid},
		{[]byte{0x80, 0x80 | byte(32), 0x7d}, ErrModeInvalid},
	}

	for _, ts := range tb {
		err := checkHeader(ts.val)
		if err != ts.expect {
			t.Errorf("error got %v, expect %v", err, ts.expect)
		}
	}

}
