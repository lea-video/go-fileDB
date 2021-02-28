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
	port uint
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

func (ctx *processCtx) GetPort() uint {
	return ctx.port
}

func (ctx *processCtx) GetPortStr()  string {
	return fmt.Sprintf(":%d", ctx.GetPort())
}

// define interface
type Context interface {
	GetDataDir() string
	GetTmpDir() string
	GetMaxTmpAge() time.Duration
	GetPort() uint
	GetPortStr() string
}

// define factory
func ParseArgs() Context {
	rootDir := flag.String("root", "<none>", "the root file directory")
	maxAge := flag.Uint("maxage", 24*60*60, "the max age of a file in tmp (in sec)")
	port := flag.Uint("port", 8080, "the port to listen on for http requests")
	flag.Parse()

	if *rootDir == "<none>" {
		_, _ = fmt.Fprintln(os.Stderr, "ArgError: root Parameter is required")
		os.Exit(1)
	}

	ctx := processCtx{
		rootDir: *rootDir,
		maxAge: time.Duration(*maxAge) * time.Second,
		port: *port,
	}
	return &ctx
}