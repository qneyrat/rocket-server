package main

import (
	"fmt"

	"rocket-server/crypto"
)

func main() {
	msg := `ShWqxl0gX0xqDW+Ljrs2DIH1lCGDGV+yv1QwzLJ8nDT4M9xi+re0`
	c := &crypto.Crypto{Key: []byte("0123456789012345")}
	res, _ :=c.Decrypt([]byte(msg))
	fmt.Println(string(res))
}
