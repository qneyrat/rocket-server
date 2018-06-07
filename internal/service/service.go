package service

import (
	"fmt"
)

type Service struct {}

func  (Service) SayHello(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
