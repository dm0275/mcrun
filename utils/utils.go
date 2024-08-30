package utils

import (
	"fmt"
	"log"
)

var (
	logger = log.Default()
)

func CheckErr(err error) {
	if err != nil {
		logger.Fatalf(fmt.Sprintf("ERROR: %s", err))
	}
}
