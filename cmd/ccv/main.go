// Package main implements the ccv command line tool.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/smlx/ccv"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	next, err := ccv.NextVersion(`.`)
	if err != nil {
		log.Error("couldn't get next version", slog.Any("error", err))
		os.Exit(1)
	}
	fmt.Println(next)
}
