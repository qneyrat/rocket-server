package main

import (
	"fmt"

	"rocket-server/crypto"
)

func main() {
	msg := `{"method": "hello"}`
	c := &crypto.Crypto{Key: []byte("0123456789012345")}
	res, _ :=c.Encrypt([]byte(msg))
	fmt.Println(string(res))
}
