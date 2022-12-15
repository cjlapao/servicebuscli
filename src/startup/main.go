package startup

import (
	"os"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/servicebuscli/common"
)

func Start() {
	if !helper.FileExists(common.TempFolder()) {
		os.MkdirAll(common.TempFolder(), os.ModePerm)
	}
}

func Cleanup() error {
	if helper.FileExists(common.TempFolder()) {
		err := os.Remove(common.TempFolder())
		if err != nil {
			return err
		}
	}

	return nil
}

func Exit(exitCode int) {
	// Cleanup()
	os.Exit(exitCode)
}
