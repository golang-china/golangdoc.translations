// +build ingore

// Package template (html/template) implements data-driven templates for
// generating HTML output safe against code injection. It provides the same
// interface as package text/template and should be used instead of
// text/template whenever the output is HTML.
//
// The documentation here focuses on the security features of the package. For
// information about how to program the templates themselves, see the
// documentation for text/template.
//
//
// Introduction
//
// This package wraps package text/template so you can share its template API to
// parse and execute HTML templates safely.
//
// 	tmpl, err := template.New("name").Parse(...)
// 	// Error checking elided
// 	err = tmpl.Execute(out, data)
//
// If successful, tmpl will now be injection-safe. Otherwise, err is an error
// defined in the docs for ErrorCode.
//
// HTML templates treat data values as plain text which should be encoded so
// they can be safely embedded in an HTML document. The escaping is contextual,
// so actions can appear within JavaScript, CSS, and URI contexts.
//
// The security model used by this package assumes that template authors are
// trusted, while Execute's data parameter is not. More details are provided
// below.
//
// Example
//
// 	import "text/template"
// 	...
// 	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
// 	err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")
//
// produces
//
// 	Hello, <script>alert('you have been pwned')</script>!
//
// but the contextual autoescaping in html/template
//
// 	import "html/template"
// 	...
// 	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
// 	err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")
//
// produces safe, escaped HTML output
//
// 	Hello, &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;!
//
//
// Contexts
//
// This package understands HTML, CSS, JavaScript, and URIs. It adds sanitizing
// functions to each simple action pipeline, so given the excerpt
//
// 	<a href="/search?q={{.}}">{{.}}</a>
//
// At parse time each {{.}} is overwritten to add escaping functions as
// necessary. In this case it becomes
//
// 	<a href="/search?q={{. | urlquery}}">{{. | html}}</a>
//
//
// Errors
//
// See the documentation of ErrorCode for details.
//
//
// A fuller picture
//
// The rest of this package comment may be skipped on first reading; it includes
// details necessary to understand escaping contexts and error messages. Most
// users will not need to understand these details.
//
//
// Contexts
//
// Assuming {{.}} is `O'Reilly: How are <i>you</i>?`, the table below shows how
// {{.}} appears when used in the context to the left.
//
// 	Context                          {{.}} After
// 	{{.}}                            O'Reilly: How are &lt;i&gt;you&lt;/i&gt;?
// 	<a title='{{.}}'>                O&#39;Reilly: How are you?
// 	<a href="/{{.}}">                O&#39;Reilly: How are %3ci%3eyou%3c/i%3e?
// 	<a href="?q={{.}}">              O&#39;Reilly%3a%20How%20are%3ci%3e...%3f
// 	<a onx='f("{{.}}")'>             O\x27Reilly: How are \x3ci\x3eyou...?
// 	<a onx='f({{.}})'>               "O\x27Reilly: How are \x3ci\x3eyou...?"
// 	<a onx='pattern = /{{.}}/;'>     O\x27Reilly: How are \x3ci\x3eyou...\x3f
//
// If used in an unsafe context, then the value might be filtered out:
//
// 	Context                          {{.}} After
// 	<a href="{{.}}">                 #ZgotmplZ
//
// since "O'Reilly:" is not an allowed protocol like "http:".
//
// If {{.}} is the innocuous word, `left`, then it can appear more widely,
//
// 	Context                              {{.}} After
// 	{{.}}                                left
// 	<a title='{{.}}'>                    left
// 	<a href='{{.}}'>                     left
// 	<a href='/{{.}}'>                    left
// 	<a href='?dir={{.}}'>                left
// 	<a style="border-{{.}}: 4px">        left
// 	<a style="align: {{.}}">             left
// 	<a style="background: '{{.}}'>       left
// 	<a style="background: url('{{.}}')>  left
// 	<style>p.{{.}} {color:red}</style>   left
//
// Non-string values can be used in JavaScript contexts. If {{.}} is
//
// 	struct{A,B string}{ "foo", "bar" }
//
// in the escaped template
//
// 	<script>var pair = {{.}};</script>
//
// then the template output is
//
// 	<script>var pair = {"A": "foo", "B": "bar"};</script>
//
// See package json to understand how non-string content is marshalled for
// embedding in JavaScript contexts.
//
//
// Typed Strings
//
// By default, this package assumes that all pipelines produce a plain text
// string. It adds escaping pipeline stages necessary to correctly and safely
// embed that plain text string in the appropriate context.
//
// When a data value is not plain text, you can make sure it is not over-escaped
// by marking it with its type.
//
// Types HTML, JS, URL, and others from content.go can carry safe content that
// is exempted from escaping.
//
// The template
//
// 	Hello, {{.}}!
//
// can be invoked with
//
// 	tmpl.Execute(out, template.HTML(`<b>World</b>`))
//
// to produce
//
// 	Hello, <b>World</b>!
//
// instead of the
//
// 	Hello, &lt;b&gt;World&lt;b&gt;!
//
// that would have been produced if {{.}} was a regular string.
//
//
// Security Model
//
// https://rawgit.com/mikesamuel/sanitized-jquery-templates/trunk/safetemplate.html#problem_definition
// defines "safe" as used by this package.
//
// This package assumes that template authors are trusted, that Execute's data
// parameter is not, and seeks to preserve the properties below in the face of
// untrusted data:
//
// Structure Preservation Property: "... when a template author writes an HTML
// tag in a safe templating language, the browser will interpret the
// corresponding portion of the output as a tag regardless of the values of
// untrusted data, and similarly for other structures such as attribute
// boundaries and JS and CSS string boundaries."
//
// Code Effect Property: "... only code specified by the template author should
// run as a result of injecting the template output into a page and all code
// specified by the template author should run as a result of the same."
//
// Least Surprise Property: "A developer (or code reviewer) familiar with HTML,
// CSS, and JavaScript, who knows that contextual autoescaping happens should be
// able to look at a {{.}} and correctly infer what sanitization happens."

// Package template (html/template) implements data-driven templates for
// generating HTML output safe against code injection. It provides the same
// interface as package text/template and should be used instead of
// text/template whenever the output is HTML.
//
// The documentation here focuses on the security features of the package. For
// information about how to program the templates themselves, see the
// documentation for text/template.
//
//
// Introduction
//
// This package wraps package text/template so you can share its template API to
// parse and execute HTML templates safely.
//
// 	tmpl, err := template.New("name").Parse(...)
// 	// Error checking elided
// 	err = tmpl.Execute(out, data)
//
// If successful, tmpl will now be injection-safe. Otherwise, err is an error
// defined in the docs for ErrorCode.
//
// HTML templates treat data values as plain text which should be encoded so
// they can be safely embedded in an HTML document. The escaping is contextual,
// so actions can appear within JavaScript, CSS, and URI contexts.
//
// The security model used by this package assumes that template authors are
// trusted, while Execute's data parameter is not. More details are provided
// below.
//
// Example
//
// 	import "text/template"
// 	...
// 	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
// 	err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")
//
// produces
//
// 	Hello, <script>alert('you have been pwned')</script>!
//
// but the contextual autoescaping in html/template
//
// 	import "html/template"
// 	...
// 	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
// 	err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")
//
// produces safe, escaped HTML output
//
// 	Hello, &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;!
//
//
// Contexts
//
// This package understands HTML, CSS, JavaScript, and URIs. It adds sanitizing
// functions to each simple action pipeline, so given the excerpt
//
// 	<a href="/search?q={{.}}">{{.}}</a>
//
// At parse time each {{.}} is overwritten to add escaping functions as
// necessary. In this case it becomes
//
// 	<a href="/search?q={{. | urlquery}}">{{. | html}}</a>
//
//
// Errors
//
// See the documentation of ErrorCode for details.
//
//
// A fuller picture
//
// The rest of this package comment may be skipped on first reading; it includes
// details necessary to understand escaping contexts and error messages. Most
// users will not need to understand these details.
//
//
// Contexts
//
// Assuming {{.}} is `O'Reilly: How are <i>you</i>?`, the table below shows how
// {{.}} appears when used in the context to the left.
//
// 	Context                          {{.}} After
// 	{{.}}                            O'Reilly: How are &lt;i&gt;you&lt;/i&gt;?
// 	<a title='{{.}}'>                O&#39;Reilly: How are you?
// 	<a href="/{{.}}">                O&#39;Reilly: How are %3ci%3eyou%3c/i%3e?
// 	<a href="?q={{.}}">              O&#39;Reilly%3a%20How%20are%3ci%3e...%3f
// 	<a onx='f("{{.}}")'>             O\x27Reilly: How are \x3ci\x3eyou...?
// 	<a onx='f({{.}})'>               "O\x27Reilly: How are \x3ci\x3eyou...?"
// 	<a onx='pattern = /{{.}}/;'>     O\x27Reilly: How are \x3ci\x3eyou...\x3f
//
// If used in an unsafe context, then the value might be filtered out:
//
// 	Context                          {{.}} After
// 	<a href="{{.}}">                 #ZgotmplZ
//
// since "O'Reilly:" is not an allowed protocol like "http:".
//
// If {{.}} is the innocuous word, `left`, then it can appear more widely,
//
// 	Context                              {{.}} After
// 	{{.}}                                left
// 	<a title='{{.}}'>                    left
// 	<a href='{{.}}'>                     left
// 	<a href='/{{.}}'>                    left
// 	<a href='?dir={{.}}'>                left
// 	<a style="border-{{.}}: 4px">        left
// 	<a style="align: {{.}}">             left
// 	<a style="background: '{{.}}'>       left
// 	<a style="background: url('{{.}}')>  left
// 	<style>p.{{.}} {color:red}</style>   left
//
// Non-string values can be used in JavaScript contexts. If {{.}} is
//
// 	struct{A,B string}{ "foo", "bar" }
//
// in the escaped template
//
// 	<script>var pair = {{.}};</script>
//
// then the template output is
//
// 	<script>var pair = {"A": "foo", "B": "bar"};</script>
//
// See package json to understand how non-string content is marshalled for
// embedding in JavaScript contexts.
//
//
// Typed Strings
//
// By default, this package assumes that all pipelines produce a plain text
// string. It adds escaping pipeline stages necessary to correctly and safely
// embed that plain text string in the appropriate context.
//
// When a data value is not plain text, you can make sure it is not over-escaped
// by marking it with its type.
//
// Types HTML, JS, URL, and others from content.go can carry safe content that
// is exempted from escaping.
//
// The template
//
// 	Hello, {{.}}!
//
// can be invoked with
//
// 	tmpl.Execute(out, template.HTML(`<b>World</b>`))
//
// to produce
//
// 	Hello, <b>World</b>!
//
// instead of the
//
// 	Hello, &lt;b&gt;World&lt;b&gt;!
//
// that would have been produced if {{.}} was a regular string.
//
//
// Security Model
//
// http://js-quasis-libraries-and-repl.googlecode.com/svn/trunk/safetemplate.html#problem_definition
// defines "safe" as used by this package.
//
// This package assumes that template authors are trusted, that Execute's data
// parameter is not, and seeks to preserve the properties below in the face of
// untrusted data:
//
// Structure Preservation Property: "... when a template author writes an HTML
// tag in a safe templating language, the browser will interpret the
// corresponding portion of the output as a tag regardless of the values of
// untrusted data, and similarly for other structures such as attribute
// boundaries and JS and CSS string boundaries."
//
// Code Effect Property: "... only code specified by the template author should
// run as a result of injecting the template output into a page and all code
// specified by the template author should run as a result of the same."
//
// Least Surprise Property: "A developer (or code reviewer) familiar with HTML,
// CSS, and JavaScript, who knows that contextual autoescaping happens should be
// able to look at a {{.}} and correctly infer what sanitization happens."
package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"text/template"
	"text/template/parse"
	"unicode"
	"unicode/utf8"
)

// We define codes for each error that manifests while escaping templates, but
// escaped templates may also fail at runtime.
//
// Output: "ZgotmplZ" Example:
//
// 	<img src="{{.X}}">
// 	where {{.X}} evaluates to `javascript:...`
//
// Discussion:
//
// 	"ZgotmplZ" is a special value that indicates that unsafe content reached a
// 	CSS or URL context at runtime. The output of the example will be
// 	  <img src="#ZgotmplZ">
// 	If the data comes from a trusted source, use content types to exempt it
// 	from filtering: URL(`javascript:...`).

// 我们为转义模板时的所有错误都定义了错误码，但经过转义修正的模板仍可能在运行时
// 出错：
//
// 输出"ZgotmplZ"的例子：
//
// 	<img src="{{.X}}">
// 	其中{{.X}}执行结果为`javascript:...`
//
// 讨论：
//
// 	"ZgotmplZ"是一个特殊值，表示运行时在CSS或URL上下文环境生成的不安全内容。本例的输出为：
// 	  <img src="#ZgotmplZ">
// 	如果数据来源可信，请转换内容类型来避免被滤除：URL(`javascript:...`)
const (
	// OK indicates the lack of an error.

	// OK表示没有出错
	OK ErrorCode = iota

	// ErrAmbigContext: "... appears in an ambiguous URL context"
	// Example:
	//   <a href="
	//      {{if .C}}
	//        /path/
	//      {{else}}
	//        /search?q=
	//      {{end}}
	//      {{.X}}
	//   ">
	// Discussion:
	//   {{.X}} is in an ambiguous URL context since, depending on {{.C}},
	//  it may be either a URL suffix or a query parameter.
	//   Moving {{.X}} into the condition removes the ambiguity:
	//   <a href="{{if .C}}/path/{{.X}}{{else}}/search?q={{.X}}">

	// 当上下文环境有歧义时导致ErrAmbigContext：
	// 举例：
	//   <a href="{{if .C}}/path/{{else}}/search?q={{end}}{{.X}}"&rt;
	// 说明：
	//   {{.X}}的URL上下文环境有歧义，因为根据{{.C}}的值，
	//   它可以是URL的后缀，或者是查询的参数。
	//   将{{.X}}移动到如下情况可以消除歧义：
	//   <a href="{{if .C}}/path/{{.X}}{{else}}/search?q={{.X}}{{end}}"&rt;
	ErrAmbigContext

	// ErrBadHTML: "expected space, attr name, or end of tag, but got ...",
	//   "... in unquoted attr", "... in attribute name"
	// Example:
	//   <a href = /search?q=foo>
	//   <href=foo>
	//   <form na<e=...>
	//   <option selected<
	// Discussion:
	//   This is often due to a typo in an HTML element, but some runes
	//   are banned in tag names, attribute names, and unquoted attribute
	//   values because they can tickle parser ambiguities.
	//   Quoting all attributes is the best policy.

	// 期望空白、属性名、标签结束标志而没有时，标签名或无引号标签值包含非法字符
	// 时， 会导致ErrBadHTML；举例：
	//
	// 	<a href = /search?q=foo&rt;
	// 	<href=foo&rt;
	// 	<form na<e=...&rt;
	// 	<option selected<
	//
	// 讨论：
	//
	// 	一般是因为HTML元素输入了错误的标签名、属性名或者未用引号的属性值，导致解析失败
	// 	将所有的属性都用引号括起来是最好的策略
	ErrBadHTML

	// ErrBranchEnd: "{{if}} branches end in different contexts"
	// Example:
	//   {{if .C}}<a href="{{end}}{{.X}}
	// Discussion:
	//   Package html/template statically examines each path through an
	//   {{if}}, {{range}}, or {{with}} to escape any following pipelines.
	//   The example is ambiguous since {{.X}} might be an HTML text node,
	//   or a URL prefix in an HTML attribute. The context of {{.X}} is
	//   used to figure out how to escape it, but that context depends on
	//   the run-time value of {{.C}} which is not statically known.
	//
	//   The problem is usually something like missing quotes or angle
	//   brackets, or can be avoided by refactoring to put the two contexts
	//   into different branches of an if, range or with. If the problem
	//   is in a {{range}} over a collection that should never be empty,
	//   adding a dummy {{else}} can help.

	// {{if}}等分支不在相同上下文开始和结束时，导致ErrBranchEnd 示例：
	//
	// 	{{if .C}}<a href="{{end}}{{.X}}
	//
	// 讨论：
	//
	// 	html/template包会静态的检验{{if}}、{{range}}或{{with}}的每一个分支，
	// 	以对后续的pipeline进行转义。该例出现了歧义，{{.X}}可能是HTML文本节点，
	// 	或者是HTML属性值的URL的前缀，{{.X}}的上下文环境可以确定如何转义，但该
	// 	上下文环境却是由运行时{{.C}}的值决定的，不能在编译期获知。
	// 	这种问题一般是因为缺少引号或者角括号引起的，另一些则可以通过重构将两个上下文
	// 	放进if、range、with的不同分支里来避免，如果问题出现在参数长度一定非0的
	// 	{{range}}的分支里，可以通过添加无效{{else}}分支解决。
	ErrBranchEnd

	// ErrEndContext: "... ends in a non-text context: ..."
	// Examples:
	//   <div
	//   <div title="no close quote>
	//   <script>f()
	// Discussion:
	//   Executed templates should produce a DocumentFragment of HTML.
	//   Templates that end without closing tags will trigger this error.
	//   Templates that should not be used in an HTML context or that
	//   produce incomplete Fragments should not be executed directly.
	//
	//   {{define "main"}} <script>{{template "helper"}}</script> {{end}}
	//   {{define "helper"}} document.write(' <div title=" ') {{end}}
	//
	//   "helper" does not produce a valid document fragment, so should
	//   not be Executed directly.

	// 如果以非文本上下文结束，则导致ErrEndContext 示例：
	//
	// 	<div
	// 	<div title="no close quote&rt;
	// 	<script>f()
	//
	// 讨论：
	//
	// 	执行模板必须生成HTML的一个文档片段，以未闭合标签结束的模板都会引发本错误。
	// 	不用在HTML上下文或者生成不完整片段的模板不应直接执行。
	// 	{{define "main"}} <script&rt;{{template "helper"}}</script> {{end}}
	// 	{{define "helper"}} document.write(' <div title=" ') {{end}}
	// 	模板"helper"不能生成合法的文档片段，所以不直接执行，用js生成。
	ErrEndContext

	// ErrNoSuchTemplate: "no such template ..."
	// Examples:
	//   {{define "main"}}<div {{template "attrs"}}>{{end}}
	//   {{define "attrs"}}href="{{.URL}}"{{end}}
	// Discussion:
	//   Package html/template looks through template calls to compute the
	//   context.
	//   Here the {{.URL}} in "attrs" must be treated as a URL when called
	//   from "main", but you will get this error if "attrs" is not defined
	//   when "main" is parsed.

	// 调用不存在的模板时导致ErrNoSuchTemplate 示例：
	//
	// 	{{define "main"}}<div {{template "attrs"}}&rt;{{end}}
	// 	{{define "attrs"}}href="{{.URL}}"{{end}}
	//
	// 讨论：
	//
	// 	html/template包略过模板调用计算上下文环境。
	// 	此例中，当被"main"模板调用时，"attrs"模板的{{.URL}}必须视为一个URL；
	// 	但如果解析"main"时，"attrs"还未被定义，就会导致本错误
	ErrNoSuchTemplate

	// ErrOutputContext: "cannot compute output context for template ..."
	// Examples:
	//   {{define "t"}}{{if .T}}{{template "t" .T}}{{end}}{{.H}}",{{end}}
	// Discussion:
	//   A recursive template does not end in the same context in which it
	//   starts, and a reliable output context cannot be computed.
	//   Look for typos in the named template.
	//   If the template should not be called in the named start context,
	//   look for calls to that template in unexpected contexts.
	//   Maybe refactor recursive templates to not be recursive.

	// 不能计算输出位置的上下文环境时，导致ErrOutputContext 示例：
	//
	// 	{{define "t"}}{{if .T}}{{template "t" .T}}{{end}}{{.H}}",{{end}}
	//
	// 讨论：
	//
	// 	一个递归的模板，其起始和结束的上下文环境不同时；
	// 	不能计算出可信的输出位置上下文环境时，就可能导致本错误。
	// 	检查各个命名模板是否有错误；
	// 	如果模板不应在命名的起始上下文环境调用，检查在不期望上下文环境中对该模板的调用；
	// 	或者将递归模板重构为非递归模板；
	ErrOutputContext

	// ErrPartialCharset: "unfinished JS regexp charset in ..."
	// Example:
	//     <script>var pattern = /foo[{{.Chars}}]/</script>
	// Discussion:
	//   Package html/template does not support interpolation into regular
	//   expression literal character sets.

	// 尚未支持JS正则表达式插入字符集
	// 示例：
	//     <script>var pattern = /foo[{{.Chars}}]/</script&rt;
	// 讨论：
	//   html/template不支持向JS正则表达式里插入字面值字符集
	ErrPartialCharset

	// ErrPartialEscape: "unfinished escape sequence in ..."
	// Example:
	//   <script>alert("\{{.X}}")</script>
	// Discussion:
	//   Package html/template does not support actions following a
	//   backslash.
	//   This is usually an error and there are better solutions; for
	//   example
	//     <script>alert("{{.X}}")</script>
	//   should work, and if {{.X}} is a partial escape sequence such as
	//   "xA0", mark the whole sequence as safe content: JSStr(`\xA0`)

	// 部分转义序列尚未支持
	// 示例：
	//   <script>alert("\{{.X}}")</script&rt;
	// 讨论：
	//   html/template包不支持紧跟在反斜杠后面的action
	//   这一般是错误的，有更好的解决方法，例如：
	//     <script>alert("{{.X}}")</script&rt;
	//   可以工作，如果{{.X}}是部分转义序列，如"xA0"，
	//   可以将整个序列标记为安全文本：JSStr(`\xA0`)
	ErrPartialEscape

	// ErrRangeLoopReentry: "on range loop re-entry: ..."
	// Example:
	//   <script>var x = [{{range .}}'{{.}},{{end}}]</script>
	// Discussion:
	//   If an iteration through a range would cause it to end in a
	//   different context than an earlier pass, there is no single context.
	//   In the example, there is missing a quote, so it is not clear
	//   whether {{.}} is meant to be inside a JS string or in a JS value
	//   context. The second iteration would produce something like
	//
	//     <script>var x = ['firstValue,'secondValue]</script>

	// range循环的重入口出错，导致ErrRangeLoopReentry 示例：
	//
	// 	<script>var x = [{{range .}}'{{.}},{{end}}]</script&rt;
	//
	// 讨论：
	//
	// 	如果range的迭代部分导致其结束于上一次循环的另一上下文，将不会有唯一的上下文环境
	// 	此例中，缺少一个引号，因此无法确定{{.}}是存在于一个JS字符串里，还是一个JS值文本里。
	// 	第二次迭代生成类似下面的输出：
	// 	  <script>var x = ['firstValue,'secondValue]</script&rt;
	ErrRangeLoopReentry

	// ErrSlashAmbig: '/' could start a division or regexp.
	// Example:
	//   <script>
	//     {{if .C}}var x = 1{{end}}
	//     /-{{.N}}/i.test(x) ? doThis : doThat();
	//   </script>
	// Discussion:
	//   The example above could produce `var x = 1/-2/i.test(s)...`
	//   in which the first '/' is a mathematical division operator or it
	//   could produce `/-2/i.test(s)` in which the first '/' starts a
	//   regexp literal.
	//   Look for missing semicolons inside branches, and maybe add
	//   parentheses to make it clear which interpretation you intend.

	// 斜杠可以开始一个除法或者正则表达式 示例：
	//
	// 	<script&rt;
	// 	  {{if .C}}var x = 1{{end}}
	// 	  /-{{.N}}/i.test(x) ? doThis : doThat();
	// 	</script&rt;
	//
	// 讨论：
	//
	// 	上例可以生成`var x = 1/-2/i.test(s)...`，其中第一个斜杠作为除号；
	// 	或者它也可以生成`/-2/i.test(s)`，其中第一个斜杠生成一个正则表达式字面值
	// 	检查分支中是否缺少分号，或者使用括号来明确你的意图
	ErrSlashAmbig
)

// Strings of content from a trusted source.
type (
	// CSS encapsulates known safe content that matches any of:
	//
	// 	1. The CSS3 stylesheet production, such as `p { color: purple }`.
	// 	2. The CSS3 rule production, such as `a[href=~"https:"].foo#bar`.
	// 	3. CSS3 declaration productions, such as `color: red; margin: 2px`.
	// 	4. The CSS3 value production, such as `rgba(0, 0, 255, 127)`.
	//
	// See http://www.w3.org/TR/css3-syntax/#parsing and
	// https://web.archive.org/web/20090211114933/http://w3.org/TR/css3-syntax#style
	//
	// Use of this type presents a security risk: the encapsulated content
	// should come from a trusted source, as it will be included verbatim in the
	// template output.

	// CSS用于包装匹配如下任一条的已知安全的内容：
	//
	// 	1. CSS3样式表，如`p { color: purple }`
	// 	2. CSS3规则，如`a[href=~"https:"].foo#bar`
	// 	3. CSS3声明，如`color: red; margin: 2px`
	// 	4. CSS3规则，如`rgba(0, 0, 255, 127)`
	//
	// 参见：http://www.w3.org/TR/css3-syntax/#parsing
	//
	// 以及：
	// https://web.archive.org/web/20090211114933/http://w3.org/TR/css3-syntax#style
	CSS string

	// HTML encapsulates a known safe HTML document fragment.
	// It should not be used for HTML from a third-party, or HTML with
	// unclosed tags or comments. The outputs of a sound HTML sanitizer
	// and a template escaped by this package are fine for use with HTML.
	//
	// Use of this type presents a security risk:
	// the encapsulated content should come from a trusted source,
	// as it will be included verbatim in the template output.
	HTML string

	// HTMLAttr encapsulates an HTML attribute from a trusted source,
	// for example, ` dir="ltr"`.
	//
	// Use of this type presents a security risk:
	// the encapsulated content should come from a trusted source,
	// as it will be included verbatim in the template output.
	HTMLAttr string

	// JS encapsulates a known safe EcmaScript5 Expression, for example,
	// `(x + y * z())`.
	// Template authors are responsible for ensuring that typed expressions
	// do not break the intended precedence and that there is no
	// statement/expression ambiguity as when passing an expression like
	// "{ foo: bar() }\n['foo']()", which is both a valid Expression and a
	// valid Program with a very different meaning.
	//
	// Use of this type presents a security risk:
	// the encapsulated content should come from a trusted source,
	// as it will be included verbatim in the template output.
	//
	// Using JS to include valid but untrusted JSON is not safe.
	// A safe alternative is to parse the JSON with json.Unmarshal and then
	// pass the resultant object into the template, where it will be
	// converted to sanitized JSON when presented in a JavaScript context.
	JS string

	// JSStr encapsulates a sequence of characters meant to be embedded
	// between quotes in a JavaScript expression.
	// The string must match a series of StringCharacters:
	//   StringCharacter :: SourceCharacter but not `\` or LineTerminator
	//                    | EscapeSequence
	// Note that LineContinuations are not allowed.
	// JSStr("foo\\nbar") is fine, but JSStr("foo\\\nbar") is not.
	//
	// Use of this type presents a security risk:
	// the encapsulated content should come from a trusted source,
	// as it will be included verbatim in the template output.
	JSStr string

	// URL encapsulates a known safe URL or URL substring (see RFC 3986).
	// A URL like `javascript:checkThatFormNotEditedBeforeLeavingPage()`
	// from a trusted source should go in the page, but by default dynamic
	// `javascript:` URLs are filtered out since they are a frequently
	// exploited injection vector.
	//
	// Use of this type presents a security risk:
	// the encapsulated content should come from a trusted source,
	// as it will be included verbatim in the template output.
	URL string
)

// Error describes a problem encountered during template Escaping.

// Error描述在模板转义时出现的错误。
type Error struct {
	// ErrorCode describes the kind of error.
	ErrorCode ErrorCode

	// Node is the node that caused the problem, if known.
	// If not nil, it overrides Name and Line.
	Node parse.Node

	// Name is the name of the template in which the error was encountered.
	Name string

	// Line is the line number of the error in the template source or 0.
	Line int

	// Description is a human-readable description of the problem.
	Description string
}

// ErrorCode is a code for a kind of error.

// ErrorCode是代表错误种类的错误码。
type ErrorCode int

// FuncMap is the type of the map defining the mapping from names to
// functions. Each function must have either a single return value, or two
// return values of which the second has type error. In that case, if the
// second (error) argument evaluates to non-nil during execution, execution
// terminates and Execute returns that error. FuncMap has the same base type
// as FuncMap in "text/template", copied here so clients need not import
// "text/template".

// FuncMap类型定义了函数名字符串到函数的映射，每个函数都必须有1到2个返回值，如果
// 有2个则后一个必须是error接口类型；如果有2个返回值的方法返回的error非nil，模板
// 执行会中断并返回给调用者该错误。该类型拷贝自text/template包的同名类型，因此不
// 需要导入该包以使用该类型。
type FuncMap map[string]interface{}

// Template is a specialized Template from "text/template" that produces a safe
// HTML document fragment.

// Template类型是text/template包的Template类型的特化版本，用于生成安全的HTML文本
// 片段。
type Template struct {
	// The underlying template's parse tree, updated to be HTML-safe.
	Tree *parse.Tree
}

// HTMLEscape writes to w the escaped HTML equivalent of the plain text data b.

// 函数向w中写入b的HTML转义等价表示。
func HTMLEscape(w io.Writer, b []byte)

// HTMLEscapeString returns the escaped HTML equivalent of the plain text data
// s.

// 返回s的HTML转义等价表示字符串。
func HTMLEscapeString(s string) string

// HTMLEscaper returns the escaped HTML equivalent of the textual
// representation of its arguments.

// 函数返回其所有参数文本表示的HTML转义等价表示字符串。
func HTMLEscaper(args ...interface{}) string

// IsTrue reports whether the value is 'true', in the sense of not the zero of
// its type, and whether the value has a meaningful truth value. This is the
// definition of truth used by if and other such actions.
func IsTrue(val interface{}) (truth, ok bool)

// JSEscape writes to w the escaped JavaScript equivalent of the plain text data
// b.

// 函数向w中写入b的JavaScript转义等价表示。
func JSEscape(w io.Writer, b []byte)

// JSEscapeString returns the escaped JavaScript equivalent of the plain text
// data s.

// 返回s的JavaScript转义等价表示字符串。
func JSEscapeString(s string) string

// JSEscaper returns the escaped JavaScript equivalent of the textual
// representation of its arguments.

// 函数返回其所有参数文本表示的JavaScript转义等价表示字符串。
func JSEscaper(args ...interface{}) string

// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//
// 	var t = template.Must(template.New("name").Parse("html"))

// Must函数用于包装返回(*Template,
// error)的函数/方法调用，它会在err非nil时panic，一般用于变量初始化：
//
//     var t = template.Must(template.New("name").Parse("html"))
func Must(t *Template, err error) *Template

// New allocates a new HTML template with the given name.

// 创建一个名为name的模板。
func New(name string) *Template

// ParseFiles creates a new Template and parses the template definitions from
// the named files. The returned template's name will have the (base) name and
// (parsed) contents of the first file. There must be at least one file. If an
// error occurs, parsing stops and the returned *Template is nil.
//
// When parsing multiple files with the same name in different directories, the
// last one mentioned will be the one that results. For instance,
// ParseFiles("a/foo", "b/foo") stores "b/foo" as the template named "foo",
// while "a/foo" is unavailable.

// ParseFiles函数创建一个模板并解析filenames指定的文件里的模板定义。返回的模板的
// 名字是第一个文件的文件名（不含扩展名），内容为解析后的第一个文件的内容。至少
// 要提供一个文件。如果发生错误，会停止解析并返回nil。
func ParseFiles(filenames ...string) (*Template, error)

// ParseGlob creates a new Template and parses the template definitions from the
// files identified by the pattern, which must match at least one file. The
// returned template will have the (base) name and (parsed) contents of the
// first file matched by the pattern. ParseGlob is equivalent to calling
// ParseFiles with the list of files matched by the pattern.
//
// When parsing multiple files with the same name in different directories,
// the last one mentioned will be the one that results.

// ParseGlob创建一个模板并解析匹配pattern的文件（参见glob规则）里的模板定义。返
// 回的模板的名字是第一个匹配的文件的文件名（不含扩展名），内容为解析后的第一个
// 文件的内容。至少要存在一个匹配的文件。如果发生错误，会停止解析并返回nil。
// ParseGlob等价于使用匹配pattern的文件的列表为参数调用ParseFiles。
func ParseGlob(pattern string) (*Template, error)

// URLQueryEscaper returns the escaped value of the textual representation of
// its arguments in a form suitable for embedding in a URL query.

// 函数返回其所有参数文本表示的可以嵌入URL查询的转义等价表示字符串。
func URLQueryEscaper(args ...interface{}) string

func (e *Error) Error() string

// AddParseTree creates a new template with the name and parse tree
// and associates it with t.
//
// It returns an error if t has already been executed.

// AddParseTree方法使用name和tree创建一个模板并使它和t相关联。
//
// 如果t已经执行过了，会返回错误。
func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error)

// Clone returns a duplicate of the template, including all associated
// templates. The actual representation is not copied, but the name space of
// associated templates is, so further calls to Parse in the copy will add
// templates to the copy but not to the original. Clone can be used to prepare
// common templates and use them with variant definitions for other templates by
// adding the variants after the clone is made.
//
// It returns an error if t has already been executed.

// Clone方法返回模板的一个副本，包括所有相关联的模板。模板的底层表示树并未拷贝，
// 而是拷贝了命名空间，因此拷贝调用Parse方法不会修改原模板的命名空间。Clone方法
// 用于准备模板的公用部分，向拷贝中加入其他关联模板后再进行使用。
//
// 如果t已经执行过了，会返回错误。
func (t *Template) Clone() (*Template, error)

// DefinedTemplates returns a string listing the defined templates,
// prefixed by the string "; defined templates are: ". If there are none,
// it returns the empty string. Used to generate an error message.
func (t *Template) DefinedTemplates() string

// Delims sets the action delimiters to the specified strings, to be used in
// subsequent calls to Parse, ParseFiles, or ParseGlob. Nested template
// definitions will inherit the settings. An empty delimiter stands for the
// corresponding default: {{ or }}.
// The return value is the template, so calls can be chained.

// Delims方法用于设置action的分界字符串，应用于之后的Parse、ParseFiles、
// ParseGlob方法。嵌套模板定义会继承这种分界符设置。空字符串分界符表示相应的默认
// 分界符：{{或}}。返回值就是t，以便进行链式调用。
func (t *Template) Delims(left, right string) *Template

// Execute applies a parsed template to the specified data object,
// writing the output to wr.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
// A template may be executed safely in parallel.

// Execute方法将解析好的模板应用到data上，并将输出写入wr。如果执行时出现错误，会
// 停止执行，但有可能已经写入wr部分数据。模板可以安全的并发执行。
func (t *Template) Execute(wr io.Writer, data interface{}) error

// ExecuteTemplate applies the template associated with t that has the given
// name to the specified data object and writes the output to wr.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
// A template may be executed safely in parallel.

// ExecuteTemplate方法类似Execute，但是使用名为name的t关联的模板产生输出。
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error

// Funcs adds the elements of the argument map to the template's function map.
// It panics if a value in the map is not a function with appropriate return
// type. However, it is legal to overwrite elements of the map. The return
// value is the template, so calls can be chained.

// Funcs方法向模板t的函数字典里加入参数funcMap内的键值对。如果funcMap某个键值对
// 的值不是函数类型或者返回值不符合要求会panic。但是，可以对t函数列表的成员进行
// 重写。方法返回t以便进行链式调用。
func (t *Template) Funcs(funcMap FuncMap) *Template

// Lookup returns the template with the given name that is associated with t,
// or nil if there is no such template.

// Lookup方法返回与t关联的名为name的模板，如果没有这个模板会返回nil。
func (t *Template) Lookup(name string) *Template

// Name returns the name of the template.

// 返回模板t的名字。
func (t *Template) Name() string

// New allocates a new HTML template associated with the given one
// and with the same delimiters. The association, which is transitive,
// allows one template to invoke another with a {{template}} action.

// New方法创建一个和t关联的名字为name的模板并返回它。这种可以传递的关联允许一个
// 模板使用template action调用另一个模板。
func (t *Template) New(name string) *Template

// Option sets options for the template. Options are described by
// strings, either a simple string or "key=value". There can be at
// most one equals sign in an option string. If the option string
// is unrecognized or otherwise invalid, Option panics.
//
// Known options:
//
// missingkey: Control the behavior during execution if a map is
// indexed with a key that is not present in the map.
// 	"missingkey=default" or "missingkey=invalid"
// 		The default behavior: Do nothing and continue execution.
// 		If printed, the result of the index operation is the string
// 		"<no value>".
// 	"missingkey=zero"
// 		The operation returns the zero value for the map type's element.
// 	"missingkey=error"
// 		Execution stops immediately with an error.
func (t *Template) Option(opt ...string) *Template

// Parse parses a string into a template. Nested template definitions
// will be associated with the top-level template t. Parse may be
// called multiple times to parse definitions of templates to associate
// with t. It is an error if a resulting template is non-empty (contains
// content other than template definitions) and would replace a
// non-empty template with the same name.  (In multiple calls to Parse
// with the same receiver template, only one call can contain text
// other than space, comments, and template definitions.)

// Parse方法将字符串text解析为模板。嵌套定义的模板会关联到最顶层的t。Parse可以多
// 次调用，但只有第一次调用可以包含空格、注释和模板定义之外的文本。如果后面的调
// 用在解析后仍剩余文本会引发错误、返回nil且丢弃剩余文本；如果解析得到的模板已有
// 相关联的同名模板，会覆盖掉原模板。
func (t *Template) Parse(src string) (*Template, error)

// ParseFiles parses the named files and associates the resulting templates with
// t. If an error occurs, parsing stops and the returned template is nil;
// otherwise it is t. There must be at least one file.
//
// When parsing multiple files with the same name in different directories,
// the last one mentioned will be the one that results.

// ParseGlob方法解析filenames指定的文件里的模板定义并将解析结果与t关联。如果发生
// 错误，会停止解析并返回nil，否则返回(t, nil)。至少要提供一个文件。
func (t *Template) ParseFiles(filenames ...string) (*Template, error)

// ParseGlob parses the template definitions in the files identified by the
// pattern and associates the resulting templates with t. The pattern is
// processed by filepath.Glob and must match at least one file. ParseGlob is
// equivalent to calling t.ParseFiles with the list of files matched by the
// pattern.
//
// When parsing multiple files with the same name in different directories,
// the last one mentioned will be the one that results.

// ParseFiles方法解析匹配pattern的文件里的模板定义并将解析结果与t关联。如果发生
// 错误，会停止解析并返回nil，否则返回(t, nil)。至少要存在一个匹配的文件。
func (t *Template) ParseGlob(pattern string) (*Template, error)

// Templates returns a slice of the templates associated with t, including t
// itself.

// Templates方法返回与t相关联的模板的切片，包括t自己。
func (t *Template) Templates() []*Template

