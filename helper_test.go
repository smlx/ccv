package main

// NextVersion is a test fixture to expose the private nextVersion() function.
func NextVersion(path string) (string, error) {
	return nextVersion(path)
}
