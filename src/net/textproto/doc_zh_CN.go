// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package textproto implements generic support for text-based request/response
// protocols in the style of HTTP, NNTP, and SMTP.
//
// The package provides:
//
// Error, which represents a numeric error response from a server.
//
// Pipeline, to manage pipelined requests and responses in a client.
//
// Reader, to read numeric response code lines, key: value headers, lines wrapped
// with leading spaces on continuation lines, and whole text blocks ending with a
// dot on a line by itself.
//
// Writer, to write dot-encoded text blocks.
//
// Conn, a convenient packaging of Reader, Writer, and Pipeline for use with a
// single network connection.

// textproto实现了对基于文本的请求/回复协议的一般性支持，包括HTTP、NNTP和SMTP。
//
// 本包提供：
//
// 错误，代表服务端回复的错误码。Pipeline，以管理客户端中的管道化的请求/回复。Reader，读取数值回复码行，键值对形式的头域，一个作为后续行先导的空行，以及以只有一个"."的一行为结尾的整个文本块。Writer，写入点编码的文本。Conn，对Reader、Writer和Pipline的易用的包装，用于单个网络连接。
package textproto

// CanonicalMIMEHeaderKey returns the canonical format of the MIME header key s.
// The canonicalization converts the first letter and any letter following a hyphen
// to upper case; the rest are converted to lowercase. For example, the canonical
// key for "accept-encoding" is "Accept-Encoding". MIME header keys are assumed to
// be ASCII only.

// 返回一个MIME头的键的规范格式。该标准会将首字母和所有"-"之后的字符改为大写，其余字母改为小写。举个例子，"accept-encoding"作为键的标准格式是"Accept-Encoding"。MIME头的键必须是ASCII码构成。
func CanonicalMIMEHeaderKey(s string) string

// TrimBytes returns b without leading and trailing ASCII space.

// 去掉b前后的ASCII码空白（不去Unicode空白）
func TrimBytes(b []byte) []byte

// TrimString returns s without leading and trailing ASCII space.

// 去掉s前后的ASCII码空白（不去Unicode空白）
func TrimString(s string) string

// A Conn represents a textual network protocol connection. It consists of a Reader
// and Writer to manage I/O and a Pipeline to sequence concurrent requests on the
// connection. These embedded types carry methods with them; see the documentation
// of those types for details.

// Conn代表一个文本网络协议的连接。它包含一个Reader和一个Writer来管理读写，一个Pipeline来对连接中并行的请求进行排序。匿名嵌入的类型字段是Conn可以调用它们的方法。
type Conn struct {
	Reader
	Writer
	Pipeline
	// contains filtered or unexported fields
}

// Dial connects to the given address on the given network using net.Dial and then
// returns a new Conn for the connection.

// Dial函数使用net.Dial在给定网络上和给定地址建立网络连接，并返回用于该连接的Conn。
func Dial(network, addr string) (*Conn, error)

// NewConn returns a new Conn using conn for I/O.

// NewConn函数返回以I/O为底层的Conn。
func NewConn(conn io.ReadWriteCloser) *Conn

// Close closes the connection.

// Close方法关闭连接。
func (c *Conn) Close() error

// Cmd is a convenience method that sends a command after waiting its turn in the
// pipeline. The command text is the result of formatting format with args and
// appending \r\n. Cmd returns the id of the command, for use with StartResponse
// and EndResponse.
//
// For example, a client might run a HELP command that returns a dot-body by using:
//
//	id, err := c.Cmd("HELP")
//	if err != nil {
//		return nil, err
//	}
//
//	c.StartResponse(id)
//	defer c.EndResponse(id)
//
//	if _, _, err = c.ReadCodeLine(110); err != nil {
//		return nil, err
//	}
//	text, err := c.ReadDotBytes()
//	if err != nil {
//		return nil, err
//	}
//	return c.ReadCodeLine(250)

// Cmd方法用于在管道中等待轮到它执行，并发送命令。命令文本是用给定的format字符串和参数格式化生成的。并会在最后添加上\r\n。Cmd函数返回该命令的Pipeline
// id，用于StartResponse和EndResponse方法。
//
// 例如，一个客户端可以使用如下代码执行HELP命令并返回解码后的点编码文本：
//
//	id, err := c.Cmd("HELP")
//	if err != nil {
//		return nil, err
//	}
//	c.StartResponse(id)
//	defer c.EndResponse(id)
//	if _, _, err = c.ReadCodeLine(110); err != nil {
//		return nil, err
//	}
//	text, err := c.ReadDotBytes()
//	if err != nil {
//		return nil, err
//	}
//	return c.ReadCodeLine(250)
func (c *Conn) Cmd(format string, args ...interface{}) (id uint, err error)

// An Error represents a numeric error response from a server.

// Error代表一个服务端返回的数值状态码/错误码。
type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string

// A MIMEHeader represents a MIME-style header mapping keys to sets of values.

// MIMEHeader代表一个MIME头，将键映射为值的集合。
type MIMEHeader map[string][]string

// Add adds the key, value pair to the header. It appends to any existing values
// associated with key.

// Add方法向h中添加键值对，它会把新的值添加到键对应的值的集合里。
func (h MIMEHeader) Add(key, value string)

// Del deletes the values associated with key.

// Del方法删除键对应的值集。
func (h MIMEHeader) Del(key string)

// Get gets the first value associated with the given key. If there are no values
// associated with the key, Get returns "". Get is a convenience method. For more
// complex queries, access the map directly.

// Get方法返回键对应的值集的第一个值。如果键没有关联值，返回""。如要获得键对应的值集直接用map。
func (h MIMEHeader) Get(key string) string

// Set sets the header entries associated with key to the single element value. It
// replaces any existing values associated with key.

// Set方法将键对应的值集设置为只含有value一个值。没有就新建，有则删掉原有的值。
func (h MIMEHeader) Set(key, value string)

// A Pipeline manages a pipelined in-order request/response sequence.
//
// To use a Pipeline p to manage multiple clients on a connection, each client
// should run:
//
//	id := p.Next()	// take a number
//
//	p.StartRequest(id)	// wait for turn to send request
//	«send request»
//	p.EndRequest(id)	// notify Pipeline that request is sent
//
//	p.StartResponse(id)	// wait for turn to read response
//	«read response»
//	p.EndResponse(id)	// notify Pipeline that response is read
//
// A pipelined server can use the same calls to ensure that responses computed in
// parallel are written in the correct order.

// Pipeline管理管道化的有序请求/回复序列。
//
// 为了使用Pipeline管理一个连接的多个客户端，每个客户端应像下面一样运行：
//
//	id := p.Next()      // 获取一个数字id
//	p.StartRequest(id)  // 等待轮到该id发送请求
//	«send request»
//	p.EndRequest(id)    // 通知Pipeline请求发送完毕
//	p.StartResponse(id) // 等待该id读取回复
//	«read response»
//	p.EndResponse(id)   // 通知Pipeline回复已经读取
//
// 一个管道化的服务器可以使用相同的调用来保证回复并行的生成并以正确的顺序写入。
type Pipeline struct {
	// contains filtered or unexported fields
}

// EndRequest notifies p that the request with the given id has been sent (or, if
// this is a server, received).

// 通知p，给定id的request的操作已经结束了。
func (p *Pipeline) EndRequest(id uint)

// EndResponse notifies p that the response with the given id has been received
// (or, if this is a server, sent).

// 通知p，给定id的response的操作已经结束了。
func (p *Pipeline) EndResponse(id uint)

// Next returns the next id for a request/response pair.

// 返回下一对request/response的id。
func (p *Pipeline) Next() uint

// StartRequest blocks until it is time to send (or, if this is a server, receive)
// the request with the given id.

// 阻塞程序，直到轮到给定id来发送（读取）request。
func (p *Pipeline) StartRequest(id uint)

// StartResponse blocks until it is time to receive (or, if this is a server, send)
// the request with the given id.

// 阻塞程序，直到轮到给定id来读取（发送）response。
func (p *Pipeline) StartResponse(id uint)

// A ProtocolError describes a protocol violation such as an invalid response or a
// hung-up connection.

// ProtocolError描述一个违反协议的错误，如不合法的回复或者挂起的连接。
type ProtocolError string

func (p ProtocolError) Error() string

// A Reader implements convenience methods for reading requests or responses from a
// text protocol network connection.

// Reader实现了从一个文本协议网络连接中方便的读取请求/回复的方法。
type Reader struct {
	R *bufio.Reader
	// contains filtered or unexported fields
}

// NewReader returns a new Reader reading from r.

// NewReader返回一个从r读取数据的Reader。
func NewReader(r *bufio.Reader) *Reader

// DotReader returns a new Reader that satisfies Reads using the decoded text of a
// dot-encoded block read from r. The returned Reader is only valid until the next
// call to a method on r.
//
// Dot encoding is a common framing used for data blocks in text protocols such as
// SMTP. The data consists of a sequence of lines, each of which ends in "\r\n".
// The sequence itself ends at a line containing just a dot: ".\r\n". Lines
// beginning with a dot are escaped with an additional dot to avoid looking like
// the end of the sequence.
//
// The decoded form returned by the Reader's Read method rewrites the "\r\n" line
// endings into the simpler "\n", removes leading dot escapes if present, and stops
// with error io.EOF after consuming (and discarding) the end-of-sequence line.

// DotReader方法返回一个io.Reader，该接口自动解码r中读取的点编码块。注意该接口仅在下一次调用r的方法之前才有效。点编码是文本协议如SMTP用于文本块的通用框架。数据包含多个行，每行以"\r\n"结尾。数据本身以一个只含有一个点的一行".\r\n"来结尾。以点起始的行会添加额外的点，来避免看起来像是文本的结尾。
//
// 返回接口的Read方法会将行尾的"\r\n"修改为"\n"，去掉起头的转义点，并在底层读取到（并抛弃掉）表示文本结尾的行时停止解码并返回io.EOF错误。
func (r *Reader) DotReader() io.Reader

// ReadCodeLine reads a response code line of the form
//
//	code message
//
// where code is a three-digit status code and the message extends to the rest of
// the line. An example of such a line is:
//
//	220 plan9.bell-labs.com ESMTP
//
// If the prefix of the status does not match the digits in expectCode,
// ReadCodeLine returns with err set to &Error{code, message}. For example, if
// expectCode is 31, an error will be returned if the status is not in the range
// [310,319].
//
// If the response is multi-line, ReadCodeLine returns an error.
//
// An expectCode <= 0 disables the check of the status code.

// 方法读取回复的状态码行，格式如下：
//
//	code message
//
// 状态码是3位数字，message进一步描述状态，例如：
//
//	220 plan9.bell-labs.com ESMTP
//
// 如果状态码字符串的前缀不匹配expectCode，方法返回错误&Error{code,
// message}。例如expectCode是31，则如果状态码不在区间[310,
// 319]内就会返回错误。如果回复是多行的则会返回错误。
//
// 如果expectCode <= 0，将不会检查状态码。
func (r *Reader) ReadCodeLine(expectCode int) (code int, message string, err error)

// ReadContinuedLine reads a possibly continued line from r, eliding the final
// trailing ASCII white space. Lines after the first are considered continuations
// if they begin with a space or tab character. In the returned data, continuation
// lines are separated from the previous line only by a single space: the newline
// and leading white space are removed.
//
// For example, consider this input:
//
//	Line 1
//	  continued...
//	Line 2
//
// The first call to ReadContinuedLine will return "Line 1 continued..." and the
// second will return "Line 2".
//
// A line consisting of only white space is never continued.

// ReadContinuedLine从r中读取可能有后续的行，会将该行尾段的ASCII空白剔除，并将该行后面所有以空格或者tab起始的行视为其后续，后续部分会剔除行头部的空白，所有这些行包括第一行以单个空格连接起来返回。
//
// 举例如下：
//
//	Line 1
//	  continued...
//	Line 2
//
// 第一次调用ReadContinuedLine会返回"Line 1 continued..."，第二次会返回"Line 2"
//
// 只有空格的行不被视为有后续的行。
func (r *Reader) ReadContinuedLine() (string, error)

// ReadContinuedLineBytes is like ReadContinuedLine but returns a []byte instead of
// a string.

// ReadContinuedLineBytes类似ReadContinuedLine但返回[]byte切片。
func (r *Reader) ReadContinuedLineBytes() ([]byte, error)

// ReadDotBytes reads a dot-encoding and returns the decoded data.
//
// See the documentation for the DotReader method for details about dot-encoding.

// ReadDotBytes读取点编码文本返回解码后的数据，点编码详见DotReader方法。
func (r *Reader) ReadDotBytes() ([]byte, error)

// ReadDotLines reads a dot-encoding and returns a slice containing the decoded
// lines, with the final \r\n or \n elided from each.
//
// See the documentation for the DotReader method for details about dot-encoding.

// ReadDotLines方法读取一个点编码文本块并返回一个包含解码后各行的切片，各行最后的\r\n或\n去掉。
func (r *Reader) ReadDotLines() ([]string, error)

// ReadLine reads a single line from r, eliding the final \n or \r\n from the
// returned string.

// ReadLine方法从r读取单行，去掉最后的\r\n或\n。
func (r *Reader) ReadLine() (string, error)

// ReadLineBytes is like ReadLine but returns a []byte instead of a string.

// ReadLineBytes类似ReadLine但返回[]byte切片。
func (r *Reader) ReadLineBytes() ([]byte, error)

// ReadMIMEHeader reads a MIME-style header from r. The header is a sequence of
// possibly continued Key: Value lines ending in a blank line. The returned map m
// maps CanonicalMIMEHeaderKey(key) to a sequence of values in the same order
// encountered in the input.
//
// For example, consider this input:
//
//	My-Key: Value 1
//	Long-Key: Even
//	       Longer Value
//	My-Key: Value 2
//
// Given that input, ReadMIMEHeader returns the map:
//
//	map[string][]string{
//		"My-Key": {"Value 1", "Value 2"},
//		"Long-Key": {"Even Longer Value"},
//	}

// ReadMIMEHeader从r读取MIME风格的头域。该头域包含一系列可能有后续的键值行，以空行结束。返回的map映射CanonicalMIMEHeaderKey(key)到值的序列（顺序与输入相同）。
//
// 举例如下：
//
//	My-Key: Value 1
//	Long-Key: Even
//	       Longer Value
//	My-Key: Value 2
//
// 对此输入，ReadMIMEHeader返回：
//
//	map[string][]string{
//		"My-Key": {"Value 1", "Value 2"},
//		"Long-Key": {"Even Longer Value"},
//	}
func (r *Reader) ReadMIMEHeader() (MIMEHeader, error)

// ReadResponse reads a multi-line response of the form:
//
//	code-message line 1
//	code-message line 2
//	...
//	code message line n
//
// where code is a three-digit status code. The first line starts with the code and
// a hyphen. The response is terminated by a line that starts with the same code
// followed by a space. Each line in message is separated by a newline (\n).
//
// See page 36 of RFC 959 (http://www.ietf.org/rfc/rfc959.txt) for details.
//
// If the prefix of the status does not match the digits in expectCode,
// ReadResponse returns with err set to &Error{code, message}. For example, if
// expectCode is 31, an error will be returned if the status is not in the range
// [310,319].
//
// An expectCode <= 0 disables the check of the status code.

// ReadResponse方法读取如下格式的多行回复：
//
//	code-message line 1
//	code-message line 2
//	...
//	code message line n
//
// 其中code是三位数的状态码。第一行以code和连字符开始，最后以同code后跟空格的行结束。返回值message每行以\n分隔。细节参见RFC
// 959(http://www.ietf.org/rfc/rfc959.txt)第36页。
//
// 如果状态码字符串的前缀不匹配expectCode，方法返回时err设为&Error{code,
// message}。例如expectCode是31，则如果状态码不在区间[310,
// 319]内就会返回错误。如果回复是多行的则会返回错误。
//
// 如果expectCode <= 0，将不会检查状态码。
func (r *Reader) ReadResponse(expectCode int) (code int, message string, err error)

// A Writer implements convenience methods for writing requests or responses to a
// text protocol network connection.

// Writer实现了方便的方法在一个文本协议网络连接中写入请求/回复。
type Writer struct {
	W *bufio.Writer
	// contains filtered or unexported fields
}

// NewWriter returns a new Writer writing to w.

// NewWriter函数返回一个底层写入w的Writer。
func NewWriter(w *bufio.Writer) *Writer

// DotWriter returns a writer that can be used to write a dot-encoding to w. It
// takes care of inserting leading dots when necessary, translating line-ending \n
// into \r\n, and adding the final .\r\n line when the DotWriter is closed. The
// caller should close the DotWriter before the next call to a method on w.
//
// See the documentation for Reader's DotReader method for details about
// dot-encoding.

// DotWriter方法返回一个io.WriteCloser，用于将点编码文本写入w。返回的接口会在必要时添加转义点，将行尾的\n替换为\r\n，并在关闭时添加最后的.\r\n行。调用者必须在下一次调用w的方法前关闭该接口。点编码文本格式参见Reader.DotReader方法。
func (w *Writer) DotWriter() io.WriteCloser

// PrintfLine writes the formatted output followed by \r\n.

// PrintfLine方法将格式化的输出写入底层并在最后写入\r\n。
func (w *Writer) PrintfLine(format string, args ...interface{}) error
