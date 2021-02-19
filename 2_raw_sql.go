package geeorm

import (
	"database/sql"

	"strings"
)

// Session 持有*sql.DB指针，并提供操作各种数据库的功能。
type Session struct {
	db      *sql.DB         // sql.Open() 方法连接数据库成功之后返回的指针,该指针具备操作数据库的能力。
	sql     strings.Builder // 用来拼接 SQL 语句
	sqlVars []interface{}   // 用来存储 SQL语句中占位符的对应值
}

// New 返回一个Session实例。
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// Clear 初始化Session，清空Session中的SQL语句和它对应的SQL参数值。
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// DB 返回 *sql.DB
func (s *Session) DB() *sql.DB {
	return s.db
}

// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		Error(err)
	}
	return
}

// QueryRow 从数据库获取一条记录。
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// Query 从数据库获取记录列表
func (s *Session) Query() (rows *sql.Rows, err error) {
	defer s.Clear()
	Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		Error(err)
	}
	return
}

// Append 分别向 s.Sql 和 s.sqlVars 中增加一条SQL语句。
func (s *Session) Append(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}
