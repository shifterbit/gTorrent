package bencode

import (
	"errors"
	"strconv"
	"strings"
)

type BencodeValue interface {
	/*
	 Converts a Given `BencodeValue` to a plain Go value
	   - `BencodeString` becomes `string`
	   - `BencodeInt` becomes `int`
	*/
	Value() any
}

type BencodeString struct {
	Length int
	String    string
}

type BencodeInt int

func (s *BencodeString) Value() any {
	return s.String
}

func (i *BencodeInt) Value() any {
	return int(*i)
}

type LeadingZeroError struct{}

func (e *LeadingZeroError) Error() string {
	return "cannot start an integer with a leading zero"
}

// Parses bencoded data and returns a `BencodeValue`
func Parse(text string) (BencodeValue, error) {
	// TODO
	return nil, nil

}

// Parses a bencoded string, returning `BencodeString`
func ParseString(str string) (BencodeString, error) {
	val := strings.SplitN(str, ":", 2)
	length, err := strconv.Atoi(val[0])
	if err != nil {
		return BencodeString{}, err
	}

	return BencodeString{
		Length: length,
		String:    val[1],
	}, nil
}

// Parses a bencoded integer, returning `BencodeInteger`
func ParseInt(str string) (BencodeInt, error) {
	start := str[0]
	end := str[len(str)-1]
	if start != 'i' {
		return 0, errors.New("Bencode: Error Parsing Int")
	}
	if end != 'e' {
		return 0, errors.New("Bencode: Unexpected End of File")
	}

	stringifiedNumber := str[1 : len(str)-1]
	// We need to check for leading zeroes as integers with leading zeroes
	// are considred invalid
	if stringifiedNumber[0] == '0' && len(stringifiedNumber) > 1 {
		return 0, &LeadingZeroError{}
	}

	num, err := strconv.Atoi(stringifiedNumber)

	if err != nil {
		return 0, err
	}

	return BencodeInt(num), nil
}
