// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package time provides functionality for measuring and displaying time.
//
// The calendrical calculations always assume a Gregorian calendar.
package time

// These are predefined layouts for use in Time.Format and Time.Parse. The
// reference time used in the layouts is the specific time:
//
//	Mon Jan 2 15:04:05 MST 2006
//
// which is Unix time 1136239445. Since MST is GMT-0700, the reference time can be
// thought of as
//
//	01/02 03:04:05PM '06 -0700
//
// To define your own format, write down what the reference time would look like
// formatted your way; see the values of constants like ANSIC, StampMicro or
// Kitchen for examples. The model is to demonstrate what the reference time looks
// like so that the Format and Parse methods can apply the same transformation to a
// general time value.
//
// Within the format string, an underscore _ represents a space that may be
// replaced by a digit if the following number (a day) has two digits; for
// compatibility with fixed-width Unix time formats.
//
// A decimal point followed by one or more zeros represents a fractional second,
// printed to the given number of decimal places. A decimal point followed by one
// or more nines represents a fractional second, printed to the given number of
// decimal places, with trailing zeros removed. When parsing (only), the input may
// contain a fractional second field immediately after the seconds field, even if
// the layout does not signify its presence. In that case a decimal point followed
// by a maximal series of digits is parsed as a fractional second.
//
// Numeric time zone offsets format as follows:
//
//	-0700  ±hhmm
//	-07:00 ±hh:mm
//
// Replacing the sign in the format with a Z triggers the ISO 8601 behavior of
// printing Z instead of an offset for the UTC zone. Thus:
//
//	Z0700  Z or ±hhmm
//	Z07:00 Z or ±hh:mm
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
	_ = iota
)

// After waits for the duration to elapse and then sends the current time on the
// returned channel. It is equivalent to NewTimer(d).C.
func After(d Duration) <-chan Time

// Sleep pauses the current goroutine for at least the duration d. A negative or
// zero duration causes Sleep to return immediately.
func Sleep(d Duration)

// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only. Useful for clients that have no need to shut down the ticker.
func Tick(d Duration) <-chan Time

// A Duration represents the elapsed time between two instants as an int64
// nanosecond count. The representation limits the largest representable duration
// to approximately 290 years.
type Duration int64

// Common durations. There is no definition for units of Day or larger to avoid
// confusion across daylight savings time zone transitions.
//
// To count the number of units in a Duration, divide:
//
//	second := time.Second
//	fmt.Print(int64(second/time.Millisecond)) // prints 1000
//
// To convert an integer number of units to a Duration, multiply:
//
//	seconds := 10
//	fmt.Print(time.Duration(seconds)*time.Second) // prints 10s
const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

// ParseDuration parses a duration string. A duration string is a possibly signed
// sequence of decimal numbers, each with optional fraction and a unit suffix, such
// as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms",
// "s", "m", "h".
func ParseDuration(s string) (Duration, error)

// Since returns the time elapsed since t. It is shorthand for time.Now().Sub(t).
func Since(t Time) Duration

// Hours returns the duration as a floating point number of hours.
func (d Duration) Hours() float64

// Minutes returns the duration as a floating point number of minutes.
func (d Duration) Minutes() float64

// Nanoseconds returns the duration as an integer nanosecond count.
func (d Duration) Nanoseconds() int64

// Seconds returns the duration as a floating point number of seconds.
func (d Duration) Seconds() float64

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure that
// the leading digit is non-zero. The zero duration formats as 0, with no unit.
func (d Duration) String() string

// A Location maps time instants to the zone in use at that time. Typically, the
// Location represents the collection of time offsets in use in a geographical
// area, such as CEST and CET for central Europe.
type Location struct {
	// contains filtered or unexported fields
}

// Local represents the system's local time zone.
var Local *Location = &localLoc

// UTC represents Universal Coordinated Time (UTC).
var UTC *Location = &utcLoc

// FixedZone returns a Location that always uses the given zone name and offset
// (seconds east of UTC).
func FixedZone(name string, offset int) *Location

// LoadLocation returns the Location with the given name.
//
// If the name is "" or "UTC", LoadLocation returns UTC. If the name is "Local",
// LoadLocation returns Local.
//
// Otherwise, the name is taken to be a location name corresponding to a file in
// the IANA Time Zone database, such as "America/New_York".
//
// The time zone database needed by LoadLocation may not be present on all systems,
// especially non-Unix systems. LoadLocation looks in the directory or uncompressed
// zip file named by the ZONEINFO environment variable, if any, then looks in known
// installation locations on Unix systems, and finally looks in
// $GOROOT/lib/time/zoneinfo.zip.
func LoadLocation(name string) (*Location, error)

// String returns a descriptive name for the time zone information, corresponding
// to the argument to LoadLocation.
func (l *Location) String() string

// A Month specifies a month of the year (January = 1, ...).
type Month int

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

// String returns the English name of the month ("January", "February", ...).
func (m Month) String() string

// ParseError describes a problem parsing a time string.
type ParseError struct {
	Layout     string
	Value      string
	LayoutElem string
	ValueElem  string
	Message    string
}

// Error returns the string representation of a ParseError.
func (e *ParseError) Error() string

// A Ticker holds a channel that delivers `ticks' of a clock at intervals.
type Ticker struct {
	C <-chan Time // The channel on which the ticks are delivered.
	// contains filtered or unexported fields
}

// NewTicker returns a new Ticker containing a channel that will send the time with
// a period specified by the duration argument. It adjusts the intervals or drops
// ticks to make up for slow receivers. The duration d must be greater than zero;
// if not, NewTicker will panic. Stop the ticker to release associated resources.
func NewTicker(d Duration) *Ticker

// Stop turns off a ticker. After Stop, no more ticks will be sent. Stop does not
// close the channel, to prevent a read from the channel succeeding incorrectly.
func (t *Ticker) Stop()

// A Time represents an instant in time with nanosecond precision.
//
// Programs using times should typically store and pass them as values, not
// pointers. That is, time variables and struct fields should be of type time.Time,
// not *time.Time. A Time value can be used by multiple goroutines simultaneously.
//
// Time instants can be compared using the Before, After, and Equal methods. The
// Sub method subtracts two instants, producing a Duration. The Add method adds a
// Time and a Duration, producing a Time.
//
// The zero value of type Time is January 1, year 1, 00:00:00.000000000 UTC. As
// this time is unlikely to come up in practice, the IsZero method gives a simple
// way of detecting a time that has not been initialized explicitly.
//
// Each Time has associated with it a Location, consulted when computing the
// presentation form of the time, such as in the Format, Hour, and Year methods.
// The methods Local, UTC, and In return a Time with a specific location. Changing
// the location in this way changes only the presentation; it does not change the
// instant in time being denoted and therefore does not affect the computations
// described in earlier paragraphs.
//
// Note that the Go == operator compares not just the time instant but also the
// Location. Therefore, Time values should not be used as map or database keys
// without first guaranteeing that the identical Location has been set for all
// values, which can be achieved through use of the UTC or Local method.
type Time struct {
	// contains filtered or unexported fields
}

// Date returns the Time corresponding to
//
//	yyyy-mm-dd hh:mm:ss + nsec nanoseconds
//
// in the appropriate zone for that time in the given location.
//
// The month, day, hour, min, sec, and nsec values may be outside their usual
// ranges and will be normalized during the conversion. For example, October 32
// converts to November 1.
//
// A daylight savings time transition skips or repeats times. For example, in the
// United States, March 13, 2011 2:15am never occurred, while November 6, 2011
// 1:15am occurred twice. In such cases, the choice of time zone, and therefore the
// time, is not well-defined. Date returns a time that is correct in one of the two
// zones involved in the transition, but it does not guarantee which.
//
// Date panics if loc is nil.
func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time

// Now returns the current local time.
func Now() Time

// Parse parses a formatted string and returns the time value it represents. The
// layout defines the format by showing how the reference time, defined to be
//
//	Mon Jan 2 15:04:05 -0700 MST 2006
//
// would be interpreted if it were the value; it serves as an example of the input
// format. The same interpretation will then be made to the input string.
// Predefined layouts ANSIC, UnixDate, RFC3339 and others describe standard and
// convenient representations of the reference time. For more information about the
// formats and the definition of the reference time, see the documentation for
// ANSIC and the other constants defined by this package.
//
// Elements omitted from the value are assumed to be zero or, when zero is
// impossible, one, so parsing "3:04pm" returns the time corresponding to Jan 1,
// year 0, 15:04:00 UTC (note that because the year is 0, this time is before the
// zero Time). Years must be in the range 0000..9999. The day of the week is
// checked for syntax but it is otherwise ignored.
//
// In the absence of a time zone indicator, Parse returns a time in UTC.
//
// When parsing a time with a zone offset like -0700, if the offset corresponds to
// a time zone used by the current location (Local), then Parse uses that location
// and zone in the returned time. Otherwise it records the time as being in a
// fabricated location with time fixed at the given zone offset.
//
// When parsing a time with a zone abbreviation like MST, if the zone abbreviation
// has a defined offset in the current location, then that offset is used. The zone
// abbreviation "UTC" is recognized as UTC regardless of location. If the zone
// abbreviation is unknown, Parse records the time as being in a fabricated
// location with the given zone abbreviation and a zero offset. This choice means
// that such a time can be parsed and reformatted with the same layout losslessly,
// but the exact instant used in the representation will differ by the actual zone
// offset. To avoid such problems, prefer time layouts that use a numeric zone
// offset, or use ParseInLocation.
func Parse(layout, value string) (Time, error)

// ParseInLocation is like Parse but differs in two important ways. First, in the
// absence of time zone information, Parse interprets a time as UTC;
// ParseInLocation interprets the time as in the given location. Second, when given
// a zone offset or abbreviation, Parse tries to match it against the Local
// location; ParseInLocation uses the given location.
func ParseInLocation(layout, value string, loc *Location) (Time, error)

// Unix returns the local Time corresponding to the given Unix time, sec seconds
// and nsec nanoseconds since January 1, 1970 UTC. It is valid to pass nsec outside
// the range [0, 999999999].
func Unix(sec int64, nsec int64) Time

// Add returns the time t+d.
func (t Time) Add(d Duration) Time

// AddDate returns the time corresponding to adding the given number of years,
// months, and days to t. For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does, so, for example,
// adding one month to October 31 yields December 1, the normalized form for
// November 31.
func (t Time) AddDate(years int, months int, days int) Time

// After reports whether the time instant t is after u.
func (t Time) After(u Time) bool

// Before reports whether the time instant t is before u.
func (t Time) Before(u Time) bool

// Clock returns the hour, minute, and second within the day specified by t.
func (t Time) Clock() (hour, min, sec int)

// Date returns the year, month, and day in which t occurs.
func (t Time) Date() (year int, month Month, day int)

// Day returns the day of the month specified by t.
func (t Time) Day() int

// Equal reports whether t and u represent the same time instant. Two times can be
// equal even if they are in different locations. For example, 6:00 +0200 CEST and
// 4:00 UTC are Equal. This comparison is different from using t == u, which also
// compares the locations.
func (t Time) Equal(u Time) bool

// Format returns a textual representation of the time value formatted according to
// layout, which defines the format by showing how the reference time, defined to
// be
//
//	Mon Jan 2 15:04:05 -0700 MST 2006
//
// would be displayed if it were the value; it serves as an example of the desired
// output. The same display rules will then be applied to the time value.
// Predefined layouts ANSIC, UnixDate, RFC3339 and others describe standard and
// convenient representations of the reference time. For more information about the
// formats and the definition of the reference time, see the documentation for
// ANSIC and the other constants defined by this package.
func (t Time) Format(layout string) string

// GobDecode implements the gob.GobDecoder interface.
func (t *Time) GobDecode(data []byte) error

// GobEncode implements the gob.GobEncoder interface.
func (t Time) GobEncode() ([]byte, error)

// Hour returns the hour within the day specified by t, in the range [0, 23].
func (t Time) Hour() int

// ISOWeek returns the ISO 8601 year and week number in which t occurs. Week ranges
// from 1 to 53. Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year
// n-1, and Dec 29 to Dec 31 might belong to week 1 of year n+1.
func (t Time) ISOWeek() (year, week int)

// In returns t with the location information set to loc.
//
// In panics if loc is nil.
func (t Time) In(loc *Location) Time

// IsZero reports whether t represents the zero time instant, January 1, year 1,
// 00:00:00 UTC.
func (t Time) IsZero() bool

// Local returns t with the location set to local time.
func (t Time) Local() Time

// Location returns the time zone information associated with t.
func (t Time) Location() *Location

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (t Time) MarshalBinary() ([]byte, error)

// MarshalJSON implements the json.Marshaler interface. The time is a quoted string
// in RFC 3339 format, with sub-second precision added if present.
func (t Time) MarshalJSON() ([]byte, error)

// MarshalText implements the encoding.TextMarshaler interface. The time is
// formatted in RFC 3339 format, with sub-second precision added if present.
func (t Time) MarshalText() ([]byte, error)

// Minute returns the minute offset within the hour specified by t, in the range
// [0, 59].
func (t Time) Minute() int

// Month returns the month of the year specified by t.
func (t Time) Month() Month

// Nanosecond returns the nanosecond offset within the second specified by t, in
// the range [0, 999999999].
func (t Time) Nanosecond() int

// Round returns the result of rounding t to the nearest multiple of d (since the
// zero time). The rounding behavior for halfway values is to round up. If d <= 0,
// Round returns t unchanged.
func (t Time) Round(d Duration) Time

// Second returns the second offset within the minute specified by t, in the range
// [0, 59].
func (t Time) Second() int

// String returns the time formatted using the format string
//
//	"2006-01-02 15:04:05.999999999 -0700 MST"
func (t Time) String() string

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration will
// be returned. To compute t-d for a duration d, use t.Add(-d).
func (t Time) Sub(u Time) Duration

// Truncate returns the result of rounding t down to a multiple of d (since the
// zero time). If d <= 0, Truncate returns t unchanged.
func (t Time) Truncate(d Duration) Time

// UTC returns t with the location set to UTC.
func (t Time) UTC() Time

// Unix returns t as a Unix time, the number of seconds elapsed since January 1,
// 1970 UTC.
func (t Time) Unix() int64

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in nanoseconds
// cannot be represented by an int64. Note that this means the result of calling
// UnixNano on the zero Time is undefined.
func (t Time) UnixNano() int64

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *Time) UnmarshalBinary(data []byte) error

// UnmarshalJSON implements the json.Unmarshaler interface. The time is expected to
// be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(data []byte) (err error)

// UnmarshalText implements the encoding.TextUnmarshaler interface. The time is
// expected to be in RFC 3339 format.
func (t *Time) UnmarshalText(data []byte) (err error)

// Weekday returns the day of the week specified by t.
func (t Time) Weekday() Weekday

// Year returns the year in which t occurs.
func (t Time) Year() int

// YearDay returns the day of the year specified by t, in the range [1,365] for
// non-leap years, and [1,366] in leap years.
func (t Time) YearDay() int

// Zone computes the time zone in effect at time t, returning the abbreviated name
// of the zone (such as "CET") and its offset in seconds east of UTC.
func (t Time) Zone() (name string, offset int)

// The Timer type represents a single event. When the Timer expires, the current
// time will be sent on C, unless the Timer was created by AfterFunc. A Timer must
// be created with NewTimer or AfterFunc.
type Timer struct {
	C <-chan Time
	// contains filtered or unexported fields
}

// AfterFunc waits for the duration to elapse and then calls f in its own
// goroutine. It returns a Timer that can be used to cancel the call using its Stop
// method.
func AfterFunc(d Duration, f func()) *Timer

// NewTimer creates a new Timer that will send the current time on its channel
// after at least duration d.
func NewTimer(d Duration) *Timer

// Reset changes the timer to expire after duration d. It returns true if the timer
// had been active, false if the timer had expired or been stopped.
func (t *Timer) Reset(d Duration) bool

// Stop prevents the Timer from firing. It returns true if the call stops the
// timer, false if the timer has already expired or been stopped. Stop does not
// close the channel, to prevent a read from the channel succeeding incorrectly.
func (t *Timer) Stop() bool

// A Weekday specifies a day of the week (Sunday = 0, ...).
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// String returns the English name of the day ("Sunday", "Monday", ...).
func (d Weekday) String() string
