package fileDB

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// define type implementing GetRequest
type getRequestData struct {
	startByte int64
	endByte int64
	path string
}

func (rd *getRequestData) GetStartByte() int64 {
	return rd.startByte
}

func (rd *getRequestData) GetEndByte() int64 {
	return rd.endByte
}

func (rd *getRequestData) GetChunkLength() int64 {
	return rd.endByte - rd.startByte
}

func (rd *getRequestData) GetPath() string {
	return rd.path
}

// define interface
type GetRequest interface {
	GetStartByte() int64
	GetEndByte() int64
	GetChunkLength() int64
	GetPath() string
}

// define factory
func NewGetRequest(startByte, endByte int64, path string) GetRequest {
	return &getRequestData{
		startByte:   startByte,
		endByte:     endByte,
		path:        path,
	}
}

// other stuff
type GetResponse struct {
	req GetRequest
	data []byte
	isAtEnd bool
	fSize int64
	err error
}

func HandleGet(ctx Context, req GetRequest) GetResponse {
	// Merge request path with root directory
	p := filepath.Join(ctx.GetDataDir(), "." + req.GetPath())

	// Resolve to absolute Path
	absPath, err := filepath.Abs(p)
	if err != nil {
		return GetResponse{
			req: req,
			err: err,
		}
	}
	// Check that the relative path is still inside the root Directory
	if !strings.HasPrefix(absPath, ctx.GetDataDir()) {
		return GetResponse{
			req: req,
			err: errors.New("not in root Directory"),
		}
	}

	// Open the File
	f, err := os.Open(p)
	if err != nil {
		return GetResponse{
			req: req,
			err: err,
		}
	}

	// Get the file size
	stat, err := f.Stat()
	if err != nil {
		return GetResponse{
			req: req,
			err: err,
		}
	}

	fSize := stat.Size()
	// Calculate if req.chunkLength is longer then the file
	toFetch := min(fSize - req.GetStartByte(), req.GetChunkLength())
	atEnd := false
	if fSize == req.GetStartByte() + toFetch {
		atEnd = true
	}

	// Read file content
	b := make([]byte, toFetch)
	_, err = f.ReadAt(b, req.GetStartByte())
	if err != nil {
		return GetResponse{
			req: req,
			err: err,
		}
	}

	return GetResponse{
		req: req,
		isAtEnd: atEnd,
		fSize: fSize,
		data: b,
	}
}
