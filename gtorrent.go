package main

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
	"gtorrent/bencode"
	"gtorrent/torrentfile"
	"os"
	"strings"
)

func main() {
	if os.Args[1] != "" {
		file, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
		var v torrentfile.TorrentFile
		err = bencode.Unmarshall(file, &v)
		// t, err := bencode.Parse(string(file))
		if err != nil {
			fmt.Println(err)
		}

		pp.Println(v)
		// pp.Println("raw value:", t.Value())

	}

}

func rs(str string) string {
	return strings.ReplaceAll(str, " ", "")
}
