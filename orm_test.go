package geeorm

import (
	"database/sql"
	"github.com/kotlin2018/geeorm/dialect"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

func TestSetLevel(t *testing.T) {
	SetLevel(ErrorLevel)
	if infoLog.Writer() == os.Stdout || errorLog.Writer() != os.Stdout {
		t.Fatal("failed to set log level")
	}
	SetLevel(Disabled)
	if infoLog.Writer() == os.Stdout || errorLog.Writer() == os.Stdout {
		t.Fatal("failed to set log level")
	}
}

var TestDB *sql.DB

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "../gee.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Append("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Append("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Append("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRow(t *testing.T) {
	s := NewSession()
	_, _ = s.Append("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Append("CREATE TABLE User(Name text);").Exec()
	row := s.Append("SELECT count(*) FROM User").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := Open("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func TestEngine_Session(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}

func TestSchema_RecordValues(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	values := schema.RecordValues(&User{"Tom", 18})

	name := values[0].(string)
	age := values[1].(int)

	if name != "Tom" || age != 18 {
		t.Fatal("failed to get values")
	}
}

type UserTest struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

func (u *UserTest) TableName() string {
	return "ns_user_test"
}

func TestSchema_TableName(t *testing.T) {
	schema := Parse(&UserTest{}, TestDial)
	if schema.Name != "ns_user_test" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
}






