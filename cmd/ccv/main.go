// Package main implements the ccv command line tool.
package main

import (
	"fmt"

	"github.com/smlx/ccv"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	next, err := ccv.NextVersion(`.`)
	if err != nil {
		log.Fatal("couldn't get next version", zap.Error(err))
	}
	fmt.Println(next)
}
