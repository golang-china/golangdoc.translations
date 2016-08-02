// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Trace is a tool for viewing trace files.
//
// Trace files can be generated with:
//     - runtime/trace.Start
//     - net/http/pprof package
//     - go test -trace
//
// Example usage:
// Generate a trace file with 'go test':
//     go test -trace trace.out pkg
// View the trace in a web browser:
//     go tool trace pkg.test trace.out
package main // go get cmd/trace

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "html/template"
    "internal/trace"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "sort"
    "strconv"
    "strings"
    "sync"
)

type NameArg struct {
    Name string `json:"name"`
}

// Record represents one entry in pprof-like profiles.
type Record struct {
    stk  []*trace.Frame
    n    uint64
    time int64
}

type SortIndexArg struct {
    Index int `json:"sort_index"`
}

type ViewerData struct {
    Events   []*ViewerEvent         `json:"traceEvents"`
    Frames   map[string]ViewerFrame `json:"stackFrames"`
    TimeUnit string                 `json:"displayTimeUnit"`
}

type ViewerEvent struct {
    Name     string      `json:"name,omitempty"`
    Phase    string      `json:"ph"`
    Scope    string      `json:"s,omitempty"`
    Time     float64     `json:"ts"`
    Dur      float64     `json:"dur,omitempty"`
    Pid      uint64      `json:"pid"`
    Tid      uint64      `json:"tid"`
    ID       uint64      `json:"id,omitempty"`
    Stack    int         `json:"sf,omitempty"`
    EndStack int         `json:"esf,omitempty"`
    Arg      interface{} `json:"args,omitempty"`
}

type ViewerFrame struct {
    Name   string `json:"name"`
    Parent int    `json:"parent,omitempty"`
}

