package infrastructure

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"test_aws/interfaces"
)

type SqliteHandler struct {
	Conn *sql.DB
}

type SqliteRow struct {
	Rows *sql.Rows
}

func (handler *SqliteHandler) Execute(statement string, args ...interface{}) error {
	stmt, err := handler.Conn.Prepare(statement)
	_, err = stmt.Exec(args...)

	return err
}

func (handler *SqliteHandler) Query(statement string) interfaces.Row {
	rows, err := handler.Conn.Query(statement)

	if err != nil {
		fmt.Println(err)
		return new(SqliteRow)
	}

	row := new(SqliteRow)
	row.Rows = rows
	return row
}

func (r SqliteRow) Scan(dest ...interface{}) error {
	error := r.Rows.Scan(dest...)
	return error
}

func (r SqliteRow) Next() bool {
	return r.Rows.Next()
}

func (r SqliteRow) Close() {
	r.Rows.Close()
}

func NewSqliteHandler(dbfileName string) *SqliteHandler {
	conn, _ := sql.Open("sqlite3", dbfileName)
	sqliteHandler := new(SqliteHandler)
	sqliteHandler.Conn = conn
	return sqliteHandler
}

