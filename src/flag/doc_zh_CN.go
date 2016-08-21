// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package flag implements command-line flag parsing.
//
// Usage:
//
// Define flags using flag.String(), Bool(), Int(), etc.
//
// This declares an integer flag, -flagname, stored in the pointer ip, with type
// *int.
//
//     import "flag"
//     var ip = flag.Int("flagname", 1234, "help message for flagname")
//
// If you like, you can bind the flag to a variable using the Var() functions.
//
//     var flagvar int
//     func init() {
//         flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
//     }
//
// Or you can create custom flags that satisfy the Value interface (with pointer
// receivers) and couple them to flag parsing by
//
//     flag.Var(&flagVal, "name", "help message for flagname")
//
// For such flags, the default value is just the initial value of the variable.
//
// After all flags are defined, call
//
//     flag.Parse()
//
// to parse the command line into the defined flags.
//
// Flags may then be used directly. If you're using the flags themselves, they
// are all pointers; if you bind to variables, they're values.
//
//     fmt.Println("ip has value ", *ip)
//     fmt.Println("flagvar has value ", flagvar)
//
// After parsing, the arguments following the flags are available as the slice
// flag.Args() or individually as flag.Arg(i). The arguments are indexed from 0
// through flag.NArg()-1.
//
// Command line flag syntax:
//
//     -flag
//     -flag=x
//     -flag x  // non-boolean flags only
//
// One or two minus signs may be used; they are equivalent. The last form is not
// permitted for boolean flags because the meaning of the command
//
//     cmd -x *
//
// will change if there is a file called 0, false, etc. You must use the
// -flag=false form to turn off a boolean flag.
//
// Flag parsing stops just before the first non-flag argument ("-" is a non-flag
// argument) or after the terminator "--".
//
// Integer flags accept 1234, 0664, 0x1234 and may be negative. Boolean flags
// may be:
//
//     1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False
//
// Duration flags accept any input valid for time.ParseDuration.
//
// The default set of command-line flags is controlled by top-level functions.
// The FlagSet type allows one to define independent sets of flags, such as to
// implement subcommands in a command-line interface. The methods of FlagSet are
// analogous to the top-level functions for the command-line flag set.

// flag 包实现命令行标签解析.
//
// 使用：
//
// 定义标签需要使用flag.String(),Bool(),Int()等方法。
//
// 下面的代码定义了一个interger标签，标签名是flagname，标签解析的结果存放在ip指
// 针（*int）指向的值中
//
//     import "flag"
//     var ip = flag.Int("flagname", 1234, "help message for flagname")
//
// 你还可以选择使用Var()函数将标签绑定到指定变量中
//
//     var flagvar int
//     func init() {
//         flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
//     }
//
// 你也可以传入自定义类型的标签，只要标签满足对应的值接口（接收指针指向的接收者
// ）。像下面代码一样定义标签
//
//     flag.Var(&flagVal, "name", "help message for flagname")
//
// 这样的标签，默认值就是自定义类型的初始值。
//
// 所有的标签都定义好了，就可以调用
//
//     flag.Parse()
//
// 来解析命令行参数并传入到定义好的标签了。
//
// 标签可以被用来直接使用。如果你直接使用标签（没有绑定变量），那他们都是指针类
// 型。如果你将他们绑定到变量上，他们就是值类型。
//
//     fmt.Println("ip has value ", *ip)
//     fmt.Println("flagvar has value ", flagvar)
//
// 在解析之后，标签对应的参数可以从flag.Args()获取到，它返回的slice，也可以使用
// flag.Arg(i)来获取单个参数。 参数列的索引是从0到flag.NArg()-1。
//
// 命令行标签格式：
//
//     -flag
//     -flag=x
//     -flag x  // 只有非boolean标签能这么用
//
// 减号可以使用一个或者两个，效果是一样的。 上面最后一种方式不能被boolean类型的
// 标签使用。因为当有个文件的名字是0或者false这样的词的话，下面的命令
//
//     cmd -x *
//
// 的原意会被改变。你必须使用-flag=false的方式来解析boolean标签。
//
// 一个标签的解析会在下次出现第一个非标签参数（“-”就是一个非标签参数）的时候停
// 止，或者是在终止符号“--”的时候停止。
//
// Interger标签接受如1234，0664，0x1234和负数这样的值。 Boolean标签接受1，0，t，
// f，true，false，TRUE，FALSE，True，False。 Duration标签接受任何可被
// time.ParseDuration解析的值。
//
// 默认的命令行标签是由最高层的函数来控制的。FlagSet类型允许每个包定义独立的标签
// 集合，例如在命令行接口中实现子命令。 FlagSet的方法就是模拟使用最高层函数来控
// 制命令行标签集的行为的。
package flag

import (
    "errors"
    "fmt"
    "io"
    "os"
    "sort"
    "strconv"
    "time"
)

// These constants cause FlagSet.Parse to behave as described if the parse
// fails.
const (
    ContinueOnError ErrorHandling = iota
    ExitOnError
    PanicOnError
)

// CommandLine is the default set of command-line flags, parsed from os.Args.
// The top-level functions such as BoolVar, Arg, and so on are wrappers for the
// methods of CommandLine.

// CommandLine 是命令行标记的默认集合，从 os.Args 解析而来。像 BoolVar、Arg
// 等这样的顶级函数为 CommandLine 方法的包装。
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

// ErrHelp is the error returned if the -help or -h flag is invoked
// but no such flag is defined.

// ErrHelp 在 -help 或 -h
// 标志未定义却调用了它时返回一个错误。
var ErrHelp = errors.New("flag: help requested")

// Usage prints to standard error a usage message documenting all defined
// command-line flags. It is called when an error occurs while parsing flags.
// The function is a variable that may be changed to point to a custom function.
// By default it prints a simple header and calls PrintDefaults; for details
// about the format of the output and how to control it, see the documentation
// for PrintDefaults.

// Usage打印出标准的错误信息，包含所有定义过的命令行标签说明。
// 这个函数赋值到一个变量上去，当然也可以将这个变量指向到自定义的函数。
var Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
    PrintDefaults()
}

// ErrorHandling defines how FlagSet.Parse behaves if the parse fails.

// ErrorHandling定义了如何处理标签解析的错误
type ErrorHandling int

// A Flag represents the state of a flag.

// Flag表示标签的状态
type Flag struct {
    Name     string // name as it appears on command line
    Usage    string // help message
    Value    Value  // value as set
    DefValue string // default value (as text); for usage message
}

// A FlagSet represents a set of defined flags.  The zero value of a FlagSet
// has no name and has ContinueOnError error handling.

// FlagSet 是已经定义好的标签的集合。FlagSet 的零值没有名字且拥有
// ContinueOnError 错误处理。
type FlagSet struct {
    // Usage is the function called when an error occurs while parsing flags.
    // The field is a function (not a method) that may be changed to point to
    // a custom error handler.
    Usage func()
}

// Getter is an interface that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided by
// this package satisfy the Getter interface.
type Getter interface {
    Value
    Get() interface{}
}

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
//
// If a Value has an IsBoolFlag() bool method returning true,
// the command-line parser makes -name equivalent to -name=true
// rather than using the next command-line argument.
//
// Set is called once, in command line order, for each flag present.

// Value接口是定义了标签对应的具体的参数值。 （默认值是string类型）
//
// 若 Value 拥有的 IsBoolFlag() bool 方法返回 ture，则命令行解析器会使 -name 等
// 价于 -name=true，而非使用下一个命令行实参。
type Value interface {
    String() string
    Set(string) error
}

// Arg returns the i'th command-line argument. Arg(0) is the first remaining
// argument after flags have been processed. Arg returns an empty string if the
// requested element does not exist.

// Arg返回第i个命令行参数。当有标签被解析之后，Arg(0)就成为了保留参数。
func Arg(i int) string

// Args returns the non-flag command-line arguments.

// Args返回非标签的命令行参数。
func Args() []string

// Bool defines a bool flag with specified name, default value, and usage
// string. The return value is the address of a bool variable that stores the
// value of the flag.

// Bool定义了一个有指定名字，默认值，和用法说明的bool标签。
// 返回值是一个存储标签解析值的bool变量地址。
func Bool(name string, value bool, usage string) *bool

// BoolVar defines a bool flag with specified name, default value, and usage
// string. The argument p points to a bool variable in which to store the value
// of the flag.

// BoolVar定义了一个有指定名字，默认值，和用法说明的bool标签。
// 参数p指向一个存储标签解析值的bool变量。
func BoolVar(p *bool, name string, value bool, usage string)

// Duration defines a time.Duration flag with specified name, default value, and
// usage string. The return value is the address of a time.Duration variable
// that stores the value of the flag. The flag accepts a value acceptable to
// time.ParseDuration.

// Duration定义了一个有指定名字，默认值，和用法说明的time.Duration标签。 返回值
// 是一个存储标签解析值的time.Duration变量地址。 此标记接受一个
// time.ParseDuration 可接受的值。
func Duration(name string, value time.Duration, usage string) *time.Duration

// DurationVar defines a time.Duration flag with specified name, default value,
// and usage string. The argument p points to a time.Duration variable in which
// to store the value of the flag. The flag accepts a value acceptable to
// time.ParseDuration.

// DurationVar定义了一个有指定名字，默认值，和用法说明的time.Duration标签。 参数
// p指向一个存储标签解析值的time.Duration变量。 此标记接受一个
// time.ParseDuration 可接受的值。
func DurationVar(p *time.Duration, name string, value time.Duration, usage string)

// Float64 defines a float64 flag with specified name, default value, and usage
// string. The return value is the address of a float64 variable that stores the
// value of the flag.

// Float64定义了一个有指定名字，默认值，和用法说明的float64标签。
// 返回值是一个存储标签解析值的float64变量地址。
func Float64(name string, value float64, usage string) *float64

// Float64Var defines a float64 flag with specified name, default value, and
// usage string. The argument p points to a float64 variable in which to store
// the value of the flag.

// Float64Var定义了一个有指定名字，默认值，和用法说明的float64标签。
// 参数p指向一个存储标签解析值的float64变量。
func Float64Var(p *float64, name string, value float64, usage string)

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of
// the flag.

// Int定义了一个有指定名字，默认值，和用法说明的int标签。
// 返回值是一个存储标签解析值的int变量地址。
func Int(name string, value int, usage string) *int

// Int64 defines an int64 flag with specified name, default value, and usage
// string. The return value is the address of an int64 variable that stores the
// value of the flag.

// Int64定义了一个有指定名字，默认值，和用法说明的int64标签。
// 返回值是一个存储标签解析值的int64变量地址。
func Int64(name string, value int64, usage string) *int64

// Int64Var defines an int64 flag with specified name, default value, and usage
// string. The argument p points to an int64 variable in which to store the
// value of the flag.

// Int64Var定义了一个有指定名字，默认值，和用法说明的int64标签。
// 参数p指向一个存储标签解析值的int64变量。
func Int64Var(p *int64, name string, value int64, usage string)

// IntVar defines an int flag with specified name, default value, and usage
// string. The argument p points to an int variable in which to store the value
// of the flag.

// IntVar定义了一个有指定名字，默认值，和用法说明的int标签。
// 参数p指向一个存储标签解析值的int变量。
func IntVar(p *int, name string, value int, usage string)

// Lookup returns the Flag structure of the named command-line flag,
// returning nil if none exists.

// Lookup返回命令行已经定义过的标签，如果标签不存在的话，返回nil。
func Lookup(name string) *Flag

// NArg is the number of arguments remaining after flags have been processed.

// 在命令行标签被解析之后，NArg就返回解析后参数的个数。
func NArg() int

// NFlag returns the number of command-line flags that have been set.

// NFlag返回解析过的命令行标签的数量。
func NFlag() int

// NewFlagSet returns a new, empty flag set with the specified name and
// error handling property.

// NewFlagSet
// 通过设置一个特定的名字和错误处理属性，返回一个新的，空的FlagSet。
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet

// Parse parses the command-line flags from os.Args[1:].  Must be called
// after all flags are defined and before flags are accessed by the program.

// Parse从参数os.Args[1:]中解析命令行标签。
// 这个方法调用时间点必须在FlagSet的所有标签都定义之后，程序访问这些标签之前。
func Parse()

// Parsed reports whether the command-line flags have been parsed.

// Parsed 返回是否命令行标签已经被解析过。
func Parsed() bool

// PrintDefaults prints, to standard error unless configured otherwise,
// a usage message showing the default settings of all defined
// command-line flags.
// For an integer valued flag x, the default output has the form
//     -x int
//         usage-message-for-x (default 7)
// The usage message will appear on a separate line for anything but
// a bool flag with a one-byte name. For bool flags, the type is
// omitted and if the flag name is one byte the usage message appears
// on the same line. The parenthetical default is omitted if the
// default is the zero value for the type. The listed type, here int,
// can be changed by placing a back-quoted name in the flag's usage
// string; the first such item in the message is taken to be a parameter
// name to show in the message and the back quotes are stripped from
// the message when displayed. For instance, given
//     flag.String("I", "", "search `directory` for include files")
// the output will be
//     -I directory
//         search directory for include files.

// PrintDefaults打印出标准错误，就是所有命令行中定义好的标签的默认信息。
func PrintDefaults()

// Set sets the value of the named command-line flag.

// Set设置命令行中已经定义过的标签的值。
func Set(name, value string) error

// String defines a string flag with specified name, default value, and usage
// string. The return value is the address of a string variable that stores the
// value of the flag.

// String定义了一个有指定名字，默认值，和用法说明的string标签。
// 返回值是一个存储标签解析值的string变量地址。
func String(name string, value string, usage string) *string

// StringVar defines a string flag with specified name, default value, and usage
// string. The argument p points to a string variable in which to store the
// value of the flag.

// StringVar定义了一个有指定名字，默认值，和用法说明的string标签。
// 参数p指向一个存储标签解析值的string变量。
func StringVar(p *string, name string, value string, usage string)

// Uint defines a uint flag with specified name, default value, and usage
// string. The return value is the address of a uint variable that stores the
// value of the flag.

// Uint定义了一个有指定名字，默认值，和用法说明的uint标签。
// 返回值是一个存储标签解析值的uint变量地址。
func Uint(name string, value uint, usage string) *uint

// Uint64 defines a uint64 flag with specified name, default value, and usage
// string. The return value is the address of a uint64 variable that stores the
// value of the flag.

// Uint64定义了一个有指定名字，默认值，和用法说明的uint64标签。
// 返回值是一个存储标签解析值的uint64变量地址。
func Uint64(name string, value uint64, usage string) *uint64

// Uint64Var defines a uint64 flag with specified name, default value, and usage
// string. The argument p points to a uint64 variable in which to store the
// value of the flag.

// Uint64Var定义了一个有指定名字，默认值，和用法说明的uint64标签。
// 参数p指向一个存储标签解析值的uint64变量。
func Uint64Var(p *uint64, name string, value uint64, usage string)

// UintVar defines a uint flag with specified name, default value, and usage
// string. The argument p points to a uint variable in which to store the value
// of the flag.

// UintVar定义了一个有指定名字，默认值，和用法说明的uint标签。
// 参数p指向一个存储标签解析值的uint变量。
func UintVar(p *uint, name string, value uint, usage string)

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.

// Var定义了一个有指定名字和用法说明的标签。标签的类型和值是由第一个参数指定的，
// 这个参数 是Value类型，并且是用户自定义的实现了Value接口的类型。举个例子，调用
// 者可以定义一种标签，这种标签会把 逗号分隔的字符串变成字符串slice，并提供出这
// 种转换的方法。这样，Set（FlagSet）就会将逗号分隔 的字符串转换成为slice。
func Var(value Value, name string, usage string)

// Visit visits the command-line flags in lexicographical order, calling fn
// for each.  It visits only those flags that have been set.

// Visit按照字典顺序遍历命令行标签，并且对每个标签调用fn。
// 这个函数只遍历定义过的标签。
func Visit(fn func(*Flag))

// VisitAll visits the command-line flags in lexicographical order, calling
// fn for each.  It visits all flags, even those not set.

// VisitAll按照字典顺序遍历控制台标签，并且对每个标签调用fn。
// 这个函数会遍历所有标签，包括那些没有定义的标签。
func VisitAll(fn func(*Flag))

// Arg returns the i'th argument.  Arg(0) is the first remaining argument
// after flags have been processed. Arg returns an empty string if the
// requested element does not exist.

// Arg返回第i个参数。当有标签被解析之后，Arg(0)就成为了保留参数。
func (*FlagSet) Arg(i int) string

// Args returns the non-flag arguments.

// Args返回非标签的参数。
func (*FlagSet) Args() []string

// Bool defines a bool flag with specified name, default value, and usage
// string. The return value is the address of a bool variable that stores the
// value of the flag.

// Bool定义了一个有指定名字，默认值，和用法说明的bool标签。
// 返回值是一个存储标签解析值的bool变量地址。
func (*FlagSet) Bool(name string, value bool, usage string) *bool

// BoolVar defines a bool flag with specified name, default value, and usage
// string. The argument p points to a bool variable in which to store the value
// of the flag.

// BoolVar定义了一个有指定名字，默认值，和用法说明的标签。
// 参数p指向一个存储标签值的bool变量。
func (*FlagSet) BoolVar(p *bool, name string, value bool, usage string)

// Duration defines a time.Duration flag with specified name, default value, and
// usage string. The return value is the address of a time.Duration variable
// that stores the value of the flag. The flag accepts a value acceptable to
// time.ParseDuration.

// Duration定义了一个有指定名字，默认值，和用法说明的time.Duration标签。 返回值
// 是一个存储标签解析值的time.Duration变量地址。 此标记接受一个
// time.ParseDuration 可接受的值。
func (*FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration

// DurationVar defines a time.Duration flag with specified name, default value,
// and usage string. The argument p points to a time.Duration variable in which
// to store the value of the flag. The flag accepts a value acceptable to
// time.ParseDuration.

// DurationVar定义了一个有指定名字，默认值，和用法说明的time.Duration标签。 参数
// p指向一个存储标签解析值的time.Duration变量。 此标记接受一个
// time.ParseDuration 可接受的值。
func (*FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string)

// Float64 defines a float64 flag with specified name, default value, and usage
// string. The return value is the address of a float64 variable that stores the
// value of the flag.

// Float64定义了一个有指定名字，默认值，和用法说明的float64标签。
// 返回值是一个存储标签解析值的float64变量地址。
func (*FlagSet) Float64(name string, value float64, usage string) *float64

// Float64Var defines a float64 flag with specified name, default value, and
// usage string. The argument p points to a float64 variable in which to store
// the value of the flag.

// Float64Var定义了一个有指定名字，默认值，和用法说明的float64标签。
// 参数p指向一个存储标签解析值的float64变量。
func (*FlagSet) Float64Var(p *float64, name string, value float64, usage string)

// Init sets the name and error handling property for a flag set.
// By default, the zero FlagSet uses an empty name and the
// ContinueOnError error handling policy.

// Init设置名字和错误处理标签集合的属性。
// 空标签集合默认使用一个空名字和ContinueOnError的错误处理属性。
func (*FlagSet) Init(name string, errorHandling ErrorHandling)

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of
// the flag.

// Int定义了一个有指定名字，默认值，和用法说明的int标签。
// 返回值是一个存储标签解析值的int变量地址。
func (*FlagSet) Int(name string, value int, usage string) *int

// Int64 defines an int64 flag with specified name, default value, and usage
// string. The return value is the address of an int64 variable that stores the
// value of the flag.

// Int64定义了一个有指定名字，默认值，和用法说明的int64标签。
// 返回值是一个存储标签解析值的int64变量地址。
func (*FlagSet) Int64(name string, value int64, usage string) *int64

// Int64Var defines an int64 flag with specified name, default value, and usage
// string. The argument p points to an int64 variable in which to store the
// value of the flag.

// Int64Var定义了一个有指定名字，默认值，和用法说明的int64标签。
// 参数p指向一个存储标签解析值的int64变量。
func (*FlagSet) Int64Var(p *int64, name string, value int64, usage string)

// IntVar defines an int flag with specified name, default value, and usage
// string. The argument p points to an int variable in which to store the value
// of the flag.

// IntVar定义了一个有指定名字，默认值，和用法说明的int标签。
// 参数p指向一个存储标签解析值的int变量。
func (*FlagSet) IntVar(p *int, name string, value int, usage string)

// Lookup returns the Flag structure of the named flag, returning nil if none
// exists.

// Lookup返回已经定义过的标签，如果标签不存在的话，返回nil。
func (*FlagSet) Lookup(name string) *Flag

// NArg is the number of arguments remaining after flags have been processed.

// 在标签被解析之后，NArg就返回解析后参数的个数。
func (*FlagSet) NArg() int

// NFlag returns the number of flags that have been set.

// NFlag返回解析过的标签的数量。
func (*FlagSet) NFlag() int

// Parse parses flag definitions from the argument list, which should not
// include the command name.  Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help or -h were set but not defined.

// Parse从参数列表中解析定义的标签，这个参数列表并不包含执行的命令名字。 这个方
// 法调用时间点必须在FlagSet的所有标签都定义之后，程序访问这些标签之前。 当
// -help 或 -h 标签没有定义却被调用了的时候，这个方法返回 ErrHelp。
func (*FlagSet) Parse(arguments []string) error

// Parsed reports whether f.Parse has been called.

// Parsed返回是否f.Parse已经被调用过。
func (*FlagSet) Parsed() bool

// PrintDefaults prints to standard error the default values of all
// defined command-line flags in the set. See the documentation for
// the global function PrintDefaults for more information.

// 除非有特别配置，否则PrintDefault会将内容输出到标准输出控制台中。
// PrintDefault会输出集合中所有定义好的标签的默认信息
func (*FlagSet) PrintDefaults()

// Set sets the value of the named flag.

// Set设置定义过的标签的值
func (*FlagSet) Set(name, value string) error

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.

// SetOutput设置了用法和错误信息的输出目的地。
// 如果output是nil，输出目的地就会使用os.Stderr。
func (*FlagSet) SetOutput(output io.Writer)

// String defines a string flag with specified name, default value, and usage
// string. The return value is the address of a string variable that stores the
// value of the flag.

// String定义了一个有指定名字，默认值，和用法说明的string标签。
// 返回值是一个存储标签解析值的string变量地址。
func (*FlagSet) String(name string, value string, usage string) *string

// StringVar defines a string flag with specified name, default value, and usage
// string. The argument p points to a string variable in which to store the
// value of the flag.

// StringVar定义了一个有指定名字，默认值，和用法说明的string标签。
// 参数p指向一个存储标签解析值的string变量。
func (*FlagSet) StringVar(p *string, name string, value string, usage string)

// Uint defines a uint flag with specified name, default value, and usage
// string. The return value is the address of a uint variable that stores the
// value of the flag.

// Uint定义了一个有指定名字，默认值，和用法说明的uint标签。
// 返回值是一个存储标签解析值的uint变量地址。
func (*FlagSet) Uint(name string, value uint, usage string) *uint

// Uint64 defines a uint64 flag with specified name, default value, and usage
// string. The return value is the address of a uint64 variable that stores the
// value of the flag.

// Uint64定义了一个有指定名字，默认值，和用法说明的uint64标签。
// 返回值是一个存储标签解析值的uint64变量地址。
func (*FlagSet) Uint64(name string, value uint64, usage string) *uint64

// Uint64Var defines a uint64 flag with specified name, default value, and usage
// string. The argument p points to a uint64 variable in which to store the
// value of the flag.

// Uint64Var定义了一个有指定名字，默认值，和用法说明的uint64标签。
// 参数p指向一个存储标签解析值的uint64变量。
func (*FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string)

// UintVar defines a uint flag with specified name, default value, and usage
// string. The argument p points to a uint variable in which to store the value
// of the flag.

// UintVar定义了一个有指定名字，默认值，和用法说明的uint标签。
// 参数p指向一个存储标签解析值的uint变量。
func (*FlagSet) UintVar(p *uint, name string, value uint, usage string)

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.

// Var定义了一个有指定名字和用法说明的标签。标签的类型和值是由第一个参数指定的，
// 这个参数 是Value类型，并且是用户自定义的实现了Value接口的类型。举个例子，调用
// 者可以定义一种标签，这种标签会把 逗号分隔的字符串变成字符串slice，并提供出这
// 种转换的方法。这样，Set（FlagSet）就会将逗号分隔 的字符串转换成为slice。
func (*FlagSet) Var(value Value, name string, usage string)

// Visit visits the flags in lexicographical order, calling fn for each.
// It visits only those flags that have been set.

// Visit按照字典顺序遍历标签，并且对每个标签调用fn。
// 这个函数只遍历定义过的标签。
func (*FlagSet) Visit(fn func(*Flag))

// VisitAll visits the flags in lexicographical order, calling fn for each.
// It visits all flags, even those not set.

// VisitAll按照字典顺序遍历标签，并且对每个标签调用fn。
// 这个函数会遍历所有标签，包括那些没有定义的标签。
func (*FlagSet) VisitAll(fn func(*Flag))

