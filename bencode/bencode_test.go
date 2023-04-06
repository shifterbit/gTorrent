package bencode_test

import (
	. "gtorrent/bencode"
	"testing"
)

func TestParseInt(t *testing.T) {
	type test struct {
		input string
		want  BencodeInt
	}

	tests := []test{
		{input: "i129e", want: BencodeInt(129)},
		{input: "i23e", want: BencodeInt(23)},
		{input: "i0e", want: BencodeInt(0)},
	}
	for _, test := range tests {
		got, err := ParseInt(test.input)
		if err != nil {
			t.Error(err)
		}
		if got != test.want {
			t.Errorf("got %v, %v wanted %v, %v", got, err, test.want, nil)
		}
	}
}

func TestParseIntLeadingZero(t *testing.T) {
	type test struct {
		input string
		want  BencodeInt
		err   error
	}

	tests := []test{
		{input: "i0e", want: BencodeInt(0), err: nil},
		{input: "i023e", want: BencodeInt(0), err: &LeadingZeroError{}},
		{input: "i000e", want: BencodeInt(0), err: &LeadingZeroError{}},
	}

	for _, test := range tests {
		got, err := ParseInt(test.input)
		if got != test.want && err != test.err {
			t.Errorf("got %v, %v wanted %v, %v", got, err, test.want, test.err)
		}

	}
}
