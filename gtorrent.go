package main

import (
	"fmt"
	"gtorrent/bencode"
)

func main() {
	fmt.Println(bencode.ParseString("3:cow"))
	fmt.Println(bencode.ParseInt("i12e"))
	fmt.Println(bencode.ParseInt("i0e"))

}
