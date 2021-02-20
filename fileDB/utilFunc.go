package fileDB

import (
	"os"
	"path/filepath"
	"time"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func unixSec(t time.Time) int64 {
	return t.Unix() - (t.Unix() % 1000)
}

func secAge(t1 time.Time) int64 {
	return unixSec(t1) - unixSec(time.Now())
}

func CleanRoot(ctx Context) error {
	// get files in tmp dir
	files, err := os.ReadDir(ctx.GetTmpDir())
	if err != nil {
		return err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		info, err := file.Info()
		if err != nil {
			return err
		}
		// check if file wasn't modified for specified value
		if secAge(info.ModTime()) > ctx.GetMaxTmpAge() {
			// if so, delete
			err := os.Remove(filepath.Join(ctx.GetTmpDir(), info.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
