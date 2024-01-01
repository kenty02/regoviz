package main

import (
	"fmt"
	"time"
)

func uid() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
