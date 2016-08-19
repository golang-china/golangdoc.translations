// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package driver defines interfaces to be implemented by database
// drivers as used by package sql.
//
// Most code should use package sql.

// Package driver defines interfaces to be implemented by database
// drivers as used by package sql.
//
// Most code should use package sql.
package driver

import (
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "time"
)

// Bool is a ValueConverter that converts input values to bools.
//
// The conversion rules are:
//  - booleans are returned unchanged
//  - for integer types,
//       1 is true
//       0 is false,
//       other integers are an error
//  - for strings and []byte, same rules as strconv.ParseBool
//  - all other types are an error
var Bool boolType


// DefaultParameterConverter is the default implementation of
// ValueConverter that's used when a Stmt doesn't implement
// ColumnConverter.
//
// DefaultParameterConverter returns its argument directly if
// IsValue(arg). Otherwise, if the argument implements Valuer, its
// Value method is used to return a Value. As a fallback, the provided
// argument's underlying type is used to convert it to a Value:
// underlying integer types are converted to int64, floats to float64,
// and strings to []byte. If the argument is a nil pointer,
// ConvertValue returns a nil Value. If the argument is a non-nil
// pointer, it is dereferenced and ConvertValue is called
// recursively. Other types are an error.
var DefaultParameterConverter defaultConverter


// ErrBadConn should be returned by a driver to signal to the sql
// package that a driver.Conn is in a bad state (such as the server
// having earlier closed the connection) and the sql package should
// retry on a new connection.
//
// To prevent duplicate operations, ErrBadConn should NOT be returned
// if there's a possibility that the database server might have
// performed the operation. Even if the server sends back an error,
// you shouldn't return ErrBadConn.
var ErrBadConn = errors.New("driver: bad connection")


// ErrSkip may be returned by some optional interfaces' methods to
// indicate at runtime that the fast path is unavailable and the sql
// package should continue as if the optional interface was not
// implemented. ErrSkip is only supported where explicitly
// documented.
var ErrSkip = errors.New("driver: skip fast-path; continue as if unimplemented")


// Int32 is a ValueConverter that converts input values to int64,
// respecting the limits of an int32 value.
var Int32 int32Type


// ResultNoRows is a pre-defined Result for drivers to return when a DDL
// command (such as a CREATE TABLE) succeeds. It returns an error for both
// LastInsertId and RowsAffected.
var ResultNoRows noRows


// String is a ValueConverter that converts its input to a string.
// If the value is already a string or []byte, it's unchanged.
// If the value is of another type, conversion to string is done
// with fmt.Sprintf("%v", v).
var String stringType



var _ Result = RowsAffected(0)



var _ ValueConverter = boolType{}



var _ Result = noRows{}



var _ ValueConverter = int32Type{}



var _ ValueConverter = defaultConverter{}


// ColumnConverter may be optionally implemented by Stmt if the
// statement is aware of its own columns' types and can convert from
// any type to a driver Value.
type ColumnConverter interface {
	// ColumnConverter returns a ValueConverter for the provided
	// column index. If the type of a specific column isn't known
	// or shouldn't be handled specially, DefaultValueConverter
	// can be returned.
	ColumnConverter(idx int) ValueConverter
}


// Conn is a connection to a database. It is not used concurrently
// by multiple goroutines.
//
// Conn is assumed to be stateful.
type Conn interface {
	// Prepare returns a prepared statement, bound to this connection.
	Prepare(query string) (Stmt, error)

	// Close invalidates and potentially stops any current
	// prepared statements and transactions, marking this
	// connection as no longer in use.
	//
	// Because the sql package maintains a free pool of
	// connections and only calls Close when there's a surplus of
	// idle connections, it shouldn't be necessary for drivers to
	// do their own connection caching.
	Close() error

	// Begin starts and returns a new transaction.
	Begin() (Tx, error)
}


// Driver is the interface that must be implemented by a database
// driver.
type Driver interface {
	// Open returns a new connection to the database.
	// The name is a string in a driver-specific format.
	//
	// Open may return a cached connection (one previously
	// closed), but doing so is unnecessary; the sql package
	// maintains a pool of idle connections for efficient re-use.
	//
	// The returned connection is only used by one goroutine at a
	// time.
	Open(name string) (Conn, error)
}


// Execer is an optional interface that may be implemented by a Conn.
//
// If a Conn does not implement Execer, the sql package's DB.Exec will
// first prepare a query, execute the statement, and then close the
// statement.
//
// Exec may return ErrSkip.
type Execer interface {
	Exec(query string, args []Value) (Result, error)
}


// NotNull is a type that implements ValueConverter by disallowing nil
// values but otherwise delegating to another ValueConverter.
type NotNull struct {
	Converter ValueConverter
}


// Null is a type that implements ValueConverter by allowing nil
// values but otherwise delegating to another ValueConverter.
type Null struct {
	Converter ValueConverter
}


// Queryer is an optional interface that may be implemented by a Conn.
//
// If a Conn does not implement Queryer, the sql package's DB.Query will
// first prepare a query, execute the statement, and then close the
// statement.
//
// Query may return ErrSkip.
type Queryer interface {
	Query(query string, args []Value) (Rows, error)
}


// Result is the result of a query execution.
type Result interface {
	// LastInsertId returns the database's auto-generated ID
	// after, for example, an INSERT into a table with primary
	// key.
	LastInsertId() (int64, error)

	// RowsAffected returns the number of rows affected by the
	// query.
	RowsAffected() (int64, error)
}


// Rows is an iterator over an executed query's results.
type Rows interface {
	// Columns returns the names of the columns. The number of
	// columns of the result is inferred from the length of the
	// slice. If a particular column name isn't known, an empty
	// string should be returned for that entry.
	Columns() []string

	// Close closes the rows iterator.
	Close() error

	// Next is called to populate the next row of data into
	// the provided slice. The provided slice will be the same
	// size as the Columns() are wide.
	//
	// Next should return io.EOF when there are no more rows.
	Next(dest []Value) error
}


// RowsAffected implements Result for an INSERT or UPDATE operation
// which mutates a number of rows.
type RowsAffected int64


// Stmt is a prepared statement. It is bound to a Conn and not
// used by multiple goroutines concurrently.
type Stmt interface {
	// Close closes the statement.
	//
	// As of Go 1.1, a Stmt will not be closed if it's in use
	// by any queries.
	Close() error

	// NumInput returns the number of placeholder parameters.
	//
	// If NumInput returns >= 0, the sql package will sanity check
	// argument counts from callers and return errors to the caller
	// before the statement's Exec or Query methods are called.
	//
	// NumInput may also return -1, if the driver doesn't know
	// its number of placeholders. In that case, the sql package
	// will not sanity check Exec or Query argument counts.
	NumInput() int

	// Exec executes a query that doesn't return rows, such
	// as an INSERT or UPDATE.
	Exec(args []Value) (Result, error)

	// Query executes a query that may return rows, such as a
	// SELECT.
	Query(args []Value) (Rows, error)
}


// Tx is a transaction.
type Tx interface {
	Commit() error
	Rollback() error
}


// Value is a value that drivers must be able to handle.
// It is either nil or an instance of one of these types:
//
//   int64
//   float64
//   bool
//   []byte
//   string   [*] everywhere except from Rows.Next.
//   time.Time

// Value is a value that drivers must be able to handle.
// It is either nil or an instance of one of these types:
//
//   int64
//   float64
//   bool
//   []byte
//   string
//   time.Time
type Value interface{}


// ValueConverter is the interface providing the ConvertValue method.
//
// Various implementations of ValueConverter are provided by the
// driver package to provide consistent implementations of conversions
// between drivers.  The ValueConverters have several uses:
//
//  * converting from the Value types as provided by the sql package
//    into a database table's specific column type and making sure it
//    fits, such as making sure a particular int64 fits in a
//    table's uint16 column.
//
//  * converting a value as given from the database into one of the
//    driver Value types.
//
//  * by the sql package, for converting from a driver's Value type
//    to a user's type in a scan.

// ValueConverter is the interface providing the ConvertValue method.
//
// Various implementations of ValueConverter are provided by the
// driver package to provide consistent implementations of conversions
// between drivers. The ValueConverters have several uses:
//
//  * converting from the Value types as provided by the sql package
//    into a database table's specific column type and making sure it
//    fits, such as making sure a particular int64 fits in a
//    table's uint16 column.
//
//  * converting a value as given from the database into one of the
//    driver Value types.
//
//  * by the sql package, for converting from a driver's Value type
//    to a user's type in a scan.
type ValueConverter interface {
	// ConvertValue converts a value to a driver Value.
	ConvertValue(v interface{}) (Value, error)
}


// Valuer is the interface providing the Value method.
//
// Types implementing Valuer interface are able to convert
// themselves to a driver Value.
type Valuer interface {
	// Value returns a driver Value.
	Value() (Value, error)
}


// IsScanValue reports whether v is a valid Value scan type.
// Unlike IsValue, IsScanValue does not permit the string type.

// IsScanValue is equivalent to IsValue.
// It exists for compatibility.
func IsScanValue(v interface{}) bool

// IsValue reports whether v is a valid Value parameter type.
// Unlike IsScanValue, IsValue permits the string type.

// IsValue reports whether v is a valid Value parameter type.
func IsValue(v interface{}) bool

func (NotNull) ConvertValue(v interface{}) (Value, error)

func (Null) ConvertValue(v interface{}) (Value, error)

func (RowsAffected) LastInsertId() (int64, error)

func (RowsAffected) RowsAffected() (int64, error)

