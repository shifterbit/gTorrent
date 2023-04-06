package bencode

import (
	"os"
	"strconv"
	"strings"
)

type BencodeType interface {
	Value() any
}

type BencodeString struct {
	length int
	str    string
}

func (s *BencodeString) Value() any {
	return s.str
}

// Read a bencode file
func ParseFile(path string) ([]string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return []string{}, err
	}
	result := string(bytes)
	res := strings.Split(result, ":")

	return res, nil
}

func ParseString(str string) (BencodeString, error) {
	val := strings.SplitN(str, ":", 2)
	length, err := strconv.Atoi(val[0])
	if err != nil {
		return BencodeString{}, err
	}

	return BencodeString{
		length: length,
		str:    val[1],
	}, nil
}
