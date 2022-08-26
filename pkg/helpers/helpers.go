package helpers

import "os"

// IsExistingDirOrFile expects a path and returns true if is a directory,
// false if is a file and an error if does not exist or can't be accessed.
func IsExistingDirOrFile(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !stat.IsDir() {
		return false, nil
	}
	return true, nil
}
