package fileDB

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// define request type
type deleteRequest struct {
	Path string
}
// define response type
type deleteResponse struct {
	Req *deleteRequest
	StatusCode int
	Error string
}

func (resp *deleteResponse) GetStatusCode() int {
	return resp.StatusCode
}
func (resp *deleteResponse) GetError() string {
	return resp.Error
}
func (resp *deleteResponse) DidFail() bool {
	return resp.Error == ""
}

// define deleteResponse factorys
func newFailedDeleteResponse(req *deleteRequest, err string, statusCode int) *deleteResponse {
	return &deleteResponse{
		Req:        req,
		Error:      err,
		StatusCode: statusCode,
	}
}
func newSuccessfullDeleteResponse(req *deleteRequest) *deleteResponse {
	return &deleteResponse{
		Req:        req,
		StatusCode: http.StatusOK,
	}
}

// general functions
func handleDeleteRequest(ctx Context, req *deleteRequest) *deleteResponse {
	// Merge request Path with root directory
	p := filepath.Join(ctx.GetDataDir(), "."+req.Path)

	// Resolve to absolute Path
	absPath, err := filepath.Abs(p)
	if err != nil {
		return newFailedDeleteResponse(req, "invalid Path", http.StatusNotFound)
	}
	// Check that the relative Path is still inside the root Directory
	if !strings.HasPrefix(absPath, ctx.GetDataDir()) {
		return newFailedDeleteResponse(req, "leaving the directory is not allowed", http.StatusForbidden)
	}

	err = os.Remove(absPath)
	if err != nil {
		return newFailedDeleteResponse(req, "unable to delete the file", http.StatusInternalServerError)
	}
	return newSuccessfullDeleteResponse(req)
}

func buildDeleteRequest(r *http.Request) *deleteRequest {
	return &deleteRequest{
		Path: r.URL.Path,
	}
}

func onDeleteRequest(ctx Context, r *http.Request) serverResponse {
	req := buildDeleteRequest(r)
	resp := handleDeleteRequest(ctx, req)
	return resp
}