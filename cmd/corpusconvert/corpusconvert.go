package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

/*
"TIPLOCDATA": [
	{
		"NLCDESC": "MERSEYRAIL ELECTRICS-HQ INPUT",
		"NLC": "000800",
		"TIPLOC": " ",
		"3ALPHA": " ",
		"STANOX": " ",
		"UIC": " ",
		"NLCDESC16": "MPTE HQ INPUT"
	},
*/

type Corpus struct {
	TIPLOCDATA []CorpusLocation
}

type CorpusLocation struct {
	NLCDESC string
	TIPLOC  string
}

func main() {

	corpus, err := ioutil.ReadFile("CORPUSExtract.json")

	var decoded Corpus
	err = json.Unmarshal(corpus, &decoded)

	if err != nil {
		panic(err)
	}

	tiplocs := make(map[string]string)

	for _, loc := range decoded.TIPLOCDATA {
		if len(loc.TIPLOC) > 0 {
			tiplocs[loc.TIPLOC] = loc.NLCDESC
		}
	}

	outfile, err := os.OpenFile("tiploc.bin", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	writer := bufio.NewWriter(outfile)
	enc := msgpack.NewEncoder(writer)

	err = enc.Encode(tiplocs)
	if err != nil {
		panic(err)
	}

	outfile.Close()
}
