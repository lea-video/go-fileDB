package main

import (
	"github.com/lea-video/go-fileDB/fileDB"
)

func main() {
	// Parse cli args
	ctx := fileDB.ParseArgs()

	// Enforce existence of data / tmp dir
	err := fileDB.CreateDirIfNotExists(ctx.GetDataDir())
	panicOn(err)
	err = fileDB.CreateDirIfNotExists(ctx.GetTmpDir())
	panicOn(err)

	// Cleanup the tmp Folder after Start
	err = fileDB.CleanTMP(ctx)
	panicOn(err)

	err = fileDB.RegisterServer()
	panicOn(err)
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}