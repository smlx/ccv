// Package main implements the ccv command line tool.
package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/smlx/ccv"
)

// CLI represents the command-line interface.
type CLI struct {
	VersionType bool `kong:"env='VERSION_TYPE',help='Print new version type (major, minor, patch) instead of version number'"`
}

// Run implements the CLI logic.
func (cli *CLI) Run() error {
	if cli.VersionType {
		nextVersionType, err := ccv.NextVersionType(`.`)
		if err != nil {
			return fmt.Errorf("couldn't get next version type: %v", err)
		}
		fmt.Println(nextVersionType)
	} else {
		nextVersion, err := ccv.NextVersion(`.`)
		if err != nil {
			return fmt.Errorf("couldn't get next version: %v", err)
		}
		fmt.Println(nextVersion)
	}
	return nil
}

func main() {
	cli := CLI{}
	kctx := kong.Parse(&cli, kong.UsageOnError())
	kctx.FatalIfErrorf(kctx.Run())
}
