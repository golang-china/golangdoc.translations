// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package fetch provides an extensible mechanism to fetch a profile
// from a data source.

// Package fetch provides an extensible mechanism to fetch a profile from a data
// source.
package fetch

import (
	"cmd/internal/pprof/plugin"
	"cmd/internal/pprof/profile"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// FetchProfile reads from a data source (network, file) and generates a
// profile.
func FetchProfile(source string, timeout time.Duration) (*profile.Profile, error)

// FetchURL fetches a profile from a URL using HTTP.
func FetchURL(source string, timeout time.Duration) (io.ReadCloser, error)

// Fetcher is the plugin.Fetcher version of FetchProfile.
func Fetcher(source string, timeout time.Duration, ui plugin.UI) (*profile.Profile, error)

// PostURL issues a POST to a URL over HTTP.
func PostURL(source, post string) ([]byte, error)

