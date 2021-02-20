package fileDB

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// define type implementing Context
type processCtx struct {
	rootDir string
	maxAge int64
}

func (ctx *processCtx) GetDataDir() string {
	return filepath.Join(ctx.rootDir, "./data/")
}

func (ctx *processCtx) GetTmpDir() string {
	return filepath.Join(ctx.rootDir, "./tmp/")
}

func (ctx *processCtx) GetMaxTmpAge() int64 {
	return ctx.maxAge
}

// define interface
type Context interface {
	GetDataDir() string
	GetTmpDir() string
	GetMaxTmpAge() int64
}

// define factory
func ParseArgs() Context {
	rootDir := flag.String("root", "<none>", "the root file directory")
	maxAge := flag.Int64("maxage", 24*60*60, "the max age of a file in tmp (in sec)")
	flag.Parse()

	if *rootDir == "<none>" {
		_, _ = fmt.Fprintln(os.Stderr, "root Parameter is required")
		os.Exit(1)
	}

	ctx := processCtx{
		rootDir: *rootDir,
		maxAge: *maxAge,
	}
	return &ctx
}