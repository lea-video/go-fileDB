package fileDB

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// define request type
type getRequest struct {
	StartByte  int64
	ByteLength int64
	Path       string
}

func (req *getRequest) calculateRequestLength(ctx Context, fileSize int64) int64 {
	var size int64 = 0
	if req.ByteLength < 0 {
		size = req.StartByte - fileSize
	} else {
		size = req.StartByte - req.ByteLength
	}

	if ctx.HasMaxChunkSizeLimit() {
		size = min(ctx.GetMaxChunkSize(), size)
	}
	return size
}

// define getRequest factory
func newGetRequest(startByte, length int64, path string) *getRequest {
	if startByte < 0 {
		startByte = 0
	}
	return &getRequest{
		StartByte:  startByte,
		ByteLength: length,
		Path:       path,
	}
}

// define response type
type getResponse struct {
	Req        *getRequest
	Data       []byte
	DataLength int64
	IsAtEnd    bool
	FSize      int64
	Error      string
	StatusCode int
}

func (resp *getResponse) GetStatusCode() int {
	return resp.StatusCode
}
func (resp *getResponse) GetError() string {
	return resp.Error
}
func (resp *getResponse) DidFail() bool {
	return resp.Error == ""
}

// define getResponse factorys
func newFailedGetResponse(req *getRequest, err string, statusCode int) *getResponse {
	return &getResponse{
		Req:        req,
		Error:      err,
		StatusCode: statusCode,
	}
}
func newSuccessfullGetResponse(req *getRequest, data []byte, datasize, filesize int64, isAtEnd bool) *getResponse {
	return &getResponse{
		Req:        req,
		Data:       data,
		DataLength: datasize,
		IsAtEnd:    isAtEnd,
		FSize:      filesize,
		StatusCode: http.StatusOK,
	}
}

// general functions

func handleGetRequest(ctx Context, req *getRequest) *getResponse {
	// Merge request Path with root directory
	p := filepath.Join(ctx.GetDataDir(), "."+req.Path)

	// Resolve to absolute Path
	absPath, err := filepath.Abs(p)
	if err != nil {
		return newFailedGetResponse(req, "invalid Path", http.StatusNotFound)
	}
	// Check that the relative Path is still inside the root Directory
	if !strings.HasPrefix(absPath, ctx.GetDataDir()) {
		return newFailedGetResponse(req, "leaving the directory is not allowed", http.StatusForbidden)
	}

	// Open the File
	f, err := os.Open(p)
	if err != nil {
		return newFailedGetResponse(req, "failed to open file", http.StatusBadRequest)
	}

	// Get the file size
	stat, err := f.Stat()
	if err != nil {
		return newFailedGetResponse(req, "failed to open file", http.StatusBadRequest)
	}
	fSize := stat.Size()

	// Range-Validation
	if fSize <= req.StartByte {
		return newFailedGetResponse(req, "start is outside of file", http.StatusBadRequest)
	}

	// calculate the bytes to fetch
	toFetch := req.calculateRequestLength(ctx, fSize)
	isAtEnd := req.StartByte+toFetch == fSize

	// Read file content
	b := make([]byte, toFetch)
	_, err = f.ReadAt(b, req.StartByte)
	if err != nil {
		return newFailedGetResponse(req, "failed to read file", http.StatusInternalServerError)
	}

	return newSuccessfullGetResponse(req, b, toFetch, fSize, isAtEnd)
}

func buildGetRequest(r *http.Request) *getRequest {
	// read all params
	start, _ := getDefaultInt64(r.URL.Query(), "start", 0)
	end, didDefaultEnd := getDefaultInt64(r.URL.Query(), "end", -1)
	length, didDefaultLength := getDefaultInt64(r.URL.Query(), "length", -1)
	// overwrite length with end only when no valid length was found
	if didDefaultLength && !didDefaultEnd {
		length = start - end
	}

	return newGetRequest(start, length, r.URL.Path)
}

func HandleGetRequest(ctx Context, _ http.ResponseWriter, r *http.Request) serverResponse {
	req := buildGetRequest(r)
	resp := handleGetRequest(ctx, req)
	return resp
}

func getDefaultInt64(val url.Values, key string, def int64) (int64, bool) {
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
