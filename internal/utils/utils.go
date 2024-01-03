package utils

import (
	"fmt"
	"time"
)

func GenerateId() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
