package main

import "github.com/lea-video/go-fileDB/fileDB"

func main() {
	// Parse cli args
	ctx := fileDB.ParseArgs()
	// Cleanup the tmp Folder after Start
	err := fileDB.CleanRoot(ctx)
	panicOn(err)

	// Simulate a read request
	req := fileDB.NewGetRequest(0, 10, "/yey/boy")
	fileDB.HandleGet(ctx, req)
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}