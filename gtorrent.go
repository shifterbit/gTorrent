package main

import (
	"fmt"
	"gtorrent/bencode"
	"os"
)

func main() {
	path := os.Args[1]
	_, err := bencode.ParseFile(path)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	fmt.Println(bencode.ParseString("3:cow"))

}
