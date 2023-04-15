package bencode

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type BencodeText struct {
	Length int
	String string
}

type BencodeValue interface {
	/*
	 Converts a Given `BencodeValue` to a plain Go value
	   - `BencodeString` becomes `string`
	   - `BencodeInt` becomes `int`
	   - `BencodeList` becomes []any
	   - `BencodeDict` becomes map[string]any
	*/
	Value() any
}

type BencodeString struct {
	Length int
	String string
}

type BencodeInt int

type BencodeList []BencodeValue

type BencodeDict map[string]BencodeValue

func (s BencodeString) Value() any {
	return s.String
}

func (i BencodeInt) Value() any {
	return int(i)
}

func (l BencodeList) Value() any {
	var list []any
	for _, item := range l {
		list = append(list, item.Value())
	}
	return list
}

func (d BencodeDict) Value() any {
	dict := make(map[string]any)
	for k, v := range d {
		dict[k] = v.Value()
	}
	return dict
}

type LeadingZeroError struct{}

func (e *LeadingZeroError) Error() string {
	return "cannot start an integer with a leading zero"
}

type IncorrectStringLengthError struct {
	String         string
	ExpectedLength int
	ActualLength   int
}

func (e *IncorrectStringLengthError) Error() string {
	return fmt.Sprintf("unexpected string length for %q, got %d, expected %d", e.String, e.ActualLength, e.ExpectedLength)
}

func Marshall(v any) ([]byte, error) {
	res, err := generateVal(v)
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

func Unmarshall(data []byte, v any) error {
	bencodeVal, err := Parse(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v).Elem()
	switch rv.Kind() {
	case reflect.Struct:
		var dict = bencodeVal.Value()
		decoderConfig := mapstructure.DecoderConfig{TagName: "bencode", Result: v}
		decoder, err := mapstructure.NewDecoder(&decoderConfig)
		if err != nil {
			return err
		}
		decoder.Decode(dict)
	default:
		val := bencodeVal.Value()
		rVal := reflect.ValueOf(val)
		reflect.ValueOf(v).Elem().Set(rVal)

	}

	return nil
}

// Parses bencoded data and returns a `BencodeValue`
func Parse(str string) (BencodeValue, error) {
	start := str[0]

	var res BencodeValue
	var err error = nil
	switch {
	case isDigit(string(start)):
		res, err = ParseString(str)
	case start == 'i':
		res, err = ParseInt(str)
	case start == 'l':
		res, err = ParseList(str)
	case start == 'd':
		res, err = ParseDict(str)
	default:
		res, err = nil, errors.New("Invalid bencode")
	}

	return res, err

}

// Parses a bencoded string, returning `BencodeString`
func ParseString(str string) (*BencodeString, error) {
	val := strings.SplitN(str, ":", 2)

	length, err := strconv.Atoi(val[0])
	if err != nil {
		return nil, err
	}

	if length != len(val[1]) {
		return nil,
			&IncorrectStringLengthError{
				String:         val[1],
				ActualLength:   len(val[1]),
				ExpectedLength: length,
			}
	}

	return &BencodeString{
		Length: length,
		String: val[1],
	}, nil
}

// Parses a bencoded integer, returning `BencodeInteger`
func ParseInt(str string) (*BencodeInt, error) {
	end := str[len(str)-1]
	if end != 'e' {
		return nil, errors.New("bencode: Unexpected End of File")
	}

	stringifiedNumber := str[1 : len(str)-1]
	// We need to check for leading zeroes as integers with leading zeroes
	// are considred invalid
	if stringifiedNumber[0] == '0' && len(stringifiedNumber) > 1 {
		return nil, &LeadingZeroError{}
	}

	num, err := strconv.Atoi(stringifiedNumber)

	if err != nil {
		return nil, err
	}

	result := BencodeInt(num)

	return &result, nil
}

// Parse a bencoded list, returning a list of `BencodedValue`
func ParseList(str string) (*BencodeList, error) {
	str = str[1:]
	end := str[len(str)-1]
	if end != 'e' {
		return nil, errors.New("bencode: Unexpected End of File")
	}
	var list BencodeList
	for len(str) > 1 {
		switch {
		case isDigit(string(str[0])):
			text := readString(str)
			val, err := ParseString(str[:text.Length])
			if err != nil {
				return nil, err
			}
			list = append(list, val)
			str = str[text.Length:]
		case str[0] == 'i':
			text := readInt(str)
			val, err := ParseInt(str[:text.Length])
			if err != nil {
				return nil, err
			}
			list = append(list, val)
			str = str[text.Length:]
		case str[0] == 'l':
			text := readList(str[1:])
			val, err := ParseList(str[:text.Length])
			if err != nil {
				return nil, err
			}
			list = append(list, val)
			str = str[text.Length:]
		case str[0] == 'd':
			text := readDict(str[1:])
			val, err := ParseDict(str[:text.Length])
			if err != nil {
				return nil, err
			}
			list = append(list, val)
			str = str[text.Length:]
		case str[0] == 'e':
			str = str[1:]
		}
	}

	return &list, nil
}

func ParseDict(str string) (*BencodeDict, error) {
	str = str[1:]
	end := str[len(str)-1]
	if end != 'e' {
		return nil, errors.New("bencode: Unexpected End of File")
	}
	var list BencodeList
	dict := make(BencodeDict)
	for len(str) > 1 {
		switch {
		case isDigit(string(str[0])):
			text := readString(str)
			val, err := ParseString(str[:text.Length])
			if err != nil {
				return nil, err
			}
			list = append(list, val)
			str = str[text.Length:]
		case str[0] == 'i':
			text := readInt(str)
			val, err := ParseInt(str[:text.Length])
			if err != nil {
				return nil, err
			}
			list = append(list, val)
			str = str[text.Length:]
		case str[0] == 'l':
			text := readList(str[1:])
			val, err := ParseList(str[:text.Length])
			if err != nil {
				return nil, err
			}
			list = append(list, val)
			str = str[text.Length:]
		case str[0] == 'd':
			text := readDict(str[1:])
			val, err := ParseDict(str[:text.Length])
			if err != nil {
				return nil, err
			}
			list = append(list, val)
			str = str[text.Length:]
		case str[0] == 'e':
			str = str[1:]
		}
	}

	if isEven(len(list)) == false {
		return nil, errors.New("bencode: missing entry in key value pairs")
	}

	for len(list) > 0 {
		dict[list[0].Value().(string)] = list[1]
		list = list[2:]
	}

	return &dict, nil
}

func isDigit(str string) bool {
	start := string(str[0])
	isDigit := regexp.MustCompile(`\d`)
	return isDigit.MatchString(start)
}

func isEven(num int) bool {
	return num%2 == 0
}

func mapKeys(dict map[string]any) []string {
	keys := make([]string, 0, 10)
	for k := range dict {
		keys = append(keys, k)
	}
	return keys
}

func generateVal(v any) (string, error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return generateInt(v.(int)), nil
	case reflect.String:
		return generateString(v.(string)), nil
	case reflect.Slice, reflect.Array:
		res := "l"
		for _, v := range v.([]any)[:] {
			val, err := generateVal(v)
			if err != nil {
				return "", err
			}
			res += val
		}
		res += "e"
		return res, nil
	case reflect.Map:
		res := "d"
		v := rv.Interface()
		keys := mapKeys(v.(map[string]any))
		sort.Strings(keys)
		for k, v := range v.(map[string]any) {
			res += generateString(k)
			val, err := generateVal(v)
			if err != nil {
				return "", err
			}
			res += val
		}
		res += "e"
		return res, nil
	case reflect.Struct:
		dict := new(map[string]any)
		decoderConfig := mapstructure.DecoderConfig{TagName: "bencode", Result: dict}
		decoder, err := mapstructure.NewDecoder(&decoderConfig)
		if err != nil {
			return "", err
		}
		decoder.Decode(v)
		return generateVal(&dict)
	}
	return "", nil
}

func generateString(str string) string {
	length := len(str)
	num := strconv.Itoa(length)
	res := num + ":" + str
	return res
}

func generateInt[I Integer](num I) string {
	number := strconv.Itoa(int(num))
	res := "i" + number + "e"
	return res
}

func readString(str string) BencodeText {
	length := 0
	digits := ""
	loopComplete := false
	for _, c := range str {
		if loopComplete == true {
			break
		}

		switch {
		case isDigit(string(c)):
			digits = digits + string(c)
			length = length + 1
		case c == ':':
			length = length + 1
			loopComplete = true
			break
		}
	}

	stringLength, _ := strconv.Atoi(string(digits))
	length = length + stringLength

	return BencodeText{
		Length: length,
		String: str[:length],
	}
}

func readInt(str string) BencodeText {
	digits := ""
	length := 0
	loopComplete := false
	for _, c := range str {
		if loopComplete == true {
			break
		}
		switch {
		case c == 'i':
			length = length + 1
		case isDigit(string(c)):
			digits = digits + string(c)
			length = length + 1
		case c == 'e':
			loopComplete = true
			length = length + 1
		}
	}

	return BencodeText{
		String: str[:length],
		Length: length,
	}
}

func readList(str string) BencodeText {
	length := 1
	text := "l"
	loopComplete := false
	for len(str) > 1 {
		if loopComplete == true {
			break
		}
		switch {
		case isDigit(string(str[0])):
			stringText := readString(str)
			text = text + stringText.String
			length = length + stringText.Length
			str = str[stringText.Length:]
		case str[0] == 'i':
			intText := readInt(str)
			text = text + intText.String
			length = length + intText.Length
			str = str[intText.Length:]
		case str[0] == 'l':
			listText := readList(str[1:])
			text = text + listText.String
			length = length + listText.Length
			str = str[listText.Length:]
		case str[0] == 'd':
			dictText := readDict(str)
			text = text + dictText.String
			length = length + dictText.Length
			str = str[dictText.Length:]
		case str[0] == 'e':
			text = text + "e"
			length = length + 1
			loopComplete = true
			str = str[1:]
		}
	}
	return BencodeText{String: text, Length: length}
}

func readDict(str string) BencodeText {
	length := 1
	text := "d"
	loopComplete := false
	for len(str) > 1 {
		if loopComplete == true {
			break
		}
		switch {
		case isDigit(string(str[0])):
			stringText := readString(str)
			text = text + stringText.String
			length = length + stringText.Length
			str = str[stringText.Length:]
		case str[0] == 'i':
			intText := readInt(str)
			text = text + intText.String
			length = length + intText.Length
			str = str[intText.Length:]
		case str[0] == 'l':
			listText := readList(str)
			text = text + listText.String
			length = length + listText.Length
			str = str[listText.Length:]
		case str[0] == 'd':
			dictText := readDict(str[1:])
			text = text + dictText.String
			length = length + dictText.Length
			str = str[dictText.Length:]
		case str[0] == 'e':
			text = text + "e"
			length = length + 1
			loopComplete = true
			str = str[1:]
		}
	}
	return BencodeText{String: text, Length: length}
}
