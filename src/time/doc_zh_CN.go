// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package time provides functionality for measuring and displaying time.
//
// The calendrical calculations always assume a Gregorian calendar.

// time包提供了时间的显示和测量用的函数。日历的计算采用的是公历。
package time

import (
	"errors"
	"internal/syscall/windows/registry"
	"runtime"
	"sync"
	"syscall"
)

// These are predefined layouts for use in Time.Format and Time.Parse. The
// reference time used in the layouts is the specific time:
//
//     Mon Jan 2 15:04:05 MST 2006
//
// which is Unix time 1136239445. Since MST is GMT-0700, the reference time can
// be thought of as
//
//     01/02 03:04:05PM '06 -0700
//
// To define your own format, write down what the reference time would look like
// formatted your way; see the values of constants like ANSIC, StampMicro or
// Kitchen for examples. The model is to demonstrate what the reference time
// looks like so that the Format and Parse methods can apply the same
// transformation to a general time value.
//
// Within the format string, an underscore _ represents a space that may be
// replaced by a digit if the following number (a day) has two digits; for
// compatibility with fixed-width Unix time formats.
//
// A decimal point followed by one or more zeros represents a fractional second,
// printed to the given number of decimal places. A decimal point followed by
// one or more nines represents a fractional second, printed to the given number
// of decimal places, with trailing zeros removed. When parsing (only), the
// input may contain a fractional second field immediately after the seconds
// field, even if the layout does not signify its presence. In that case a
// decimal point followed by a maximal series of digits is parsed as a
// fractional second.
//
// Numeric time zone offsets format as follows:
//
//     -0700  ±hhmm
//     -07:00 ±hh:mm
//     -07    ±hh
//
// Replacing the sign in the format with a Z triggers the ISO 8601 behavior of
// printing Z instead of an offset for the UTC zone. Thus:
//
//     Z0700  Z or ±hhmm
//     Z07:00 Z or ±hh:mm
//     Z07    Z or ±hh
//
// The executable example for time.Format demonstrates the working of the layout
// string in detail and is a good reference.
//
// Note that the RFC822, RFC850, and RFC1123 formats should be applied only to
// local times. Applying them to UTC times will use "UTC" as the time zone
// abbreviation, while strictly speaking those RFCs require the use of "GMT" in
// that case. In general RFC1123Z should be used instead of RFC1123 for servers
// that insist on that format, and RFC3339 should be preferred for new
// protocols.

// 这些预定义的版式用于Time.Format和Time.Parse函数。用在版式中的参考时间是：
//
//     Mon Jan 2 15:04:05 MST 2006
//
// 对应的Unix时间是1136239445。因为MST的时区是GMT-0700，参考时间也可以表示为如下
// ：
//
//     01/02 03:04:05PM '06 -0700
//
// 要定义你自己的格式，写下该参考时间应用于你的格式的情况；例子请参见ANSIC、
// StampMicro或Kitchen等常数的值。该模型是为了演示参考时间的格式化效果，如此一来
// Format和Parse方法可以将相同的转换规则用于一个普通的时间值。
//
// 在格式字符串中，用前置的'0'表示一个可以被可以被数字替换的'0'（如果它后面的数
// 字有两位）；使用下划线表示一个可以被数字替换的空格（如果它后面的数字有两位）
// ；以便兼容Unix定长时间格式。
//
// 小数点后跟0到多个'0'，表示秒数的小数部分，输出时会生成和'0'一样多的小数位；小
// 数点后跟0到多个'9'，表示秒数的小数部分，输出时会生成和'9'一样多的小数位但会将
// 拖尾的'0'去掉。（只有）解析时，输入可以在秒字段后面紧跟一个小数部分，即使格式
// 字符串里没有指明该部分。此时，小数点及其后全部的数字都会成为秒的小数部分。
//
// 数字表示的时区格式如下：
//
//     -0700  ±hhmm
//     -07:00 ±hh:mm
//
// 将格式字符串中的负号替换为Z会触发ISO 8601行为（当时区是UTC时，输出Z而不是时区
// 偏移量），这样：
//
//     Z0700  Z or ±hhmm
//     Z07:00 Z or ±hh:mm
const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)

const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// Common durations.  There is no definition for units of Day or larger
// to avoid confusion across daylight savings time zone transitions.
//
// To count the number of units in a Duration, divide:
//     second := time.Second
//     fmt.Print(int64(second/time.Millisecond)) // prints 1000
//
// To convert an integer number of units to a Duration, multiply:
//     seconds := 10
//     fmt.Print(time.Duration(seconds)*time.Second) // prints 10s

// 常用的时间段。没有定义一天或超过一天的单元，以避免夏时制的时区切换的混乱。
//
// 要将Duration类型值表示为某时间单元的个数，用除法：
//
//     second := time.Second
//     fmt.Print(int64(second/time.Millisecond)) // prints 1000
//
// 要将整数个某时间单元表示为Duration类型值，用乘法：
//
//     seconds := 10
//     fmt.Print(time.Duration(seconds)*time.Second) // prints 10s
const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

const (
	_ = iota
)

// Local represents the system's local time zone.

// Local代表系统本地，对应本地时区。
var Local *Location = &localLoc

// UTC represents Universal Coordinated Time (UTC).

// UTC代表通用协调时间，对应零时区。
var UTC *Location = &utcLoc

// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count.  The representation limits the
// largest representable duration to approximately 290 years.

// Duration类型代表两个时间点之间经过的时间，以纳秒为单位。可表示的最长时间段大
// 约290年。
type Duration int64

// A Location maps time instants to the zone in use at that time.
// Typically, the Location represents the collection of time offsets
// in use in a geographical area, such as CEST and CET for central Europe.

// Location 代表一个（关联到某个时间点的）地点，以及该地点所在的时区。
type Location struct {
}

// A Month specifies a month of the year (January = 1, ...).

// Month代表一年的某个月。
type Month int

// ParseError describes a problem parsing a time string.

// ParseError描述解析时间字符串时出现的错误。
type ParseError struct {
	Layout     string
	Value      string
	LayoutElem string
	ValueElem  string
	Message    string
}

// A Ticker holds a channel that delivers `ticks' of a clock
// at intervals.

// Ticker保管一个通道，并每隔一段时间向其传递"tick"。
type Ticker struct {
	C <-chan Time // The channel on which the ticks are delivered.

}

// A Time represents an instant in time with nanosecond precision.
//
// Programs using times should typically store and pass them as values,
// not pointers.  That is, time variables and struct fields should be of
// type time.Time, not *time.Time.  A Time value can be used by
// multiple goroutines simultaneously.
//
// Time instants can be compared using the Before, After, and Equal methods.
// The Sub method subtracts two instants, producing a Duration.
// The Add method adds a Time and a Duration, producing a Time.
//
// The zero value of type Time is January 1, year 1, 00:00:00.000000000 UTC.
// As this time is unlikely to come up in practice, the IsZero method gives
// a simple way of detecting a time that has not been initialized explicitly.
//
// Each Time has associated with it a Location, consulted when computing the
// presentation form of the time, such as in the Format, Hour, and Year methods.
// The methods Local, UTC, and In return a Time with a specific location.
// Changing the location in this way changes only the presentation; it does not
// change the instant in time being denoted and therefore does not affect the
// computations described in earlier paragraphs.
//
// Note that the Go == operator compares not just the time instant but also the
// Location. Therefore, Time values should not be used as map or database keys
// without first guaranteeing that the identical Location has been set for all
// values, which can be achieved through use of the UTC or Local method.

// Time代表一个纳秒精度的时间点。
//
// 程序中应使用Time类型值来保存和传递时间，而不能用指针。就是说，表示时间的变量
// 和字段，应为time.Time类型，而不是*time.Time.类型。一个Time类型值可以被多个go
// 程同时使用。时间点可以使用Before、After和Equal方法进行比较。Sub方法让两个时间
// 点相减，生成一个Duration类型值（代表时间段）。Add方法给一个时间点加上一个时间
// 段，生成一个新的Time类型时间点。
//
// Time零值代表时间点January 1, year 1, 00:00:00.000000000 UTC。因为本时间点一般
// 不会出现在使用中，IsZero方法提供了检验时间是否显式初始化的一个简单途径。
//
// 每一个时间都具有一个地点信息（及对应地点的时区信息），当计算时间的表示格式时
// ，如Format、Hour和Year等方法，都会考虑该信息。Local、UTC和In方法返回一个指定
// 时区（但指向同一时间点）的Time。修改地点/时区信息只是会改变其表示；不会修改被
// 表示的时间点，因此也不会影响其计算。
type Time struct {
}

// The Timer type represents a single event.
// When the Timer expires, the current time will be sent on C,
// unless the Timer was created by AfterFunc.
// A Timer must be created with NewTimer or AfterFunc.

// Timer类型代表单次时间事件。当Timer到期时，当时的时间会被发送给C，除非Timer是
// 被AfterFunc函数创建的。
type Timer struct {
	C <-chan Time
}

// A Weekday specifies a day of the week (Sunday = 0, ...).

// Weekday代表一周的某一天。
type Weekday int

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
// It is equivalent to NewTimer(d).C.

// After会在另一线程经过时间段d后向返回值发送当时的时间。等价于NewTimer(d).C。
func After(d Duration) <-chan Time

// AfterFunc waits for the duration to elapse and then calls f
// in its own goroutine. It returns a Timer that can
// be used to cancel the call using its Stop method.

// AfterFunc另起一个go程等待时间段d过去，然后调用f。它返回一个Timer，可以通过调
// 用其Stop方法来取消等待和对f的调用。
func AfterFunc(d Duration, f func()) *Timer

// Date returns the Time corresponding to
//     yyyy-mm-dd hh:mm:ss + nsec nanoseconds
// in the appropriate zone for that time in the given location.
//
// The month, day, hour, min, sec, and nsec values may be outside
// their usual ranges and will be normalized during the conversion.
// For example, October 32 converts to November 1.
//
// A daylight savings time transition skips or repeats times.
// For example, in the United States, March 13, 2011 2:15am never occurred,
// while November 6, 2011 1:15am occurred twice.  In such cases, the
// choice of time zone, and therefore the time, is not well-defined.
// Date returns a time that is correct in one of the two zones involved
// in the transition, but it does not guarantee which.
//
// Date panics if loc is nil.

// Date返回一个时区为loc、当地时间为：
//
//     year-month-day hour:min:sec + nsec nanoseconds
//
// 的时间点。
//
// month、day、hour、min、sec和nsec的值可能会超出它们的正常范围，在转换前函数会
// 自动将之规范化。如October 32被修正为November 1。
//
// 夏时制的时区切换会跳过或重复时间。如，在美国，March 13, 2011 2:15am从来不会出
// 现，而November 6, 2011 1:15am 会出现两次。此时，时区的选择和时间是没有良好定
// 义的。Date会返回在时区切换的两个时区其中一个时区
//
// 正确的时间，但本函数不会保证在哪一个时区正确。
//
// 如果loc为nil会panic。
func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time

// FixedZone returns a Location that always uses
// the given zone name and offset (seconds east of UTC).

// FixedZone使用给定的地点名name和时间偏移量offset（单位秒）创建并返回一个
// Location
func FixedZone(name string, offset int) *Location

// LoadLocation returns the Location with the given name.
//
// If the name is "" or "UTC", LoadLocation returns UTC.
// If the name is "Local", LoadLocation returns Local.
//
// Otherwise, the name is taken to be a location name corresponding to a file
// in the IANA Time Zone database, such as "America/New_York".
//
// The time zone database needed by LoadLocation may not be
// present on all systems, especially non-Unix systems.
// LoadLocation looks in the directory or uncompressed zip file
// named by the ZONEINFO environment variable, if any, then looks in
// known installation locations on Unix systems,
// and finally looks in $GOROOT/lib/time/zoneinfo.zip.

// LoadLocation返回使用给定的名字创建的Location。
//
// 如果name是""或"UTC"，返回UTC；如果name是"Local"，返回Local；否则name应该是
// IANA时区数据库里有记录的地点名（该数据库记录了地点和对应的时区），如
// "America/New_York"。
//
// LoadLocation函数需要的时区数据库可能不是所有系统都提供，特别是非Unix系统。此
// 时LoadLocation会查找环境变量ZONEINFO指定目录或解压该变量指定的zip文件（如果有
// 该环境变量）；然后查找Unix系统的惯例时区数据安装位置，最后查找
// $GOROOT/lib/time/zoneinfo.zip。
func LoadLocation(name string) (*Location, error)

// NewTicker returns a new Ticker containing a channel that will send the
// time with a period specified by the duration argument.
// It adjusts the intervals or drops ticks to make up for slow receivers.
// The duration d must be greater than zero; if not, NewTicker will panic.
// Stop the ticker to release associated resources.

// NewTicker返回一个新的Ticker，该Ticker包含一个通道字段，并会每隔时间段d就向该
// 通道发送当时的时间。它会调整时间间隔或者丢弃tick信息以适应反应慢的接收者。如
// 果d<=0会panic。关闭该Ticker可以释放相关资源。
func NewTicker(d Duration) *Ticker

// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.

// NewTimer创建一个Timer，它会在最少过去时间段d后到期，向其自身的C字段发送当时的
// 时间。
func NewTimer(d Duration) *Timer

// Now returns the current local time.
func Now() Time

// Parse parses a formatted string and returns the time value it represents. The
// layout defines the format by showing how the reference time, defined to be
//
//     Mon Jan 2 15:04:05 -0700 MST 2006
//
// would be interpreted if it were the value; it serves as an example of the
// input format. The same interpretation will then be made to the input string.
//
// Predefined layouts ANSIC, UnixDate, RFC3339 and others describe standard and
// convenient representations of the reference time. For more information about
// the formats and the definition of the reference time, see the documentation
// for ANSIC and the other constants defined by this package. Also, the
// executable example for time.Format demonstrates the working of the layout
// string in detail and is a good reference.
//
// Elements omitted from the value are assumed to be zero or, when zero is
// impossible, one, so parsing "3:04pm" returns the time corresponding to Jan 1,
// year 0, 15:04:00 UTC (note that because the year is 0, this time is before
// the zero Time). Years must be in the range 0000..9999. The day of the week is
// checked for syntax but it is otherwise ignored.
//
// In the absence of a time zone indicator, Parse returns a time in UTC.
//
// When parsing a time with a zone offset like -0700, if the offset corresponds
// to a time zone used by the current location (Local), then Parse uses that
// location and zone in the returned time. Otherwise it records the time as
// being in a fabricated location with time fixed at the given zone offset.
//
// No checking is done that the day of the month is within the month's valid
// dates; any one- or two-digit value is accepted. For example February 31 and
// even February 99 are valid dates, specifying dates in March and May. This
// behavior is consistent with time.Date.
//
// When parsing a time with a zone abbreviation like MST, if the zone
// abbreviation has a defined offset in the current location, then that offset
// is used. The zone abbreviation "UTC" is recognized as UTC regardless of
// location. If the zone abbreviation is unknown, Parse records the time as
// being in a fabricated location with the given zone abbreviation and a zero
// offset. This choice means that such a time can be parsed and reformatted with
// the same layout losslessly, but the exact instant used in the representation
// will differ by the actual zone offset. To avoid such problems, prefer time
// layouts that use a numeric zone offset, or use ParseInLocation.

// Parse parses a formatted string and returns the time value it represents. The
// layout defines the format by showing how the reference time, defined to be
//
//     Mon Jan 2 15:04:05 -0700 MST 2006
//
// would be interpreted if it were the value; it serves as an example of the
// input format. The same interpretation will then be made to the input string.
// Predefined layouts ANSIC, UnixDate, RFC3339 and others describe standard and
// convenient representations of the reference time. For more information about
// the formats and the definition of the reference time, see the documentation
// for ANSIC and the other constants defined by this package.
//
// Elements omitted from the value are assumed to be zero or, when zero is
// impossible, one, so parsing "3:04pm" returns the time corresponding to Jan 1,
// year 0, 15:04:00 UTC (note that because the year is 0, this time is before
// the zero Time). Years must be in the range 0000..9999. The day of the week is
// checked for syntax but it is otherwise ignored.
//
// In the absence of a time zone indicator, Parse returns a time in UTC.
//
// When parsing a time with a zone offset like -0700, if the offset corresponds
// to a time zone used by the current location (Local), then Parse uses that
// location and zone in the returned time. Otherwise it records the time as
// being in a fabricated location with time fixed at the given zone offset.
//
// When parsing a time with a zone abbreviation like MST, if the zone
// abbreviation has a defined offset in the current location, then that offset
// is used. The zone abbreviation "UTC" is recognized as UTC regardless of
// location. If the zone abbreviation is unknown, Parse records the time as
// being in a fabricated location with the given zone abbreviation and a zero
// offset. This choice means that such a time can be parsed and reformatted with
// the same layout losslessly, but the exact instant used in the representation
// will differ by the actual zone offset. To avoid such problems, prefer time
// layouts that use a numeric zone offset, or use ParseInLocation.
func Parse(layout, value string) (Time, error)

// ParseDuration parses a duration string. A duration string is a possibly
// signed sequence of decimal numbers, each with optional fraction and a unit
// suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us"
// (or "µs"), "ms", "s", "m", "h".
func ParseDuration(s string) (Duration, error)

// ParseInLocation is like Parse but differs in two important ways. First, in
// the absence of time zone information, Parse interprets a time as UTC;
// ParseInLocation interprets the time as in the given location. Second, when
// given a zone offset or abbreviation, Parse tries to match it against the
// Local location; ParseInLocation uses the given location.
func ParseInLocation(layout, value string, loc *Location) (Time, error)

// Since returns the time elapsed since t. It is shorthand for
// time.Now().Sub(t).
func Since(t Time) Duration

// Sleep pauses the current goroutine for at least the duration d.
// A negative or zero duration causes Sleep to return immediately.

// Sleep阻塞当前go程至少d代表的时间段。d<=0时，Sleep会立刻返回。
func Sleep(d Duration)

// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only. While Tick is useful for clients that have no need to shut down
// the Ticker, be aware that without a way to shut it down the underlying
// Ticker cannot be recovered by the garbage collector; it "leaks".

// Tick是NewTicker的封装，只提供对Ticker的通道的访问。如果不需要关闭Ticker，本函
// 数就很方便。
func Tick(d Duration) <-chan Time

// Unix returns the local Time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
// It is valid to pass nsec outside the range [0, 999999999].
// Not all sec values have a corresponding time value. One such
// value is 1<<63-1 (the largest int64 value).

// Unix returns the local Time corresponding to the given Unix time, sec seconds
// and nsec nanoseconds since January 1, 1970 UTC. It is valid to pass nsec
// outside the range [0, 999999999].
func Unix(sec int64, nsec int64) Time

// String returns a descriptive name for the time zone information,
// corresponding to the argument to LoadLocation.

// String返回对时区信息的描述，返回值绑定为LoadLocation或FixedZone函数创建l时的
// name参数。
func (*Location) String() string

// Error returns the string representation of a ParseError.

// Error返回ParseError的字符串表示。
func (*ParseError) Error() string

// Stop turns off a ticker. After Stop, no more ticks will be sent. Stop does
// not close the channel, to prevent a read from the channel succeeding
// incorrectly.

// Stop关闭一个Ticker。在关闭后，将不会发送更多的tick信息。Stop不会关闭通道t.C，
// 以避免从该通道的读取不正确的成功。
func (*Ticker) Stop()

// GobDecode implements the gob.GobDecoder interface.
func (*Time) GobDecode(data []byte) error

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (*Time) UnmarshalBinary(data []byte) error

// UnmarshalJSON implements the json.Unmarshaler interface. The time is expected
// to be a quoted string in RFC 3339 format.
func (*Time) UnmarshalJSON(data []byte) (err error)

// UnmarshalText implements the encoding.TextUnmarshaler interface. The time is
// expected to be in RFC 3339 format.
func (*Time) UnmarshalText(data []byte) (err error)

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active, false if the timer had
// expired or been stopped.

// Reset使t重新开始计时，（本方法返回后再）等待时间段d过去后到期。如果调用时t还
// 在等待中会返回真；如果t已经到期或者被停止了会返回假。
func (*Timer) Reset(d Duration) bool

// Stop prevents the Timer from firing. It returns true if the call stops the
// timer, false if the timer has already expired or been stopped. Stop does not
// close the channel, to prevent a read from the channel succeeding incorrectly.

// Stop停止Timer的执行。如果停止了t会返回真；如果t已经被停止或者过期了会返回假。
// Stop不会关闭通道t.C，以避免从该通道的读取不正确的成功。
func (*Timer) Stop() bool

// Hours returns the duration as a floating point number of hours.
func (Duration) Hours() float64

// Minutes returns the duration as a floating point number of minutes.
func (Duration) Minutes() float64

// Nanoseconds returns the duration as an integer nanosecond count.
func (Duration) Nanoseconds() int64

// Seconds returns the duration as a floating point number of seconds.
func (Duration) Seconds() float64

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted.  As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero.  The zero duration formats as 0,
// with no unit.

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero. The zero duration formats as 0, with no
// unit.
func (Duration) String() string

// String returns the English name of the month ("January", "February", ...).
func (Month) String() string

// Add returns the time t+d.
func (Time) Add(d Duration) Time

// AddDate returns the time corresponding to adding the given number of years,
// months, and days to t. For example, AddDate(-1, 2, 3) applied to January 1,
// 2011 returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does, so, for
// example, adding one month to October 31 yields December 1, the normalized
// form for November 31.
func (Time) AddDate(years int, months int, days int) Time

// After reports whether the time instant t is after u.
func (Time) After(u Time) bool

// Before reports whether the time instant t is before u.
func (Time) Before(u Time) bool

// Clock returns the hour, minute, and second within the day specified by t.
func (Time) Clock() (hour, min, sec int)

// Date returns the year, month, and day in which t occurs.
func (Time) Date() (year int, month Month, day int)

// Day returns the day of the month specified by t.
func (Time) Day() int

// Equal reports whether t and u represent the same time instant. Two times can
// be equal even if they are in different locations. For example, 6:00 +0200
// CEST and 4:00 UTC are Equal. This comparison is different from using t == u,
// which also compares the locations.
func (Time) Equal(u Time) bool

// Format returns a textual representation of the time value formatted
// according to layout, which defines the format by showing how the reference
// time, defined to be
//     Mon Jan 2 15:04:05 -0700 MST 2006
// would be displayed if it were the value; it serves as an example of the
// desired output. The same display rules will then be applied to the time
// value.
//
// A fractional second is represented by adding a period and zeros
// to the end of the seconds section of layout string, as in "15:04:05.000"
// to format a time stamp with millisecond precision.
//
// Predefined layouts ANSIC, UnixDate, RFC3339 and others describe standard
// and convenient representations of the reference time. For more information
// about the formats and the definition of the reference time, see the
// documentation for ANSIC and the other constants defined by this package.

// Format returns a textual representation of the time value formatted according
// to layout, which defines the format by showing how the reference time,
// defined to be
//
//     Mon Jan 2 15:04:05 -0700 MST 2006
//
// would be displayed if it were the value; it serves as an example of the
// desired output. The same display rules will then be applied to the time
// value. Predefined layouts ANSIC, UnixDate, RFC3339 and others describe
// standard and convenient representations of the reference time. For more
// information about the formats and the definition of the reference time, see
// the documentation for ANSIC and the other constants defined by this package.
func (Time) Format(layout string) string

// GobEncode implements the gob.GobEncoder interface.
func (Time) GobEncode() ([]byte, error)

// Hour returns the hour within the day specified by t, in the range [0, 23].
func (Time) Hour() int

// ISOWeek returns the ISO 8601 year and week number in which t occurs. Week
// ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to week 52 or 53
// of year n-1, and Dec 29 to Dec 31 might belong to week 1 of year n+1.
func (Time) ISOWeek() (year, week int)

// In returns t with the location information set to loc.
//
// In panics if loc is nil.
func (Time) In(loc *Location) Time

// IsZero reports whether t represents the zero time instant, January 1, year 1,
// 00:00:00 UTC.
func (Time) IsZero() bool

// Local returns t with the location set to local time.
func (Time) Local() Time

// Location returns the time zone information associated with t.
func (Time) Location() *Location

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (Time) MarshalBinary() ([]byte, error)

// MarshalJSON implements the json.Marshaler interface. The time is a quoted
// string in RFC 3339 format, with sub-second precision added if present.
func (Time) MarshalJSON() ([]byte, error)

// MarshalText implements the encoding.TextMarshaler interface. The time is
// formatted in RFC 3339 format, with sub-second precision added if present.
func (Time) MarshalText() ([]byte, error)

// Minute returns the minute offset within the hour specified by t, in the range
// [0, 59].
func (Time) Minute() int

// Month returns the month of the year specified by t.
func (Time) Month() Month

// Nanosecond returns the nanosecond offset within the second specified by t, in
// the range [0, 999999999].
func (Time) Nanosecond() int

// Round returns the result of rounding t to the nearest multiple of d (since
// the zero time). The rounding behavior for halfway values is to round up. If d
// <= 0, Round returns t unchanged.
func (Time) Round(d Duration) Time

// Second returns the second offset within the minute specified by t, in the
// range [0, 59].
func (Time) Second() int

// String returns the time formatted using the format string
//     "2006-01-02 15:04:05.999999999 -0700 MST"

// String returns the time formatted using the format string
//
//     "2006-01-02 15:04:05.999999999 -0700 MST"
func (Time) String() string

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned. To compute t-d for a duration d, use t.Add(-d).
func (Time) Sub(u Time) Duration

// Truncate returns the result of rounding t down to a multiple of d (since the
// zero time). If d <= 0, Truncate returns t unchanged.
func (Time) Truncate(d Duration) Time

// UTC returns t with the location set to UTC.
func (Time) UTC() Time

// Unix returns t as a Unix time, the number of seconds elapsed since January 1,
// 1970 UTC.
func (Time) Unix() int64

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in nanoseconds
// cannot be represented by an int64. Note that this means the result of calling
// UnixNano on the zero Time is undefined.
func (Time) UnixNano() int64

// Weekday returns the day of the week specified by t.
func (Time) Weekday() Weekday

// Year returns the year in which t occurs.
func (Time) Year() int

// YearDay returns the day of the year specified by t, in the range [1,365] for
// non-leap years, and [1,366] in leap years.
func (Time) YearDay() int

// Zone computes the time zone in effect at time t, returning the abbreviated
// name of the zone (such as "CET") and its offset in seconds east of UTC.
func (Time) Zone() (name string, offset int)

// String returns the English name of the day ("Sunday", "Monday", ...).

// String返回该日（周几）的英文名（"Sunday"、"Monday"，……）
func (Weekday) String() string
