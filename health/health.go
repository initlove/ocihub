package health

import (
	"fmt"
)

type Health interface {
	GetStatus() (string, string)
}

func RegisterHealth(name string, h Health) error {
	fmt.Println("rh")
	return nil
}
