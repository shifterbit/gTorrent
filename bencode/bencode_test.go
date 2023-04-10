package bencode_test

import (
	. "gtorrent/bencode"
	"reflect"
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
		if *got != test.want {
			t.Errorf("got %v, %v wanted %v, %v", got, err, test.want, nil)
		}
	}
}

func TestParseIntLeadingZero(t *testing.T) {
	type test struct {
		input string
		want  *BencodeInt
		err   error
	}

	zero := BencodeInt(0)
	tests := []test{
		{input: "i0e", want: &zero, err: nil},
		{input: "i023e", want: nil, err: &LeadingZeroError{}},
		{input: "i000e", want: nil, err: &LeadingZeroError{}},
	}

	for _, test := range tests {
		got, err := ParseInt(test.input)
		if got != test.want && err != test.err {
			t.Errorf("got %v, %v wanted %v, %v", got, err, test.want, test.err)
		}
	}
}

func TestParseString(t *testing.T) {
	type test struct {
		input string
		want  BencodeString
	}

	tests := []test{
		{input: "3:foo", want: BencodeString{String: "foo", Length: 3}},
		{input: "4:spam", want: BencodeString{String: "spam", Length: 4}},
		{input: "6:foobar", want: BencodeString{String: "foobar", Length: 6}},
		{input: "0:", want: BencodeString{String: "", Length: 0}},
	}
	for _, test := range tests {
		got, err := ParseString(test.input)
		if err != nil {
			t.Error(err)
		}
		if *got != test.want {
			t.Errorf("got %q, %v wanted %q, %v", got, err, test.want, nil)
		}
	}
}

func TestParseList(t *testing.T) {
	type test struct {
		input string
		want  any
	}

	tests := []test{
		{input: "li1ei2ei3ee", want: []any{1, 2, 3}},
		{input: "l3:ham6:foobare", want: []any{"ham", "foobar"}},
		{input: "lli2ei4eeli6ei8eee", want: []any{[]any{2, 4}, []any{6, 8}}},
		{input: "ll3:foo3:barel3:egg3:hamee", want: []any{[]any{"foo", "bar"}, []any{"egg", "ham"}}},
		{input: "l3:fool3:foo3:barei25ee", want: []any{"foo", []any{"foo", "bar"}, 25}},
	}

	for _, test := range tests {
		got, err := ParseList(test.input)
		if err != nil {
			t.Error(err)
		}
		if reflect.DeepEqual(got.Value(), test.want) == false {
			t.Errorf("got %q, %v wanted %q, %v", got.Value(), err, test.want, nil)
		}
	}

}

func TestParseDict(t *testing.T) {
	type test struct {
		input string
		want  any
	}

	tests := []test{
		{input: "d3:foo3:bar3:egg3:hame",
			want: map[string]any{"foo": "bar", "egg": "ham"}},
		{input: "d3:fooli2ei4ei6ee3:bard3:egg3:hamee",
			want: map[string]any{"foo": []any{2,4,6}, "bar": map[string]any{"egg":"ham"} }},
	}

	for _, test := range tests {
		got, err := ParseDict(test.input)
		if err != nil {
			t.Error(err)
		}
		if reflect.DeepEqual(got.Value(), test.want) == false {
			t.Errorf("got %q, %v wanted %q, %v", got.Value(), err, test.want, nil)
		}
	}

}
