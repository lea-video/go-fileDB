package fileDB

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// define type implementing getRequestData
type getRequestData struct {
	StartByte int64
	EndByte   int64
	Path      string
}

func (rd *getRequestData) GetStartByte() int64 {
	return rd.StartByte
}

func (rd *getRequestData) GetEndByte() int64 {
	return rd.EndByte
}

func (rd *getRequestData) HasEnd() bool {
	return rd.EndByte != -1
}

func (rd *getRequestData) GetChunkLength() int64 {
	return rd.EndByte - rd.StartByte
}

func (rd *getRequestData) GetPath() string {
	return rd.Path
}

type GetRequest interface {
	GetStartByte() int64
	GetEndByte() int64
	HasEnd() bool
	GetChunkLength() int64
	GetPath() string
}

// define factory
func newGetRequest(startByte, endByte int64, path string) GetRequest {
	if startByte < 0 {
		startByte = 0
	}
	return &getRequestData{
		StartByte: startByte,
		EndByte:   endByte,
		Path:      path,
	}
}

// other stuff
type GetResponse struct {
	Req GetRequest
	Data []byte
	DataLength int64
	IsAtEnd bool
	FSize int64
	Error string
}

func handleGet(ctx Context, req GetRequest) GetResponse {
	// Merge request path with root directory
	p := filepath.Join(ctx.GetDataDir(), "." + req.GetPath())

	// Resolve to absolute Path
	absPath, err := filepath.Abs(p)
	if err != nil {
		return GetResponse{
			Req: req,
			Error: "invalid path",
		}
	}
	// Check that the relative path is still inside the root Directory
	if !strings.HasPrefix(absPath, ctx.GetDataDir()) {
		return GetResponse{
			Req: req,
			Error: "leaving the directory is not allowed",
		}
	}

	// Open the File
	f, err := os.Open(p)
	if err != nil {
		return GetResponse{
			Req: req,
			Error: "failed to open file",
		}
	}

	// Get the file size
	stat, err := f.Stat()
	if err != nil {
		return GetResponse{
			Req: req,
			Error: "failed to open file",
		}
	}
	fSize := stat.Size()

	// Range-Validation
	// TODO: check if some of these can be default to other vars
	if fSize <= req.GetStartByte() {
		return GetResponse{
			Req: req,
			Error: "Range outside of file",
		}
	}
	if req.HasEnd() && req.GetEndByte() < req.GetStartByte() {
		return GetResponse{
			Req: req,
			Error: "End before Start",
		}
	}

	// calculate the bytes to fetch
	// TODO: max value for toFetch
	toFetch := fSize - req.GetStartByte()
	isAtEnd := true
	if req.HasEnd() {
		// Calculate if req.chunkLength is longer then the file
		toFetch = min(toFetch, req.GetChunkLength())
		if fSize > req.GetStartByte() + toFetch {
			isAtEnd = false
		}
	}

	// Read file content
	b := make([]byte, toFetch)
	_, err = f.ReadAt(b, req.GetStartByte())
	if err != nil {
		return GetResponse{
			Req: req,
			Error: "failed to read file",
		}
	}

	return GetResponse{
		Req:     req,
		IsAtEnd: isAtEnd,
		FSize:   fSize,
		Data:    b,
		DataLength: toFetch,
	}
}

func HandleGetRequest(ctx Context, w http.ResponseWriter, r *http.Request) {
	start := getDefaultInt64(r.URL.Query(), "start", 0)
	// TODO: support "length" instead of "end"
	end := getDefaultInt64(r.URL.Query(), "end", -1)
	req := newGetRequest(start, end, r.URL.Path)
	resp := handleGet(ctx, req)

	b, err := json.Marshal(resp)
	// TODO: replace error handling
	if err != nil {
		fmt.Println(err)
	} else {
		_, err = fmt.Fprint(w, string(b))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getDefaultInt64(val url.Values, key string, def int64) int64 {
	r := val.Get(key)
	if r == "" {
		return def;
	}
	i, err := strconv.ParseInt(r, 10, 64);
	if err != nil {
		return def;
	}
	return i;
}
