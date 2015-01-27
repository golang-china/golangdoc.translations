// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package commands defines and manages the basic pprof commands
package commands

// Command describes the actions for a pprof command. Includes a function for
// command-line completion, the report format to use during report generation, any
// postprocessing functions, and whether the command expects a regexp parameter
// (typically a function name).
type Command struct {
	Complete    Completer     // autocomplete for interactive mode
	Format      int           // report format to generate
	PostProcess PostProcessor // postprocessing to run on report
	HasParam    bool          // Collect a parameter from the CLI
	Usage       string        // Help text
}

// Commands describes the commands accepted by pprof.
type Commands map[string]*Command

// PProf returns the basic pprof report-generation commands
func PProf(c Completer, interactive **bool, svgpan **string) Commands

// Completer is a function for command-line autocompletion
type Completer func(prefix string) string

// NewCompleter creates an autocompletion function for a set of commands.
func NewCompleter(cs Commands) Completer

// PostProcessor is a function that applies post-processing to the report output
type PostProcessor func(input *bytes.Buffer, output io.Writer, ui plugin.UI) error
