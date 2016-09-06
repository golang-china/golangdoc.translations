// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package driver defines interfaces to be implemented by database
// drivers as used by package sql.
//
// Most code should use package sql.

// driver包定义了应被数据库驱动实现的接口，这些接口会被sql包使用。
//
// 绝大多数代码应使用sql包。
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

// Bool is a ValueConverter that converts input values to bools.
//
// The conversion rules are:
//
//     - booleans are returned unchanged
//     - for integer types,
//          1 is true
//          0 is false,
//          other integers are an error
//     - for strings and []byte, same rules as strconv.ParseBool
//     - all other types are an error
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

// DefaultParameterConverter is the default implementation of ValueConverter
// that's used when a Stmt doesn't implement ColumnConverter.
//
// DefaultParameterConverter returns the given value directly if IsValue(value).
// Otherwise integer type are converted to int64, floats to float64, and strings
// to []byte. Other types are an error.
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

// ErrBadConn should be returned by a driver to signal to the sql package that a
// driver.Conn is in a bad state (such as the server having earlier closed the
// connection) and the sql package should retry on a new connection.
//
// To prevent duplicate operations, ErrBadConn should NOT be returned if there's
// a possibility that the database server might have performed the operation.
// Even if the server sends back an error, you shouldn't return ErrBadConn.
var ErrBadConn = errors.New("driver: bad connection")

// ErrSkip may be returned by some optional interfaces' methods to
// indicate at runtime that the fast path is unavailable and the sql
// package should continue as if the optional interface was not
// implemented. ErrSkip is only supported where explicitly
// documented.

// ErrSkip may be returned by some optional interfaces' methods to indicate at
// runtime that the fast path is unavailable and the sql package should continue
// as if the optional interface was not implemented. ErrSkip is only supported
// where explicitly documented.
var ErrSkip = errors.New("driver: skip fast-path; continue as if unimplemented")

// Int32 is a ValueConverter that converts input values to int64,
// respecting the limits of an int32 value.

// Int32 is a ValueConverter that converts input values to int64, respecting the
// limits of an int32 value.
var Int32 int32Type

// ResultNoRows is a pre-defined Result for drivers to return when a DDL
// command (such as a CREATE TABLE) succeeds. It returns an error for both
// LastInsertId and RowsAffected.

// ResultNoRows is a pre-defined Result for drivers to return when a DDL command
// (such as a CREATE TABLE) succeeds. It returns an error for both LastInsertId
// and RowsAffected.
var ResultNoRows noRows

// String is a ValueConverter that converts its input to a string.
// If the value is already a string or []byte, it's unchanged.
// If the value is of another type, conversion to string is done
// with fmt.Sprintf("%v", v).

// String is a ValueConverter that converts its input to a string. If the value
// is already a string or []byte, it's unchanged. If the value is of another
// type, conversion to string is done with fmt.Sprintf("%v", v).
var String stringType

// ColumnConverter may be optionally implemented by Stmt if the
// statement is aware of its own columns' types and can convert from
// any type to a driver Value.

// 如果Stmt有自己的列类型，可以实现ColumnConverter接口，返回值可以将任意类型转换
// 为驱动的Value类型。
type ColumnConverter interface {
	// ColumnConverter returns a ValueConverter for the provided
	// column index. If the type of a specific column isn't known
	// or shouldn't be handled specially, DefaultValueConverter
	// can be returned.
	ColumnConverter(idx int)ValueConverter
}

// Conn is a connection to a database. It is not used concurrently
// by multiple goroutines.
//
// Conn is assumed to be stateful.

// Conn是与数据库的连接。该连接不会被多线程并行使用。连接被假定为具有状态的。
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
	Close()error

	// Begin starts and returns a new transaction.
	Begin() (Tx, error)
}

// Driver is the interface that must be implemented by a database
// driver.

// Driver接口必须被数据库驱动实现。
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

// Execer是一个可选的、可能被Conn接口实现的接口。
//
// 如果一个Conn未实现Execer接口，sql包的DB.Exec会首先准备一个查询，执行状态，然
// 后关闭状态。Exec可能会返回ErrSkip。
type Execer interface {
	Exec(query string, args []Value) (Result, error)
}

// NotNull is a type that implements ValueConverter by disallowing nil
// values but otherwise delegating to another ValueConverter.

// NotNull实现了ValueConverter接口，不允许nil值，否则会将值交给Converter字段处理
// 。
type NotNull struct {
	Converter ValueConverter
}

// Null is a type that implements ValueConverter by allowing nil
// values but otherwise delegating to another ValueConverter.

// Null实现了ValueConverter接口，允许nil值，否则会将值交给Converter字段处理。
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

// Queryer是一个可选的、可能被Conn接口实现的接口。 如果一个Conn未实现Queryer接
// 口，sql包的DB.Query会首先准备一个查询，执行状态， 然后关闭状态。Query可能会返
// 回ErrSkip。
type Queryer interface {
	Query(query string, args []Value) (Rows, error)
}

// Result is the result of a query execution.

// Result是查询执行的结果。
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

// Rows是执行查询得到的结果的迭代器。
type Rows interface {
	// Columns returns the names of the columns. The number of
	// columns of the result is inferred from the length of the
	// slice. If a particular column name isn't known, an empty
	// string should be returned for that entry.
	Columns()[]string

	// Close closes the rows iterator.
	Close()error

	// Next is called to populate the next row of data into
	// the provided slice. The provided slice will be the same
	// size as the Columns() are wide.
	//
	// Next should return io.EOF when there are no more rows.
	Next(dest []Value)error
}

// RowsAffected implements Result for an INSERT or UPDATE operation
// which mutates a number of rows.

// RowsAffected实现了Result接口，用于insert或update操作，这些操作会修改零到多行
// 数据。
type RowsAffected int64

// Stmt is a prepared statement. It is bound to a Conn and not
// used by multiple goroutines concurrently.

// Stmt是准备好的状态。它会绑定到一个连接，不应被多go程同时使用。
type Stmt interface {
	// Close closes the statement.
	//
	// As of Go 1.1, a Stmt will not be closed if it's in use
	// by any queries.
	Close()error

	// NumInput returns the number of placeholder parameters.
	//
	// If NumInput returns >= 0, the sql package will sanity check
	// argument counts from callers and return errors to the caller
	// before the statement's Exec or Query methods are called.
	//
	// NumInput may also return -1, if the driver doesn't know
	// its number of placeholders. In that case, the sql package
	// will not sanity check Exec or Query argument counts.
	NumInput()int

	// Exec executes a query that doesn't return rows, such
	// as an INSERT or UPDATE.
	Exec(args []Value) (Result, error)

	// Query executes a query that may return rows, such as a
	// SELECT.
	Query(args []Value) (Rows, error)
}

// Tx is a transaction.

// Tx是一次事务。
type Tx interface {
	Commit()error
	Rollback()error
}

// Value is a value that drivers must be able to handle.
// It is either nil or an instance of one of these types:
//
//   int64
//   float64
//   bool
//   []byte
//   string
//   time.Time

// Value是驱动必须能处理的值。它要么是nil，要么是如下类型的实例：
//
//         int64
//         float64
//         bool
//         []byte
//         string   [*] Rows.Next不会返回该类型值
//         time.Time
type Value interface {
}

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

// ValueConverter接口提供了ConvertValue方法。
//
// 	driver包提供了各种ValueConverter接口的实现，以保证不同驱动之间的实现和转换的一致性。ValueConverter接口有如下用途：
//
// 	   * 转换sql包提供的Value类型值到数据库指定列的类型，并保证它的匹配，
// 	     例如保证某个int64值满足一个表的uint16列。
//
// 	   * 转换数据库提供的值到驱动的Value类型。
//
// 	   * 在扫描时被sql包用于将驱动的Value类型转换为用户的类型。
type ValueConverter interface {
	// ConvertValue converts a value to a driver Value.
	ConvertValue(v interface{}) (Value, error)
}

// Valuer is the interface providing the Value method.
//
// Types implementing Valuer interface are able to convert
// themselves to a driver Value.

// Valuer是提供Value方法的接口。实现了Valuer接口的类型可以将自身转换为驱动支持的
// Value类型值。
type Valuer interface {
	// Value returns a driver Value.
	Value() (Value, error)
}

// IsScanValue is equivalent to IsValue.
// It exists for compatibility.

// IsScanValue报告v是否是合法的Value扫描类型参数。和IsValue不同，IsScanValue不接
// 受字符串类型。
func IsScanValue(v interface{}) bool

// IsValue reports whether v is a valid Value parameter type.

// IsValue报告v是否是合法的Value类型参数。和IsScanValue不同，IsValue接受字符串类
// 型。
func IsValue(v interface{}) bool

func (n NotNull) ConvertValue(v interface{}) (Value, error)

func (n Null) ConvertValue(v interface{}) (Value, error)

func (RowsAffected) LastInsertId() (int64, error)

func (v RowsAffected) RowsAffected() (int64, error)

