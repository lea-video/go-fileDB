package fileDB

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// define request type
type postRequest struct {
	ReqBody io.ReadCloser
	Offset int64
	FSize int64
	ChunkSize int64
	Path string
}
// define response type
type postResponse struct {
	Done bool
	Req *postRequest
	StatusCode int
	Error string
}

func (resp *postResponse) GetStatusCode() int {
	return resp.StatusCode
}
func (resp *postResponse) GetError() string {
	return resp.Error
}
func (resp *postResponse) DidFail() bool {
	return resp.Error == ""
}

// define postResponse factorys
func newFailedPostResponse(req *postRequest, err string, statusCode int) *postResponse {
	return &postResponse{
		Req:        req,
		Error:      err,
		StatusCode: statusCode,
	}
}
func newSuccessfullPostResponse(req *postRequest, done bool) *postResponse {
	return &postResponse{
		Req:        req,
		StatusCode: http.StatusOK,
		Done: done,
	}
}

// general functions
func handlePostRequest(ctx Context, req *postRequest) *postResponse {
	if req.FSize <= 0 {
		return newFailedPostResponse(req, "FSize invalid", http.StatusBadRequest)
	}
	if req.ChunkSize <= 0 || req.ChunkSize + req.Offset > req.FSize {
		return newFailedPostResponse(req, "ChunkSize invalid", http.StatusBadRequest)
	}

	// Merge request Path with root directory
	p := filepath.Join(ctx.GetTmpDir(), "."+req.Path)

	// Resolve to absolute Path
	absPath, err := filepath.Abs(p)
	if err != nil {
		return newFailedPostResponse(req, "invalid Path", http.StatusNotFound)
	}
	// Check that the relative Path is still inside the root Directory
	if !strings.HasPrefix(absPath, ctx.GetTmpDir()) {
		return newFailedPostResponse(req, "leaving the directory is not allowed", http.StatusForbidden)
	}

	// Open the File / create if not exists
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// TODO: f.Close()
	if err != nil {
		return newFailedPostResponse(req, "failed to open file", http.StatusBadRequest)
	}

	// Get the file size
	stat, err := f.Stat()
	if err != nil {
		return newFailedPostResponse(req, "failed to open file", http.StatusBadRequest)
	}
	fSize := stat.Size()

	// only allow appending
	if fSize != req.Offset {
		return newFailedPostResponse(req, "invalid offset", http.StatusBadRequest)
	}

	// do writing
	// TODO: check if this works
	_, err = f.Seek(req.Offset, 0)
	if err != nil {
		return newFailedPostResponse(req, "failed to write", http.StatusInternalServerError)
	}
	defer req.ReqBody.Close()
	_, err = io.Copy(f, req.ReqBody)
	if err != nil {
		return newFailedPostResponse(req, "failed to write", http.StatusInternalServerError)
	}

	// not done with file:
	if req.ChunkSize + req.Offset != req.FSize {
		return newSuccessfullPostResponse(req, false)
	}

	// done with file:
	newP := filepath.Join(ctx.GetDataDir(), "."+req.Path)
	err = os.Rename(p, newP)
	if err != nil {
		return newFailedPostResponse(req, "failed to move file", http.StatusInternalServerError)
	}
	return newSuccessfullPostResponse(req, true)
}

func buildPostRequest(r *http.Request) *postRequest {
	offset, _ := getDefaultQueryInt64(r.URL.Query(), "offset", 0)
	if offset < 0 {
		offset = 0;
	}
	chunkSize, _ := getDefaultQueryInt64(r.URL.Query(), "chunksize", -1)
	fsize, _ := getDefaultQueryInt64(r.URL.Query(), "fsize", -1)

	return &postRequest{
		Offset: offset,
		ReqBody: r.Body,
		FSize:  fsize,
		ChunkSize: chunkSize,
		Path:   r.URL.Path,
	}
}

func onPostRequest(ctx Context, r *http.Request) serverResponse {
	req := buildPostRequest(r)
	resp := handlePostRequest(ctx, req)
	return resp
}

// utility functions