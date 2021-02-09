package main

import "go.uber.org/zap"

// NextVersion is a test fixture to expose the private nextVersion() function.
func NextVersion(log *zap.Logger, path string) (string, error) {
	return nextVersion(path)
}
