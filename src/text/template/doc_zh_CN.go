// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package template implements data-driven templates for generating textual output.
//
// To generate HTML output, see package html/template, which has the same interface
// as this package but automatically secures HTML output against certain attacks.
//
// Templates are executed by applying them to a data structure. Annotations in the
// template refer to elements of the data structure (typically a field of a struct
// or a key in a map) to control execution and derive values to be displayed.
// Execution of the template walks the structure and sets the cursor, represented
// by a period '.' and called "dot", to the value at the current location in the
// structure as execution proceeds.
//
// The input text for a template is UTF-8-encoded text in any format.
// "Actions"--data evaluations or control structures--are delimited by "{{" and
// "}}"; all text outside actions is copied to the output unchanged. Actions may
// not span newlines, although comments can.
//
// Once parsed, a template may be executed safely in parallel.
//
// Here is a trivial example that prints "17 items are made of wool".
//
//	type Inventory struct {
//		Material string
//		Count    uint
//	}
//	sweaters := Inventory{"wool", 17}
//	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
//	if err != nil { panic(err) }
//	err = tmpl.Execute(os.Stdout, sweaters)
//	if err != nil { panic(err) }
//
// More intricate examples appear below.
//
//
// Actions
//
// Here is the list of actions. "Arguments" and "pipelines" are evaluations of
// data, defined in detail below.
//
//	{{/* a comment */}}
//		A comment; discarded. May contain newlines.
//		Comments do not nest and must start and end at the
//		delimiters, as shown here.
//
//	{{pipeline}}
//		The default textual representation of the value of the pipeline
//		is copied to the output.
//
//	{{if pipeline}} T1 {{end}}
//		If the value of the pipeline is empty, no output is generated;
//		otherwise, T1 is executed.  The empty values are false, 0, any
//		nil pointer or interface value, and any array, slice, map, or
//		string of length zero.
//		Dot is unaffected.
//
//	{{if pipeline}} T1 {{else}} T0 {{end}}
//		If the value of the pipeline is empty, T0 is executed;
//		otherwise, T1 is executed.  Dot is unaffected.
//
//	{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}
//		To simplify the appearance of if-else chains, the else action
//		of an if may include another if directly; the effect is exactly
//		the same as writing
//			{{if pipeline}} T1 {{else}}{{if pipeline}} T0 {{end}}{{end}}
//
//	{{range pipeline}} T1 {{end}}
//		The value of the pipeline must be an array, slice, map, or channel.
//		If the value of the pipeline has length zero, nothing is output;
//		otherwise, dot is set to the successive elements of the array,
//		slice, or map and T1 is executed. If the value is a map and the
//		keys are of basic type with a defined order ("comparable"), the
//		elements will be visited in sorted key order.
//
//	{{range pipeline}} T1 {{else}} T0 {{end}}
//		The value of the pipeline must be an array, slice, map, or channel.
//		If the value of the pipeline has length zero, dot is unaffected and
//		T0 is executed; otherwise, dot is set to the successive elements
//		of the array, slice, or map and T1 is executed.
//
//	{{template "name"}}
//		The template with the specified name is executed with nil data.
//
//	{{template "name" pipeline}}
//		The template with the specified name is executed with dot set
//		to the value of the pipeline.
//
//	{{with pipeline}} T1 {{end}}
//		If the value of the pipeline is empty, no output is generated;
//		otherwise, dot is set to the value of the pipeline and T1 is
//		executed.
//
//	{{with pipeline}} T1 {{else}} T0 {{end}}
//		If the value of the pipeline is empty, dot is unaffected and T0
//		is executed; otherwise, dot is set to the value of the pipeline
//		and T1 is executed.
//
//
// Arguments
//
// An argument is a simple value, denoted by one of the following.
//
//	- A boolean, string, character, integer, floating-point, imaginary
//	  or complex constant in Go syntax. These behave like Go's untyped
//	  constants, although raw strings may not span newlines.
//	- The keyword nil, representing an untyped Go nil.
//	- The character '.' (period):
//		.
//	  The result is the value of dot.
//	- A variable name, which is a (possibly empty) alphanumeric string
//	  preceded by a dollar sign, such as
//		$piOver2
//	  or
//		$
//	  The result is the value of the variable.
//	  Variables are described below.
//	- The name of a field of the data, which must be a struct, preceded
//	  by a period, such as
//		.Field
//	  The result is the value of the field. Field invocations may be
//	  chained:
//	    .Field1.Field2
//	  Fields can also be evaluated on variables, including chaining:
//	    $x.Field1.Field2
//	- The name of a key of the data, which must be a map, preceded
//	  by a period, such as
//		.Key
//	  The result is the map element value indexed by the key.
//	  Key invocations may be chained and combined with fields to any
//	  depth:
//	    .Field1.Key1.Field2.Key2
//	  Although the key must be an alphanumeric identifier, unlike with
//	  field names they do not need to start with an upper case letter.
//	  Keys can also be evaluated on variables, including chaining:
//	    $x.key1.key2
//	- The name of a niladic method of the data, preceded by a period,
//	  such as
//		.Method
//	  The result is the value of invoking the method with dot as the
//	  receiver, dot.Method(). Such a method must have one return value (of
//	  any type) or two return values, the second of which is an error.
//	  If it has two and the returned error is non-nil, execution terminates
//	  and an error is returned to the caller as the value of Execute.
//	  Method invocations may be chained and combined with fields and keys
//	  to any depth:
//	    .Field1.Key1.Method1.Field2.Key2.Method2
//	  Methods can also be evaluated on variables, including chaining:
//	    $x.Method1.Field
//	- The name of a niladic function, such as
//		fun
//	  The result is the value of invoking the function, fun(). The return
//	  types and values behave as in methods. Functions and function
//	  names are described below.
//	- A parenthesized instance of one the above, for grouping. The result
//	  may be accessed by a field or map key invocation.
//		print (.F1 arg1) (.F2 arg2)
//		(.StructValuedMethod "arg").Field
//
// Arguments may evaluate to any type; if they are pointers the implementation
// automatically indirects to the base type when required. If an evaluation yields
// a function value, such as a function-valued field of a struct, the function is
// not invoked automatically, but it can be used as a truth value for an if action
// and the like. To invoke it, use the call function, defined below.
//
// A pipeline is a possibly chained sequence of "commands". A command is a simple
// value (argument) or a function or method call, possibly with multiple arguments:
//
//	Argument
//		The result is the value of evaluating the argument.
//	.Method [Argument...]
//		The method can be alone or the last element of a chain but,
//		unlike methods in the middle of a chain, it can take arguments.
//		The result is the value of calling the method with the
//		arguments:
//			dot.Method(Argument1, etc.)
//	functionName [Argument...]
//		The result is the value of calling the function associated
//		with the name:
//			function(Argument1, etc.)
//		Functions and function names are described below.
//
//
// Pipelines
//
// A pipeline may be "chained" by separating a sequence of commands with pipeline
// characters '|'. In a chained pipeline, the result of the each command is passed
// as the last argument of the following command. The output of the final command
// in the pipeline is the value of the pipeline.
//
// The output of a command will be either one value or two values, the second of
// which has type error. If that second value is present and evaluates to non-nil,
// execution terminates and the error is returned to the caller of Execute.
//
//
// Variables
//
// A pipeline inside an action may initialize a variable to capture the result. The
// initialization has syntax
//
//	$variable := pipeline
//
// where $variable is the name of the variable. An action that declares a variable
// produces no output.
//
// If a "range" action initializes a variable, the variable is set to the
// successive elements of the iteration. Also, a "range" may declare two variables,
// separated by a comma:
//
//	range $index, $element := pipeline
//
// in which case $index and $element are set to the successive values of the
// array/slice index or map key and element, respectively. Note that if there is
// only one variable, it is assigned the element; this is opposite to the
// convention in Go range clauses.
//
// A variable's scope extends to the "end" action of the control structure ("if",
// "with", or "range") in which it is declared, or to the end of the template if
// there is no such control structure. A template invocation does not inherit
// variables from the point of its invocation.
//
// When execution begins, $ is set to the data argument passed to Execute, that is,
// to the starting value of dot.
//
//
// Examples
//
// Here are some example one-line templates demonstrating pipelines and variables.
// All produce the quoted word "output":
//
//	{{"\"output\""}}
//		A string constant.
//	{{`"output"`}}
//		A raw string constant.
//	{{printf "%q" "output"}}
//		A function call.
//	{{"output" | printf "%q"}}
//		A function call whose final argument comes from the previous
//		command.
//	{{printf "%q" (print "out" "put")}}
//		A parenthesized argument.
//	{{"put" | printf "%s%s" "out" | printf "%q"}}
//		A more elaborate call.
//	{{"output" | printf "%s" | printf "%q"}}
//		A longer chain.
//	{{with "output"}}{{printf "%q" .}}{{end}}
//		A with action using dot.
//	{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
//		A with action that creates and uses a variable.
//	{{with $x := "output"}}{{printf "%q" $x}}{{end}}
//		A with action that uses the variable in another action.
//	{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
//		The same, but pipelined.
//
//
// Functions
//
// During execution functions are found in two function maps: first in the
// template, then in the global function map. By default, no functions are defined
// in the template but the Funcs method can be used to add them.
//
// Predefined global functions are named as follows.
//
//	and
//		Returns the boolean AND of its arguments by returning the
//		first empty argument or the last argument, that is,
//		"and x y" behaves as "if x then y else x". All the
//		arguments are evaluated.
//	call
//		Returns the result of calling the first argument, which
//		must be a function, with the remaining arguments as parameters.
//		Thus "call .X.Y 1 2" is, in Go notation, dot.X.Y(1, 2) where
//		Y is a func-valued field, map entry, or the like.
//		The first argument must be the result of an evaluation
//		that yields a value of function type (as distinct from
//		a predefined function such as print). The function must
//		return either one or two result values, the second of which
//		is of type error. If the arguments don't match the function
//		or the returned error value is non-nil, execution stops.
//	html
//		Returns the escaped HTML equivalent of the textual
//		representation of its arguments.
//	index
//		Returns the result of indexing its first argument by the
//		following arguments. Thus "index x 1 2 3" is, in Go syntax,
//		x[1][2][3]. Each indexed item must be a map, slice, or array.
//	js
//		Returns the escaped JavaScript equivalent of the textual
//		representation of its arguments.
//	len
//		Returns the integer length of its argument.
//	not
//		Returns the boolean negation of its single argument.
//	or
//		Returns the boolean OR of its arguments by returning the
//		first non-empty argument or the last argument, that is,
//		"or x y" behaves as "if x then x else y". All the
//		arguments are evaluated.
//	print
//		An alias for fmt.Sprint
//	printf
//		An alias for fmt.Sprintf
//	println
//		An alias for fmt.Sprintln
//	urlquery
//		Returns the escaped value of the textual representation of
//		its arguments in a form suitable for embedding in a URL query.
//
// The boolean functions take any zero value to be false and a non-zero value to be
// true.
//
// There is also a set of binary comparison operators defined as functions:
//
//	eq
//		Returns the boolean truth of arg1 == arg2
//	ne
//		Returns the boolean truth of arg1 != arg2
//	lt
//		Returns the boolean truth of arg1 < arg2
//	le
//		Returns the boolean truth of arg1 <= arg2
//	gt
//		Returns the boolean truth of arg1 > arg2
//	ge
//		Returns the boolean truth of arg1 >= arg2
//
// For simpler multi-way equality tests, eq (only) accepts two or more arguments
// and compares the second and subsequent to the first, returning in effect
//
//	arg1==arg2 || arg1==arg3 || arg1==arg4 ...
//
// (Unlike with || in Go, however, eq is a function call and all the arguments will
// be evaluated.)
//
// The comparison functions work on basic types only (or named basic types, such as
// "type Celsius float32"). They implement the Go rules for comparison of values,
// except that size and exact type are ignored, so any integer value, signed or
// unsigned, may be compared with any other integer value. (The arithmetic value is
// compared, not the bit pattern, so all negative integers are less than all
// unsigned integers.) However, as usual, one may not compare an int with a float32
// and so on.
//
//
// Associated templates
//
// Each template is named by a string specified when it is created. Also, each
// template is associated with zero or more other templates that it may invoke by
// name; such associations are transitive and form a name space of templates.
//
// A template may use a template invocation to instantiate another associated
// template; see the explanation of the "template" action above. The name must be
// that of a template associated with the template that contains the invocation.
//
//
// Nested template definitions
//
// When parsing a template, another template may be defined and associated with the
// template being parsed. Template definitions must appear at the top level of the
// template, much like global variables in a Go program.
//
// The syntax of such definitions is to surround each template declaration with a
// "define" and "end" action.
//
// The define action names the template being created by providing a string
// constant. Here is a simple example:
//
//	`{{define "T1"}}ONE{{end}}
//	{{define "T2"}}TWO{{end}}
//	{{define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
//	{{template "T3"}}`
//
// This defines two templates, T1 and T2, and a third T3 that invokes the other two
// when it is executed. Finally it invokes T3. If executed this template will
// produce the text
//
//	ONE TWO
//
// By construction, a template may reside in only one association. If it's
// necessary to have a template addressable from multiple associations, the
// template definition must be parsed multiple times to create distinct *Template
// values, or must be copied with the Clone or AddParseTree method.
//
// Parse may be called multiple times to assemble the various associated templates;
// see the ParseFiles and ParseGlob functions and methods for simple ways to parse
// related templates stored in files.
//
// A template may be executed directly or through ExecuteTemplate, which executes
// an associated template identified by name. To invoke our example above, we might
// write,
//
//	err := tmpl.Execute(os.Stdout, "no data needed")
//	if err != nil {
//		log.Fatalf("execution failed: %s", err)
//	}
//
// or to invoke a particular template explicitly by name,
//
//	err := tmpl.ExecuteTemplate(os.Stdout, "T2", "no data needed")
//	if err != nil {
//		log.Fatalf("execution failed: %s", err)
//	}

// template包实现了数据驱动的用于生成文本输出的模板。
//
// 如果要生成HTML格式的输出，参见html/template包，该包提供了和本包相同的接口，但会自动将输出转化为安全的HTML格式输出，可以抵抗一些网络攻击。
//
// 通过将模板应用于一个数据结构（即该数据结构作为模板的参数）来执行，来获得输出。模板中的注释引用数据接口的元素（一般如结构体的字段或者字典的键）来控制执行过程和获取需要呈现的值。模板执行时会遍历结构并将指针表示为'.'（称之为"dot"）指向运行过程中数据结构的当前位置的值。
//
// 用作模板的输入文本必须是utf-8编码的文本。"Action"—数据运算和控制单位—由"{{"和"}}"界定；在Action之外的所有文本都不做修改的拷贝到输出中。Action内部不能有换行，但注释可以有换行。
//
// 经解析生成模板后，一个模板可以安全的并发执行。
//
// 下面是一个简单的例子，可以打印"17 of wool"。
//
//	type Inventory struct {
//		Material string
//		Count    uint
//	}
//	sweaters := Inventory{"wool", 17}
//	tmpl, err := template.New("test").Parse("{{.Count}} of {{.Material}}")
//	if err != nil { panic(err) }
//	err = tmpl.Execute(os.Stdout, sweaters)
//	if err != nil { panic(err) }
//
// 更复杂的例子在下面。
package template

// HTMLEscape writes to w the escaped HTML equivalent of the plain text data b.

// 函数向w中写入b的HTML转义等价表示。
func HTMLEscape(w io.Writer, b []byte)

// HTMLEscapeString returns the escaped HTML equivalent of the plain text data s.

// 返回s的HTML转义等价表示字符串。
func HTMLEscapeString(s string) string

// HTMLEscaper returns the escaped HTML equivalent of the textual representation of
// its arguments.

// 函数返回其所有参数文本表示的HTML转义等价表示字符串。
func HTMLEscaper(args ...interface{}) string

// JSEscape writes to w the escaped JavaScript equivalent of the plain text data b.

// 函数向w中写入b的JavaScript转义等价表示。
func JSEscape(w io.Writer, b []byte)

// JSEscapeString returns the escaped JavaScript equivalent of the plain text data
// s.

// 返回s的JavaScript转义等价表示字符串。
func JSEscapeString(s string) string

// JSEscaper returns the escaped JavaScript equivalent of the textual
// representation of its arguments.

// 函数返回其所有参数文本表示的JavaScript转义等价表示字符串。
func JSEscaper(args ...interface{}) string

// URLQueryEscaper returns the escaped value of the textual representation of its
// arguments in a form suitable for embedding in a URL query.

// 函数返回其所有参数文本表示的可以嵌入URL查询的转义等价表示字符串。
func URLQueryEscaper(args ...interface{}) string

// FuncMap is the type of the map defining the mapping from names to functions.
// Each function must have either a single return value, or two return values of
// which the second has type error. In that case, if the second (error) return
// value evaluates to non-nil during execution, execution terminates and Execute
// returns that error.

// FuncMap类型定义了函数名字符串到函数的映射，每个函数都必须有1到2个返回值，如果有2个则后一个必须是error接口类型；如果有2个返回值的方法返回的error非nil，模板执行会中断并返回给调用者该错误。
type FuncMap map[string]interface{}

// Template is the representation of a parsed template. The *parse.Tree field is
// exported only for use by html/template and should be treated as unexported by
// all other clients.

// 代表一个解析好的模板，*parse.Tree字段仅仅是暴露给html/template包使用的，因此其他人应该视字段未导出。
type Template struct {
	*parse.Tree
	// contains filtered or unexported fields
}

// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//
//	var t = template.Must(template.New("name").Parse("text"))

//函数帮助包装一个调用给方法的返回值(模板 ，错误)，如果错误不为空将打印输出。该函数用于变量，像这样初始化
//
//var t = template.Must(template.New("name").Parse("text"))
func Must(t *Template, err error) *Template

// New allocates a new template with the given name.

// 函数给新模板一个名字
func New(name string) *Template

// ParseFiles creates a new Template and parses the template definitions from the
// named files. The returned template's name will have the (base) name and (parsed)
// contents of the first file. There must be at least one file. If an error occurs,
// parsing stops and the returned *Template is nil.

//函数创建一个新模板，从命名的文件中解析模板的定义。返回模板的名字和模板的内容基于第一个文件，必须至少一个文件
//如果发生错误，将会停止解析返回的模板也将使空的
func ParseFiles(filenames ...string) (*Template, error)

// ParseGlob creates a new Template and parses the template definitions from the
// files identified by the pattern, which must match at least one file. The
// returned template will have the (base) name and (parsed) contents of the first
// file matched by the pattern. ParseGlob is equivalent to calling ParseFiles with
// the list of files matched by the pattern.

//函数创建一个新的模板，从按模式识别的文件解析模板定义，必须至少包含一个文件。返回模板的名字和模板的内容基于按第一个
//模式识别的文件。函数相当于使用 ParseFiles 调用按模式匹配的的文件列表。
func ParseGlob(pattern string) (*Template, error)

// AddParseTree creates a new template with the name and parse tree and associates
// it with t.

//函数创建一个新的模板，通过名字和解析树和模板联系在一起
func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error)

// Clone returns a duplicate of the template, including all associated templates.
// The actual representation is not copied, but the name space of associated
// templates is, so further calls to Parse in the copy will add templates to the
// copy but not to the original. Clone can be used to prepare common templates and
// use them with variant definitions for other templates by adding the variants
// after the clone is made.

//函数返回一个重复的模板，包含所有相关的模板。实际上不是复制,但命名空间和与模板相关联，所以
//将来在副本调用解析时，添加的模板是到副本上而不是原始的模板。函数常用于准备共同的模板，克隆之后用于给
//其他模板定义变量添加变量
func (t *Template) Clone() (*Template, error)

// Delims sets the action delimiters to the specified strings, to be used in
// subsequent calls to Parse, ParseFiles, or ParseGlob. Nested template definitions
// will inherit the settings. An empty delimiter stands for the corresponding
// default: {{ or }}. The return value is the template, so calls can be chained.

//函数将模板使用的分隔符设置为指定的字符串，在后来调用 Parse, ParseFiles, or ParseGlob 时使用
//嵌套模板定义将会继承这些设置。默认 {{ or }}。返回值类型是模板，所以调用会被连接。
func (t *Template) Delims(left, right string) *Template

// Execute applies a parsed template to the specified data object, and writes the
// output to wr. If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to the output
// writer. A template may be executed safely in parallel.

//函数给模板解析指定的数据对象，如果只需模板或者输出包含错误，执行操作停止。但是部分输出结果
//可能已经写到输出者那里了。一个模板可以安全的并发执行。
func (t *Template) Execute(wr io.Writer, data interface{}) (err error)

// ExecuteTemplate applies the template associated with t that has the given name
// to the specified data object and writes the output to wr. If an error occurs
// executing the template or writing its output, execution stops, but partial
// results may already have been written to the output writer. A template may be
// executed safely in parallel.

//函数用 t 应用该模板，给指定的数据类命名，写出到 wr 。如果执行模板或者写出时发生错误，执行
//操作将会停止。但是部分输出结果可能已经写出输出者那里了。一个模板可以安全的并发执行。
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error

// Funcs adds the elements of the argument map to the template's function map. It
// panics if a value in the map is not a function with appropriate return type.
// However, it is legal to overwrite elements of the map. The return value is the
// template, so calls can be chained.

//函数添加参数的map的元素到模板方法的map。如果map中的值不是一个带有适当返回值的方法将发生错误。
//然而，覆盖map中的元素是合法的。返回值是模板，调用将被连接。
func (t *Template) Funcs(funcMap FuncMap) *Template

// Lookup returns the template with the given name that is associated with t, or
// nil if there is no such template.

//Lookup方法返回与t关联的名为name的模板，nil 说明没有这样的模板。
func (t *Template) Lookup(name string) *Template

// Name returns the name of the template.

//函数返回模板的名字
func (t *Template) Name() string

// New allocates a new template associated with the given one and with the same
// delimiters. The association, which is transitive, allows one template to invoke
// another with a {{template}} action.

//函数使用给定一个相同的分隔符的方式分配一个新的模板。关联是传递的，允许一个模板通过 {{模板}} 调用另一个模板
func (t *Template) New(name string) *Template

// Parse parses a string into a template. Nested template definitions will be
// associated with the top-level template t. Parse may be called multiple times to
// parse definitions of templates to associate with t. It is an error if a
// resulting template is non-empty (contains content other than template
// definitions) and would replace a non-empty template with the same name. (In
// multiple calls to Parse with the same receiver template, only one call can
// contain text other than space, comments, and template definitions.)

//Parse方法将字符串text解析为模板。嵌套定义的模板会关联到最顶层的t。
//Parse可以多次调用，但只有第一次调用可以包含空格、注释和模板定义之外的文本。
//如果后面的调用在解析后仍剩余文本会引发错误、返回nil且丢弃剩余文本；如果解析得到的模板已有相关联的同名模板，会覆盖掉原模板。

func (t *Template) Parse(text string) (*Template, error)

// ParseFiles parses the named files and associates the resulting templates with t.
// If an error occurs, parsing stops and the returned template is nil; otherwise it
// is t. There must be at least one file.

//函数解析filenames指定的文件里的模板定义并将解析结果与t关联。如果产生错误，解析将会停止，返回的模板是nil。否则
//返回模板t 。至少要提供一个文件。
func (t *Template) ParseFiles(filenames ...string) (*Template, error)

// ParseGlob parses the template definitions in the files identified by the pattern
// and associates the resulting templates with t. The pattern is processed by
// filepath.Glob and must match at least one file. ParseGlob is equivalent to
// calling t.ParseFiles with the list of files matched by the pattern.

//函数解析匹配pattern的文件里的模板定义并将解析结果与t关联。至少要提供一个文件
func (t *Template) ParseGlob(pattern string) (*Template, error)

// Templates returns a slice of the templates associated with t, including t
// itself.

//函数返回与t相关联的模板的切片，包括模板t的本身
func (t *Template) Templates() []*Template
