// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package tabwriter implements a write filter (tabwriter.Writer) that translates
// tabbed columns in input into properly aligned text.
//
// The package is using the Elastic Tabstops algorithm described at
// http://nickgravgaard.com/elastictabstops/index.html.

// tabwriter包实现了写入过滤器（tabwriter.Writer），可以将输入的缩进修正为正确的对齐文本。
//
// 本包采用的Elastic
// Tabstops算法参见http://nickgravgaard.com/elastictabstops/index.html
package tabwriter

// Formatting can be controlled with these flags.

// 这些标志用于控制格式化。
//
//	const Escape = '\xff'
//
// 用于包围转义字符，避免该字符被转义；例如字符串"Ignore this tab:
// \xff\t\xff"中的'\t'不被转义，不结束单元；格式化时Escape视为长度1的单字符。
//
// 选择'\xff'是因为该字符不能出现在合法的utf-8序列里。
const (
	// Ignore html tags and treat entities (starting with '&'
	// and ending in ';') as single characters (width = 1).
	FilterHTML uint = 1 << iota

	// Strip Escape characters bracketing escaped text segments
	// instead of passing them through unchanged with the text.
	StripEscape

	// Force right-alignment of cell content.
	// Default is left-alignment.
	AlignRight

	// Handle empty columns as if they were not present in
	// the input in the first place.
	DiscardEmptyColumns

	// Always use tabs for indentation columns (i.e., padding of
	// leading empty cells on the left) independent of padchar.
	TabIndent

	// Print a vertical bar ('|') between columns (after formatting).
	// Discarded columns appear as zero-width columns ("||").
	Debug
)

// To escape a text segment, bracket it with Escape characters. For instance, the
// tab in this string "Ignore this tab: \xff\t\xff" does not terminate a cell and
// constitutes a single character of width one for formatting purposes.
//
// The value 0xff was chosen because it cannot appear in a valid UTF-8 sequence.
const Escape = '\xff'

// A Writer is a filter that inserts padding around tab-delimited columns in its
// input to align them in the output.
//
// The Writer treats incoming bytes as UTF-8 encoded text consisting of cells
// terminated by (horizontal or vertical) tabs or line breaks (newline or formfeed
// characters). Cells in adjacent lines constitute a column. The Writer inserts
// padding as needed to make all cells in a column have the same width, effectively
// aligning the columns. It assumes that all characters have the same width except
// for tabs for which a tabwidth must be specified. Note that cells are
// tab-terminated, not tab-separated: trailing non-tab text at the end of a line
// does not form a column cell.
//
// The Writer assumes that all Unicode code points have the same width; this may
// not be true in some fonts.
//
// If DiscardEmptyColumns is set, empty columns that are terminated entirely by
// vertical (or "soft") tabs are discarded. Columns terminated by horizontal (or
// "hard") tabs are not affected by this flag.
//
// If a Writer is configured to filter HTML, HTML tags and entities are passed
// through. The widths of tags and entities are assumed to be zero (tags) and one
// (entities) for formatting purposes.
//
// A segment of text may be escaped by bracketing it with Escape characters. The
// tabwriter passes escaped text segments through unchanged. In particular, it does
// not interpret any tabs or line breaks within the segment. If the StripEscape
// flag is set, the Escape characters are stripped from the output; otherwise they
// are passed through as well. For the purpose of formatting, the width of the
// escaped text is always computed excluding the Escape characters.
//
// The formfeed character ('\f') acts like a newline but it also terminates all
// columns in the current line (effectively calling Flush). Cells in the next line
// start new columns. Unless found inside an HTML tag or inside an escaped text
// segment, formfeed characters appear as newlines in the output.
//
// The Writer must buffer input internally, because proper spacing of one line may
// depend on the cells in future lines. Clients must call Flush when done calling
// Write.

// Writer是一个过滤器，会在输入的tab划分的列进行填充，在输出中对齐它们。
//
// 它会将输入的序列视为utf-8编码的文本，包含一系列被垂直制表符、水平制表符、换行符、回车符分割的单元。临近的单元组成一列，根据需要填充空格使所有的单元有相同的宽度，高效对齐各列。它假设所有的字符都有相同的宽度，除了tab，tab宽度应该被指定。注意单元以tab截止，而不是被tab分隔，行最后的非tab文本不被视为列的单元。
//
// Writer假设所有的unicode字符有着同样的宽度，这一点其实在很多字体里是错误的。
//
// 如果设置了DiscardEmptyColumns，以垂直制表符结尾的空列会被丢弃，水平制表符截止的空列则不会被影响。
//
// 如果设置了FilterHTML，HTML标签和实体会被放过，标签宽度视为0，实体宽度视为1。文本段可能被转义字符包围，此时tabwriter不会修改该文本段，不会打断段中的任何tab或换行。
//
// 如果设置了StripEscape，则不会计算转义字符的宽度（输出中也会去除转义字符）。
//
// 进纸符'\f'被视为换行，但也会截止当前行的所有列（有效的刷新）；除非在HTML标签内或者在转义文本段内，输出中'\f'都被作为换行。
//
// Writer会在内部缓存输入以便有效的对齐，调用者必须在写完后执行Flush方法。
type Writer struct {
	// contains filtered or unexported fields
}

// NewWriter allocates and initializes a new tabwriter.Writer. The parameters are
// the same as for the Init function.

// 创建并初始化一个tabwriter.Writer，参数用法和Init函数类似。
func NewWriter(output io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) *Writer

// Flush should be called after the last call to Write to ensure that any data
// buffered in the Writer is written to output. Any incomplete escape sequence at
// the end is considered complete for formatting purposes.
func (b *Writer) Flush() (err error)

// A Writer must be initialized with a call to Init. The first parameter (output)
// specifies the filter output. The remaining parameters control the formatting:
//
//	minwidth	minimal cell width including any padding
//	tabwidth	width of tab characters (equivalent number of spaces)
//	padding		padding added to a cell before computing its width
//	padchar		ASCII char used for padding
//			if padchar == '\t', the Writer will assume that the
//			width of a '\t' in the formatted output is tabwidth,
//			and cells are left-aligned independent of align_left
//			(for correct-looking results, tabwidth must correspond
//			to the tab width in the viewer displaying the result)
//	flags		formatting control

// 初始化一个Writer，第一个参数指定格式化后的输出目标，其余的参数控制格式化：
//
//	minwidth 最小单元长度
//	tabwidth tab字符的宽度
//	padding  计算单元宽度时会额外加上它
//	padchar  用于填充的ASCII字符，
//	         如果是'\t'，则Writer会假设tabwidth作为输出中tab的宽度，且单元必然左对齐。
//	flags    格式化控制
func (b *Writer) Init(output io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) *Writer

// Write writes buf to the writer b. The only errors returned are ones encountered
// while writing to the underlying output stream.
func (b *Writer) Write(buf []byte) (n int, err error)
