// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package sql provides a generic interface around SQL (or SQL-like) databases.
//
// The sql package must be used in conjunction with a database driver. See
// http://golang.org/s/sqldrivers for a list of drivers.
//
// For more usage examples, see the wiki page at http://golang.org/s/sqlwiki.

// sql 包提供了通用的SQL（或类SQL）数据库接口.
//
// sql 包必须与数据库驱动结合使用。驱动列表见 http://golang.org/s/sqldrivers。
//
// 更多使用范例见 http://golang.org/s/sqlwiki 的维基页面。
package sql

// ErrNoRows is returned by Scan when QueryRow doesn't return a row. In such a
// case, QueryRow returns a placeholder *Row value that defers this error until a
// Scan.

// ErrNoRows是QueryRow的时候，当没有返回任何数据，Scan会返回的错误。
// 在这种情况下，QueryRow会返回一个*Row的标示符，直到调用Scan的时候才返回这个error。
var ErrNoRows = errors.New("sql: no rows in result set")

var ErrTxDone = errors.New("sql: Transaction has already been committed or rolled back")

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string

// Register makes a database driver available by the provided name. If Register is
// called twice with the same name or if driver is nil, it panics.

// Register使得数据库驱动可以使用事先定义好的名字使用。
// 如果使用同样的名字注册，或者是注册的的sql驱动是空的，Register会panic。
func Register(name string, driver driver.Driver)

// DB is a database handle representing a pool of zero or more underlying
// connections. It's safe for concurrent use by multiple goroutines.
//
// The sql package creates and frees connections automatically; it also maintains a
// free pool of idle connections. If the database has a concept of per-connection
// state, such state can only be reliably observed within a transaction. Once
// DB.Begin is called, the returned Tx is bound to a single connection. Once Commit
// or Rollback is called on the transaction, that transaction's connection is
// returned to DB's idle connection pool. The pool size can be controlled with
// SetMaxIdleConns.

// DB is a database handle representing a pool of zero or more underlying
// connections. It's safe for concurrent use by multiple goroutines.
//
// The sql package creates and frees connections automatically; it also maintains a
// free pool of idle connections. If the database has a concept of per-connection
// state, such state can only be reliably observed within a transaction. Once
// DB.Begin is called, the returned Tx is bound to a single connection. Once Commit
// or Rollback is called on the transaction, that transaction's connection is
// returned to DB's idle connection pool. The pool size can be controlled with
// SetMaxIdleConns. TODO：待译
type DB struct {
	// contains filtered or unexported fields
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a database name
// and connection information.
//
// Most users will open a database via a driver-specific connection helper function
// that returns a *DB. No database drivers are included in the Go standard library.
// See http://golang.org/s/sqldrivers for a list of third-party drivers.
//
// Open may just validate its arguments without creating a connection to the
// database. To verify that the data source name is valid, call Ping.
//
// The returned DB is safe for concurrent use by multiple goroutines and maintains
// its own pool of idle connections. Thus, the Open function should be called just
// once. It is rarely necessary to close a DB.

// Open打开一个数据库，这个数据库是由其驱动名称和驱动制定的数据源信息打开的，这个数据源信息通常
// 是由至少一个数据库名字和连接信息组成的。
//
// 多数用户通过指定的驱动连接辅助函数来打开一个数据库。打开数据库之后会返回*DB。
//
// TODO：待译
func Open(driverName, dataSourceName string) (*DB, error)

// Begin starts a transaction. The isolation level is dependent on the driver.

// Begin开始一个事务。事务的隔离级别是由驱动决定的。
func (db *DB) Begin() (*Tx, error)

// Close closes the database, releasing any open resources.
//
// It is rare to Close a DB, as the DB handle is meant to be long-lived and shared
// between many goroutines.

// Close关闭数据库，释放一些使用中的资源。 TODO: 待译
func (db *DB) Close() error

// Driver returns the database's underlying driver.

// Driver返回了数据库的底层驱动。
func (db *DB) Driver() driver.Driver

// Exec executes a query without returning any rows. The args are for any
// placeholder parameters in the query.

// Exec 执行query操作，而不返回任何行。 args
// 为查询中的任意占位符形参。
func (db *DB) Exec(query string, args ...interface{}) (Result, error)

// Ping verifies a connection to the database is still alive, establishing a
// connection if necessary.

// Ping verifies a connection to the database is still alive, establishing a
// connection if necessary. TODO：待译
func (db *DB) Ping() error

// Prepare creates a prepared statement for later queries or executions. Multiple
// queries or executions may be run concurrently from the returned statement.

// Prepare
// 为以后的查询或执行操作事先创建了语句。
// 多个查询或执行操作可在返回的语句中并发地运行。
func (db *DB) Prepare(query string) (*Stmt, error)

// Query executes a query that returns rows, typically a SELECT. The args are for
// any placeholder parameters in the query.

// Query执行了一个有返回行的查询操作，比如SELECT。 args
// 形参为该查询中的任何占位符。
func (db *DB) Query(query string, args ...interface{}) (*Rows, error)

// QueryRow executes a query that is expected to return at most one row. QueryRow
// always return a non-nil value. Errors are deferred until Row's Scan method is
// called.

// QueryRow执行一个至多只返回一行记录的查询操作。
// QueryRow总是返回一个非空值。Error只会在调用行的Scan方法的时候才返回。
func (db *DB) QueryRow(query string, args ...interface{}) *Row

// SetMaxIdleConns sets the maximum number of connections in the idle connection
// pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns then the
// new MaxIdleConns will be reduced to match the MaxOpenConns limit
//
// If n <= 0, no idle connections are retained.

// SetMaxIdleConns sets the maximum number of connections in the idle connection
// pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns then the
// new MaxIdleConns will be reduced to match the MaxOpenConns limit
//
// If n <= 0, no idle connections are retained. TODO：待译
func (db *DB) SetMaxIdleConns(n int)

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new MaxOpenConns
// limit
//
// If n <= 0, then there is no limit on the number of open connections. The default
// is 0 (unlimited).
func (db *DB) SetMaxOpenConns(n int)

// NullBool represents a bool that may be null. NullBool implements the Scanner
// interface so it can be used as a scan destination, similar to NullString.

// NullBool代表了可空的bool类型。
// NullBool实现了Scanner接口，所以它和NullString一样可以被当做scan的目标变量。
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface.

// Scan实现了Scanner接口。
func (n *NullBool) Scan(value interface{}) error

// Value implements the driver Valuer interface.

// Value实现了driver的Valuer接口。
func (n NullBool) Value() (driver.Value, error)

// NullFloat64 represents a float64 that may be null. NullFloat64 implements the
// Scanner interface so it can be used as a scan destination, similar to
// NullString.

// NullFloat64代表了可空的float64类型。
// NullFloat64实现了Scanner接口，所以它和NullString一样可以被当做scan的目标变量。
type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

// Scan implements the Scanner interface.

// Scan实现了Scanner接口。
func (n *NullFloat64) Scan(value interface{}) error

// Value implements the driver Valuer interface.

// Value实现了driver的Valuer接口。
func (n NullFloat64) Value() (driver.Value, error)

// NullInt64 represents an int64 that may be null. NullInt64 implements the Scanner
// interface so it can be used as a scan destination, similar to NullString.

// NullInt64代表了可空的int64类型。
// NullInt64实现了Scanner接口，所以它和NullString一样可以被当做scan的目标变量。
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// Scan implements the Scanner interface.

// Scan实现了Scaner接口。
func (n *NullInt64) Scan(value interface{}) error

// Value implements the driver Valuer interface.

// Value实现了driver Valuer接口。
func (n NullInt64) Value() (driver.Value, error)

// NullString represents a string that may be null. NullString implements the
// Scanner interface so it can be used as a scan destination:
//
//	var s NullString
//	err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
//	...
//	if s.Valid {
//	   // use s.String
//	} else {
//	   // NULL value
//	}

// NullString代表一个可空的string。
// NUllString实现了Scanner接口，所以它可以被当做scan的目标变量使用:
//
//	var s NullString
//	err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
//	...
//	if s.Valid {
//	   // use s.String
//	} else {
//	   // NULL value
//	}
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.

// Scan实现了Scanner接口。
func (ns *NullString) Scan(value interface{}) error

// Value implements the driver Valuer interface.

// Value实现了driver Valuer接口。
func (ns NullString) Value() (driver.Value, error)

// RawBytes is a byte slice that holds a reference to memory owned by the database
// itself. After a Scan into a RawBytes, the slice is only valid until the next
// call to Next, Scan, or Close.

// RawBytes是一个字节数组，它是由数据库自己维护的一个内存空间。
// 当一个Scan被放入到RawBytes中之后，你下次调用Next，Scan或者Close就可以获取到slice了。
type RawBytes []byte

// A Result summarizes an executed SQL command.

// 一个Result结构代表了一个执行过的SQL命令。
type Result interface {
	// LastInsertId returns the integer generated by the database
	// in response to a command. Typically this will be from an
	// "auto increment" column when inserting a new row. Not all
	// databases support this feature, and the syntax of such
	// statements varies.
	LastInsertId() (int64, error)

	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected() (int64, error)
}

// Row is the result of calling QueryRow to select a single row.

// Row是调用QueryRow的结果，代表了查询操作的一行数据。
type Row struct {
	// contains filtered or unexported fields
}

// Scan copies the columns from the matched row into the values pointed at by dest.
// If more than one row matches the query, Scan uses the first row and discards the
// rest. If no row matches the query, Scan returns ErrNoRows.

// Scan将符合的行的对应列拷贝到dest指的对应值中。
// 如果多于一个的行满足查询条件，Scan使用第一行，而忽略其他行。
// 如果没有行满足查询条件，Scan返回ErrNoRows。
func (r *Row) Scan(dest ...interface{}) error

// Rows is the result of a query. Its cursor starts before the first row of the
// result set. Use Next to advance through the rows:
//
//	rows, err := db.Query("SELECT ...")
//	...
//	defer rows.Close()
//	for rows.Next() {
//	    var id int
//	    var name string
//	    err = rows.Scan(&id, &name)
//	    ...
//	}
//	err = rows.Err() // get any error encountered during iteration
//	...

// Rows代表查询的结果。它的指针最初指向结果集的第一行数据，需要使用Next来进一步操作。
//
//	rows, err := db.Query("SELECT ...")
//	...
//	for rows.Next() {
//	    var id int
//	    var name string
//	    err = rows.Scan(&id, &name)
//	    ...
//	}
//	err = rows.Err() // get any error encountered during iteration
//	...
type Rows struct {
	// contains filtered or unexported fields
}

// Close closes the Rows, preventing further enumeration. If Next returns false,
// the Rows are closed automatically and it will suffice to check the result of
// Err. Close is idempotent and does not affect the result of Err.

// Close 关闭 Rows，阻止了进一步枚举。若 Next 返回 false，则 Rows
// 会自动关闭并能够检查 Err 的结果。Close 是幂等的，并不会影响 Err 的结果。
func (rs *Rows) Close() error

// Columns returns the column names. Columns returns an error if the rows are
// closed, or if the rows are from QueryRow and there was a deferred error.

// Columns返回列名字。 当rows设置了closed，Columns方法会返回error。
func (rs *Rows) Columns() ([]string, error)

// Err returns the error, if any, that was encountered during iteration. Err may be
// called after an explicit or implicit Close.

// Err
// 返回错误。如果有错误的话，就会在循环过程中捕获到。 Err 可能会在一个显式或隐式的 Close 后调用。
func (rs *Rows) Err() error

// Next prepares the next result row for reading with the Scan method. It returns
// true on success, or false if there is no next result row or an error happened
// while preparing it. Err should be consulted to distinguish between the two
// cases.
//
// Every call to Scan, even the first one, must be preceded by a call to Next.

// Next获取下一行的数据以便给Scan调用。
// 在成功的时候返回true，在没有下一行数据，或在准备过程中发生错误时返回false。 应通过 Err
// 来区分这两种情况。
//
// 每次调用来Scan获取数据，甚至是第一行数据，都需要调用Next来处理。
func (rs *Rows) Next() bool

// Scan copies the columns in the current row into the values pointed at by dest.
//
// If an argument has type *[]byte, Scan saves in that argument a copy of the
// corresponding data. The copy is owned by the caller and can be modified and held
// indefinitely. The copy can be avoided by using an argument of type *RawBytes
// instead; see the documentation for RawBytes for restrictions on its use.
//
// If an argument has type *interface{}, Scan copies the value provided by the
// underlying driver without conversion. If the value is of type []byte, a copy is
// made and the caller owns the result.

// Scan将当前行的列输出到dest指向的目标值中。
//
// 如果有个参数是*[]byte的类型，Scan在这个参数里面存放的是相关数据的拷贝。
// 这个拷贝是调用函数的人所拥有的，并且可以随时被修改和存取。这个拷贝能避免使用*RawBytes；
// 关于这个类型的使用限制请参考文档。
//
// 如果有个参数是*interface{}类型，Scan会将底层驱动提供的这个值不做任何转换直接拷贝返回。
// 如果值是[]byte类型，Scan就会返回一份拷贝，并且调用者获得返回结果。
func (rs *Rows) Scan(dest ...interface{}) error

// Scanner is an interface used by Scan.

// Scanner是被Scan使用的接口。
type Scanner interface {
	// Scan assigns a value from a database driver.
	//
	// The src value will be of one of the following restricted
	// set of types:
	//
	//    int64
	//    float64
	//    bool
	//    []byte
	//    string
	//    time.Time
	//    nil - for NULL values
	//
	// An error should be returned if the value can not be stored
	// without loss of information.
	Scan(src interface{}) error
}

// Stmt is a prepared statement. Stmt is safe for concurrent use by multiple
// goroutines.

// Stmt是定义好的声明。多个goroutine并发使用Stmt是安全的。
type Stmt struct {
	// contains filtered or unexported fields
}

// Close closes the statement.

// 关闭声明。
func (s *Stmt) Close() error

// Exec executes a prepared statement with the given arguments and returns a Result
// summarizing the effect of the statement.

// Exec根据给出的参数执行定义好的声明，并返回Result来显示执行的结果。
func (s *Stmt) Exec(args ...interface{}) (Result, error)

// Query executes a prepared query statement with the given arguments and returns
// the query results as a *Rows.

// Query根据传递的参数执行一个声明的查询操作，然后以*Rows的结果返回查询结果。
func (s *Stmt) Query(args ...interface{}) (*Rows, error)

// QueryRow executes a prepared query statement with the given arguments. If an
// error occurs during the execution of the statement, that error will be returned
// by a call to Scan on the returned *Row, which is always non-nil. If the query
// selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's
// Scan scans the first selected row and discards the rest.
//
// Example usage:
//
//	var name string
//	err := nameByUseridStmt.QueryRow(id).Scan(&name)

// QueryRow根据传递的参数执行一个声明的查询操作。如果在执行声明过程中发生了错误，
// 这个error就会在Scan返回的*Row的时候返回，而这个*Row永远不会是nil。
// 如果查询没有任何行数据，*Row的Scan操作就会返回ErrNoRows。
// 否则，*Rows的Scan操作就会返回第一行数据，并且忽略其他行。
//
// Example usage:
//
//	var name string
//	err := nameByUseridStmt.QueryRow(id).Scan(&name)
func (s *Stmt) QueryRow(args ...interface{}) *Row

// Tx is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
//
// After a call to Commit or Rollback, all operations on the transaction fail with
// ErrTxDone.

// Tx代表运行中的数据库事务。
//
// 必须调用Commit或者Rollback来结束事务。
//
// 在调用 Commit 或者 Rollback
// 之后，所有对事务的后续操作就会返回 ErrTxDone。
type Tx struct {
	// contains filtered or unexported fields
}

// Commit commits the transaction.

// Commit提交事务。
func (tx *Tx) Commit() error

// Exec executes a query that doesn't return rows. For example: an INSERT and
// UPDATE.

// Exec执行不返回任何行的操作。 例如：INSERT和UPDATE操作。
func (tx *Tx) Exec(query string, args ...interface{}) (Result, error)

// Prepare creates a prepared statement for use within a transaction.
//
// The returned statement operates within the transaction and can no longer be used
// once the transaction has been committed or rolled back.
//
// To use an existing prepared statement on this transaction, see Tx.Stmt.

// Prepare在一个事务中定义了一个操作的声明。
//
// 这里定义的声明操作一旦事务被调用了commited或者rollback之后就不能使用了。
//
// 关于如何使用定义好的操作声明，请参考Tx.Stmt。
func (tx *Tx) Prepare(query string) (*Stmt, error)

// Query executes a query that returns rows, typically a SELECT.

// Query执行哪些返回行的查询操作，比如SELECT。
func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error)

// QueryRow executes a query that is expected to return at most one row. QueryRow
// always return a non-nil value. Errors are deferred until Row's Scan method is
// called.

// QueryRow执行的查询至多返回一行数据。
// QueryRow总是返回非空值。只有当执行行的Scan方法的时候，才会返回Error。
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row

// Rollback aborts the transaction.

// Rollback回滚事务。
func (tx *Tx) Rollback() error

// Stmt returns a transaction-specific prepared statement from an existing
// statement.
//
// Example:
//
//	updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
//	...
//	tx, err := db.Begin()
//	...
//	res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)

// Stmt从一个已有的声明中返回指定事务的声明。
//
// 例子:
//
//	updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
//	...
//	tx, err := db.Begin()
//	...
//	res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)
func (tx *Tx) Stmt(stmt *Stmt) *Stmt
