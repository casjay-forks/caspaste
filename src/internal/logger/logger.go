
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package logger

import (
	"fmt"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Logger struct {
	TimeFormat string
}

func New(timeFormat string) Logger {
	return Logger{
		TimeFormat: timeFormat,
	}
}

func getTrace() string {
	trace := ""

	for i := 2; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			trace = trace + file + "#" + strconv.Itoa(line) + ": "

		} else {
			return trace
		}
	}
}

func (cfg Logger) Info(msg string) {
	fmt.Fprintln(os.Stdout, time.Now().Format(cfg.TimeFormat), "[INFO]   ", msg)
}

func (cfg Logger) Error(e error) {
	fmt.Fprintln(os.Stderr, time.Now().Format(cfg.TimeFormat), "[ERROR]  ", getTrace(), e.Error())
}

func (cfg Logger) HttpRequest(req *http.Request, code int) {
	fmt.Fprintln(os.Stdout, time.Now().Format(cfg.TimeFormat), "[REQUEST]", netshare.GetClientAddr(req).String(), req.Method, code, req.URL.Path, "(User-Agent: "+req.UserAgent()+")")
}

func (cfg Logger) HttpError(req *http.Request, e error) {
	fmt.Fprintln(os.Stderr, time.Now().Format(cfg.TimeFormat), "[ERROR]  ", netshare.GetClientAddr(req).String(), req.Method, 500, req.URL.Path, "(User-Agent: "+req.UserAgent()+")", "Error:", getTrace(), e.Error())
}
