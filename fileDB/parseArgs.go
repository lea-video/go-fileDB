package fileDB

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// define type implementing Context
type processCtx struct {
	rootDir string
	maxAge time.Duration
}

func (ctx *processCtx) GetDataDir() string {
	return filepath.Join(ctx.rootDir, "./data/")
}

func (ctx *processCtx) GetTmpDir() string {
	return filepath.Join(ctx.rootDir, "./tmp/")
}

func (ctx *processCtx) GetMaxTmpAge() time.Duration {
	return ctx.maxAge
}

// define interface
type Context interface {
	GetDataDir() string
	GetTmpDir() string
	GetMaxTmpAge() time.Duration
}

// define factory
func ParseArgs() Context {
	rootDir := flag.String("root", "<none>", "ArgError: the root file directory")
	maxAge := flag.Uint("maxage", 24*60*60, "ArgError: the max age of a file in tmp (in sec)")
	flag.Parse()

	if *rootDir == "<none>" {
		_, _ = fmt.Fprintln(os.Stderr, "root Parameter is required")
		os.Exit(1)
	}

	ctx := processCtx{
		rootDir: *rootDir,
		maxAge: time.Duration(*maxAge) * time.Second,
	}
	return &ctx
}