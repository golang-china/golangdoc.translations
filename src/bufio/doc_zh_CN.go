// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer
// object, creating another object (Reader or Writer) that also implements the
// interface but provides buffering and some help for textual I/O.

// bufio 包实现了带缓存的I/O操作.
// 它封装了一个io.Reader或者io.Writer对象，另外创建了一个对象
// （Reader或者Writer），这个对象也实现了一个接口，并提供缓冲和文档读写的帮助。
package bufio

const (
	// MaxScanTokenSize is the maximum size used to buffer a token.
	// The actual maximum token size may be smaller as the buffer
	// may need to include, for instance, a newline.
	MaxScanTokenSize = 64 * 1024
)

var (
	ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
	ErrInvalidUnreadRune = errors.New("bufio: invalid use of UnreadRune")
	ErrBufferFull        = errors.New("bufio: buffer full")
	ErrNegativeCount     = errors.New("bufio: negative count")
)

// Errors returned by Scanner.
var (
	ErrTooLong         = errors.New("bufio.Scanner: token too long")
	ErrNegativeAdvance = errors.New("bufio.Scanner: SplitFunc returns negative advance count")
	ErrAdvanceTooFar   = errors.New("bufio.Scanner: SplitFunc returns advance count beyond input")
)

// ScanBytes is a split function for a Scanner that returns each byte as a token.

// ScanBytes是用于Scanner类型的分割函数（符合SplitFunc），
// 本函数会将每个字节作为一个token返回。
func ScanBytes(data []byte, atEOF bool) (advance int, token []byte, err error)

// ScanLines is a split function for a Scanner that returns each line of text,
// stripped of any trailing end-of-line marker. The returned line may be empty. The
// end-of-line marker is one optional carriage return followed by one mandatory
// newline. In regular expression notation, it is `\r?\n`. The last non-empty line
// of input will be returned even if it has no newline.

// ScanLines是用于Scanner类型的分割函数（符合SplitFunc），
// 本函数会将每一行文本去掉末尾的换行标记作为一个token返回。
// 返回的行可以是空字符串。换行标记为一个可选的回车后跟一个必选的换行符。
// 最后一行即使没有换行符也会作为一个token返回。
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error)

// ScanRunes is a split function for a Scanner that returns each UTF-8-encoded rune
// as a token. The sequence of runes returned is equivalent to that from a range
// loop over the input as a string, which means that erroneous UTF-8 encodings
// translate to U+FFFD = "\xef\xbf\xbd". Because of the Scan interface, this makes
// it impossible for the client to distinguish correctly encoded replacement runes
// from encoding errors.

// ScanRunes是用于Scanner类型的分割函数（符合SplitFunc），
// 本函数会将每个utf-8编码的unicode码值作为一个token返回。
// 本函数返回的rune序列和range一个字符串的输出rune序列相同。
// 错误的utf-8编码会翻译为U+FFFD = "\xef\xbf\xbd"，
// 但只会消耗一个字节。调用者无法区分正确编码的rune和错误编码的rune。
func ScanRunes(data []byte, atEOF bool) (advance int, token []byte, err error)

// ScanWords is a split function for a Scanner that returns each space-separated
// word of text, with surrounding spaces deleted. It will never return an empty
// string. The definition of space is set by unicode.IsSpace.

// ScanWords是用于Scanner类型的分割函数（符合SplitFunc），
// 本函数会将每一行文本去掉末尾的换行标记作为一个token返回。
// 返回的行可以是空字符串。换行标记为一个可选的回车后跟一个必选的换行符。
// 最后一行即使没有换行符也会作为一个token返回。
func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error)

// ReadWriter stores pointers to a Reader and a Writer. It implements
// io.ReadWriter.

// ReadWriter存储输入输出指针。 它实现了io.ReadWriter。
type ReadWriter struct {
	*Reader
	*Writer
}

// NewReadWriter allocates a new ReadWriter that dispatches to r and w.

// NewReadWriter分配新的ReadWriter来进行r和w的调度。
func NewReadWriter(r *Reader, w *Writer) *ReadWriter

// Reader implements buffering for an io.Reader object.

// Reader实现了对一个io.Reader对象的缓冲读。
type Reader struct {
	// contains filtered or unexported fields
}

// NewReader returns a new Reader whose buffer has the default size.

// NewReader返回一个新的Reader，这个Reader的大小是默认的大小。
func NewReader(rd io.Reader) *Reader

// NewReaderSize returns a new Reader whose buffer has at least the specified size.
// If the argument io.Reader is already a Reader with large enough size, it returns
// the underlying Reader.

// NewReaderSize返回了一个新的读取器，这个读取器的缓存大小至少大于制定的大小。
// 如果io.Reader参数已经是一个有足够大缓存的读取器，它就会返回这个Reader了。
func NewReaderSize(rd io.Reader, size int) *Reader

// Buffered returns the number of bytes that can be read from the current buffer.

// Buffered返回当前缓存的可读字节数。
func (b *Reader) Buffered() int

// Peek returns the next n bytes without advancing the reader. The bytes stop being
// valid at the next read call. If Peek returns fewer than n bytes, it also returns
// an error explaining why the read is short. The error is ErrBufferFull if n is
// larger than b's buffer size.

// Peek返回没有读取的下n个字节。在下个读取的调用前，字节是不可见的。如果Peek返回的字节数少于n，
// 它一定会解释为什么读取的字节数段了。如果n比b的缓冲大小更大，返回的错误是ErrBufferFull。
func (b *Reader) Peek(n int) ([]byte, error)

// Read reads data into p. It returns the number of bytes read into p. It calls
// Read at most once on the underlying Reader, hence n may be less than len(p). At
// EOF, the count will be zero and err will be io.EOF.

// Read读取数据到p。 返回读取到p的字节数。
// 底层读取最多只会调用一次Read，因此n会小于len(p)。
// 在EOF之后，调用这个函数返回的会是0和io.Eof。
func (b *Reader) Read(p []byte) (n int, err error)

// ReadByte reads and returns a single byte. If no byte is available, returns an
// error.

// ReadByte读取和回复一个单字节。
// 如果没有字节可以读取，返回一个error。
func (b *Reader) ReadByte() (c byte, err error)

// ReadBytes reads until the first occurrence of delim in the input, returning a
// slice containing the data up to and including the delimiter. If ReadBytes
// encounters an error before finding a delimiter, it returns the data read before
// the error and the error itself (often io.EOF). ReadBytes returns err != nil if
// and only if the returned data does not end in delim. For simple uses, a Scanner
// may be more convenient.

// ReadBytes读取输入到第一次终止符发生的时候，返回的slice包含从当前到终止符的内容（包括终止符）。
// 如果ReadBytes在遇到终止符之前就捕获到一个错误，它就会返回遇到错误之前已经读取的数据，和这个捕获
// 到的错误（经常是
// io.EOF）。当返回的数据没有以终止符结束的时候，ReadBytes返回err != nil。
// 对于简单的使用，或许 Scanner 更方便。
func (b *Reader) ReadBytes(delim byte) (line []byte, err error)

// ReadLine is a low-level line-reading primitive. Most callers should use
// ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
//
// ReadLine tries to return a single line, not including the end-of-line bytes. If
// the line was too long for the buffer then isPrefix is set and the beginning of
// the line is returned. The rest of the line will be returned from future calls.
// isPrefix will be false when returning the last fragment of the line. The
// returned buffer is only valid until the next call to ReadLine. ReadLine either
// returns a non-nil line or it returns an error, never both.
//
// The text returned from ReadLine does not include the line end ("\r\n" or "\n").
// No indication or error is given if the input ends without a final line end.
// Calling UnreadByte after ReadLine will always unread the last byte read
// (possibly a character belonging to the line end) even if that byte is not part
// of the line returned by ReadLine.

// ReadLine是一个底层的原始读取命令。许多调用者或许会使用ReadBytes('\n')或者ReadString('\n')来代替这个方法。
//
// ReadLine尝试返回单个行，不包括行尾的最后一个分隔符。如果一个行大于缓存，调用的时候返回了ifPrefix，
// 就会返回行的头部。行剩余的部分就会在下次调用的时候返回。当调用行的剩余的部分的时候，isPrefix将会设为false，
// 返回的缓存只能在下次调用ReadLine的时候看到。ReadLine会返回了一个非空行，或者返回一个error，
// 但是不会两者都返回。
//
// ReadLine返回的文本不会包含行结尾（"\r\n"或者"\n"）。如果输入没有最终的行结尾的时候，不会返回
// 任何迹象或者错误。在 ReadLine 之后调用 UnreadByte
// 将总是放回读取的最后一个字节
// （可能是属于该行末的字符），即便该字节并非 ReadLine 返回的行的一部分。
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)

// ReadRune reads a single UTF-8 encoded Unicode character and returns the rune and
// its size in bytes. If the encoded rune is invalid, it consumes one byte and
// returns unicode.ReplacementChar (U+FFFD) with a size of 1.

// ReadRune读取单个的UTF-8编码的Unicode字节，并且返回rune和它的字节大小。
// 如果编码的rune是可见的，它消耗一个字节并且返回1字节的unicode.ReplacementChar (U+FFFD)。
func (b *Reader) ReadRune() (r rune, size int, err error)

// ReadSlice reads until the first occurrence of delim in the input, returning a
// slice pointing at the bytes in the buffer. The bytes stop being valid at the
// next read. If ReadSlice encounters an error before finding a delimiter, it
// returns all the data in the buffer and the error itself (often io.EOF).
// ReadSlice fails with error ErrBufferFull if the buffer fills without a delim.
// Because the data returned from ReadSlice will be overwritten by the next I/O
// operation, most clients should use ReadBytes or ReadString instead. ReadSlice
// returns err != nil if and only if line does not end in delim.

// ReadSlice从输入中读取，直到遇到第一个终止符为止，返回一个指向缓存中字节的slice。
// 在下次调用的时候这些字节就是已经被读取了。如果ReadSlice在找到终止符之前遇到了error，
// 它就会返回缓存中所有的数据和错误本身（经常是 io.EOF）。
// 如果在终止符之前缓存已经被充满了，ReadSlice会返回ErrBufferFull错误。
// 由于ReadSlice返回的数据会被下次的I/O操作重写，因此许多的客户端会选择使用ReadBytes或者ReadString代替。
// 当且仅当数据没有以终止符结束的时候，ReadSlice返回err != nil
func (b *Reader) ReadSlice(delim byte) (line []byte, err error)

// ReadString reads until the first occurrence of delim in the input, returning a
// string containing the data up to and including the delimiter. If ReadString
// encounters an error before finding a delimiter, it returns the data read before
// the error and the error itself (often io.EOF). ReadString returns err != nil if
// and only if the returned data does not end in delim. For simple uses, a Scanner
// may be more convenient.

// ReadString读取输入到第一次终止符发生的时候，返回的string包含从当前到终止符的内容（包括终止符）。
// 如果ReadString在遇到终止符之前就捕获到一个错误，它就会返回遇到错误之前已经读取的数据，和这个捕获
// 到的错误（经常是
// io.EOF）。当返回的数据没有以终止符结束的时候，ReadString返回err != nil。
// 对于简单的使用，或许 Scanner 更方便。
func (b *Reader) ReadString(delim byte) (line string, err error)

// Reset discards any buffered data, resets all state, and switches the buffered
// reader to read from r.

// Reset丢弃缓冲中的数据，清除任何错误，将b重设为其下层从r读取数据。
func (b *Reader) Reset(r io.Reader)

// UnreadByte unreads the last byte. Only the most recently read byte can be
// unread.

// UnreadByte将最后的字节标志为未读。只有最后的字节才可以被标志为未读。
func (b *Reader) UnreadByte() error

// UnreadRune unreads the last rune. If the most recent read operation on the
// buffer was not a ReadRune, UnreadRune returns an error. (In this regard it is
// stricter than UnreadByte, which will unread the last byte from any read
// operation.)

// UnreadRune将最后一个rune设置为未读。如果最新的在buffer上的操作不是ReadRune，则UnreadRune
// 就返回一个error。（在这个角度上看，这个函数比UnreadByte更严格，UnreadByte会将最后一个读取
// 的byte设置为未读。）
func (b *Reader) UnreadRune() error

// WriteTo implements io.WriterTo.

// WriteTo实现了io.WriterTo。
func (b *Reader) WriteTo(w io.Writer) (n int64, err error)

// Scanner provides a convenient interface for reading data such as a file of
// newline-delimited lines of text. Successive calls to the Scan method will step
// through the 'tokens' of a file, skipping the bytes between the tokens. The
// specification of a token is defined by a split function of type SplitFunc; the
// default split function breaks the input into lines with line termination
// stripped. Split functions are defined in this package for scanning a file into
// lines, bytes, UTF-8-encoded runes, and space-delimited words. The client may
// instead provide a custom split function.
//
// Scanning stops unrecoverably at EOF, the first I/O error, or a token too large
// to fit in the buffer. When a scan stops, the reader may have advanced
// arbitrarily far past the last token. Programs that need more control over error
// handling or large tokens, or must run sequential scans on a reader, should use
// bufio.Reader instead.

// Scanner类型提供了方便的读取数据的接口，如从换行符分隔的文本里读取每一行。
//
// 成功调用的Scan方法会逐步提供文件的token，
// 跳过token之间的字节。token由SplitFunc类型的分割函数指定；
// 默认的分割函数会将输入分割为多个行，并去掉行尾的换行标志。
// 本包预定义的分割函数可以将文件分割为行、字节、unicode码值、空白分隔的word。
// 调用者可以定制自己的分割函数。
//
// 扫描会在抵达输入流结尾、遇到的第一个I/O错误、token过大不能保存进缓冲时，
// 不可恢复的停止。当扫描停止后，当前读取位置可能会远在最后一个获得的token后面。
// 需要更多对错误管理的控制或token很大，或必须从reader连续扫描的程序， 应使用bufio.Reader代替。
type Scanner struct {
	// contains filtered or unexported fields
}

// NewScanner returns a new Scanner to read from r. The split function defaults to
// ScanLines.

// NewScanner创建并返回一个从r读取数据的Scanner，默认的分割函数是ScanLines。
func NewScanner(r io.Reader) *Scanner

// Bytes returns the most recent token generated by a call to Scan. The underlying
// array may point to data that will be overwritten by a subsequent call to Scan.
// It does no allocation.

// Bytes方法返回最近一次Scan调用生成的token。
// 底层数组指向的数据可能会被下一次Scan的调用重写。
func (s *Scanner) Bytes() []byte

// Err returns the first non-EOF error that was encountered by the Scanner.

// Err返回Scanner遇到的第一个非EOF的错误。
func (s *Scanner) Err() error

// Scan advances the Scanner to the next token, which will then be available
// through the Bytes or Text method. It returns false when the scan stops, either
// by reaching the end of the input or an error. After Scan returns false, the Err
// method will return any error that occurred during scanning, except that if it
// was io.EOF, Err will return nil. Split panics if the split function returns 100
// empty tokens without advancing the input. This is a common error mode for
// scanners.

// Scan方法获取当前位置的token（该token可以通过Bytes或Text方法获得），
// 并让Scanner的扫描位置移动到下一个token。
// 当扫描因为抵达输入流结尾或者遇到错误而停止时，
// 本方法会返回false。在Scan方法返回false后，
// Err方法将返回扫描时遇到的任何错误；除非是io.EOF，此时Err会返回nil。 若 split 函数返回了 100
// 个空标记而没有推进输入，那么它就会派错（panic）。这是 scanner 的一个常见错误。
func (s *Scanner) Scan() bool

// Split sets the split function for the Scanner. If called, it must be called
// before Scan. The default split function is ScanLines.

// Split设置该Scanner的分割函数。本方法必须在Scan之前调用。
func (s *Scanner) Split(split SplitFunc)

// Text returns the most recent token generated by a call to Scan as a newly
// allocated string holding its bytes.

// Bytes方法返回最近一次Scan调用生成的token，
// 会申请创建一个字符串保存token并返回该字符串。
func (s *Scanner) Text() string

// SplitFunc is the signature of the split function used to tokenize the input. The
// arguments are an initial substring of the remaining unprocessed data and a flag,
// atEOF, that reports whether the Reader has no more data to give. The return
// values are the number of bytes to advance the input and the next token to return
// to the user, plus an error, if any. If the data does not yet hold a complete
// token, for instance if it has no newline while scanning lines, SplitFunc can
// return (0, nil, nil) to signal the Scanner to read more data into the slice and
// try again with a longer slice starting at the same point in the input.
//
// If the returned error is non-nil, scanning stops and the error is returned to
// the client.
//
// The function is never called with an empty data slice unless atEOF is true. If
// atEOF is true, however, data may be non-empty and, as always, holds unprocessed
// text.

// SplitFunc类型代表用于对输出作词法分析的分割函数。
//
// 参数data是尚未处理的数据的一个开始部分的切片，
// 参数atEOF表示是否Reader接口不能提供更多的数据。
// 返回值是解析位置前进的字节数，将要返回给调用者的token切片，
// 以及可能遇到的错误。如果数据不足以（保证）生成一个完整的token，
// 例如需要一整行数据但data里没有换行符， SplitFunc可以返回(0, nil,
// nil)来告诉Scanner读取更多的数据
// 写入切片然后用从同一位置起始、长度更长的切片再试一次（调用SplitFunc类型函数）。
//
// 如果返回值err非nil，扫描将终止并将该错误返回给Scanner的调用者。
//
// 除非atEOF为真，永远不会使用空切片data调用SplitFunc类型函数。
// 然而，如果atEOF为真，data却可能是非空的、且包含着未处理的文本。
type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

// Writer implements buffering for an io.Writer object. If an error occurs writing
// to a Writer, no more data will be accepted and all subsequent writes will return
// the error. After all data has been written, the client should call the Flush
// method to guarantee all data has been forwarded to the underlying io.Writer.

// Writer实现了io.Writer对象的缓存。
// 如果在写数据到Writer的时候出现了一个错误，不会再有数据被写进来了，
// 并且所有随后的写操作都会返回error。当所有数据被写入后，客户端应调用 Flush
// 方法以确保所有数据已转为基本的 io.Writer
type Writer struct {
	// contains filtered or unexported fields
}

// NewWriter returns a new Writer whose buffer has the default size.

// NewWriter返回一个新的，有默认尺寸缓存的Writer。
func NewWriter(w io.Writer) *Writer

// NewWriterSize returns a new Writer whose buffer has at least the specified size.
// If the argument io.Writer is already a Writer with large enough size, it returns
// the underlying Writer.

// NewWriterSize返回一个新的Writer，它的缓存一定大于指定的size参数。
// 如果io.Writer参数已经是足够大的有缓存的Writer了，函数就会返回它底层的Writer。
func NewWriterSize(w io.Writer, size int) *Writer

// Available returns how many bytes are unused in the buffer.

// Available返回buffer中有多少的字节数未使用。
func (b *Writer) Available() int

// Buffered returns the number of bytes that have been written into the current
// buffer.

// Buffered返回已经写入到当前缓存的字节数。
func (b *Writer) Buffered() int

// Flush writes any buffered data to the underlying io.Writer.

// Flush将缓存上的所有数据写入到底层的io.Writer中。
func (b *Writer) Flush() error

// ReadFrom implements io.ReaderFrom.

// ReadFrom实现了io.ReaderFrom。
func (b *Writer) ReadFrom(r io.Reader) (n int64, err error)

// Reset discards any unflushed buffered data, clears any error, and resets b to
// write its output to w.
func (b *Writer) Reset(w io.Writer)

// Write writes the contents of p into the buffer. It returns the number of bytes
// written. If nn < len(p), it also returns an error explaining why the write is
// short.

// Writer将p中的内容写入到缓存中。 它返回写入的字节数。 如果nn < len(p),
// 它也会返回错误，用于解释为什么写入的数据会短缺。
func (b *Writer) Write(p []byte) (nn int, err error)

// WriteByte writes a single byte.

// WriterByte写单个字节。
func (b *Writer) WriteByte(c byte) error

// WriteRune writes a single Unicode code point, returning the number of bytes
// written and any error.

// WriteRune写单个的Unicode代码，返回写的字节数，和遇到的错误。
func (b *Writer) WriteRune(r rune) (size int, err error)

// WriteString writes a string. It returns the number of bytes written. If the
// count is less than len(s), it also returns an error explaining why the write is
// short.

// WriteString写一个string。 它返回写入的字节数。
// 如果字节数比len(s)少，它就会返回error来解释为什么写入的数据短缺了。
func (b *Writer) WriteString(s string) (int, error)
