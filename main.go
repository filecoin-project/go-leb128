package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"encoding/base64"
)


const USAGE = `
USAGE 
COMMANDS
        dec          decodes a LEB128 uint64 
       
        dec-big      decodes a LEB128 signed big integer

        enc          encodes a LEB128 uint64

        enc-big      encodes a LEB128 signed big integer

EXAMPLES
        To decode from leb128:
        leb128 dec <leb128 uint64>
        leb128 dec-big <leb128 big int>

        To encode to leb128:
        leb128 enc 100
        leb128 enc-big 10000000000000000000000000
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(USAGE)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "-h", "--help":
		fmt.Print(USAGE)
		os.Exit(1)
	case "dec-big":
		// Base64 string to bytes
		buf := make([]byte, 64)
		decodedBytes, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(buf)))
		if err != nil {
			fmt.Printf("Error decoding base64 string: %s\n", err.Error())
			os.Exit(1)
		}
		// leb128 bytes to big int
		b := ToBigInt(decodedBytes)
		fmt.Printf("%s", b.Text(10))
		
		os.Exit(0)
	case "dec":
		fmt.Printf("Implement me!\n")
		os.Exit(1)
	case "enc":
		fmt.Printf("Implement me!\n")
		os.Exit(1)
	case "enc-big":
		fmt.Printf("Implement me!\n")
		os.Exit(1)
	default:
		fmt.Printf("Invalid command\n")
		os.Exit(1)		
	}
}
