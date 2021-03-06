package fileDB

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func CleanTMP(ctx Context) error {
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
		if time.Now().Sub(info.ModTime()) > ctx.GetMaxTmpAge() {
			// if so, delete
			err := os.Remove(filepath.Join(ctx.GetTmpDir(), info.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateDirIfNotExists(dir string) error {
	file, err := os.Open(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New(dir + " already exists as non-directory")
	}
	return nil
}

func getDefaultQueryInt64(val url.Values, key string, def int64) (int64, bool) {
	r := val.Get(key)
	if r == "" {
		return def, true
	}
	i, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return def, true
	}
	return i, false
}
