package lofty

import (
	"os"
)

func isPiping() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return !(info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0)
}
