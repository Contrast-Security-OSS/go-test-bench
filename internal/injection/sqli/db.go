package sqli

import (
	"database/sql"
	"database/sql/driver"
)

type stmt struct{}

func (stmt) Close() error                               { return nil }
func (stmt) NumInput() int                              { return 0 }
func (stmt) Exec([]driver.Value) (driver.Result, error) { return nil, nil }
func (stmt) Query([]driver.Value) (driver.Rows, error)  { return nil, nil }

type conn struct{}

func (conn) Prepare(query string) (driver.Stmt, error) { return stmt{}, nil }
func (conn) Close() error                              { return nil }
func (conn) Begin() (driver.Tx, error)                 { return nil, nil }

type testdriver struct{}

func (testdriver) Open(name string) (driver.Conn, error) { return conn{}, nil }

func init() {
	sql.Register("test", testdriver{})
}
