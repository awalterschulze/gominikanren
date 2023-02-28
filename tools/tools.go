//go:build tools
// +build tools

package tools

import (
	_ "github.com/awalterschulze/goderive"
	_ "github.com/goccmack/gocc"
	_ "github.com/kisielk/errcheck"
	_ "github.com/kisielk/gotool"
	_ "golang.org/x/lint/golint"
)
