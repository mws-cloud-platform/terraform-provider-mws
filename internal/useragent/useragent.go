package useragent

import (
	"fmt"
	"runtime"
	"strings"

	"go.mws.cloud/util-toolset/pkg/os/env"
)

func New(env env.Env, product, version string) string {
	comments := make([]string, 0, 3)
	for _, comment := range []string{runtime.GOOS, runtime.GOARCH, DetectShell(env)} {
		if comment == "" {
			comment = "unknown"
		}
		comments = append(comments, comment)
	}

	return fmt.Sprintf("%s/%s (%s)", product, version, strings.Join(comments, "; "))
}

func DetectShell(env env.Env) string {
	if shell := env.Getenv("SHELL"); shell != "" {
		return shell[strings.LastIndex(shell, "/")+1:]
	}
	return "other"
}
