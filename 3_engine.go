package geeorm

import "database/sql"

// Engine 是orm与用户交互的入口。
//
// Session负责与数据库的交互，那交互前的准备工作（比如连接/测试数据库），交互后的收尾工作（关闭连接）等就交给 Engine 来负责了。。
type Engine struct {
	db *sql.DB
}

// Open 打开一个数据库，返回一个数据库连接对象.
func Open(driverName,dataSource string) (e *Engine,err error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		Error(err)
		return
	}

	if err = db.Ping(); err != nil {
			Error(err)
			return
	}
	//e.db = db 这样调用是错误的.
	e = &Engine{db: db}
	Info("连接数据库成功了")
	return
}

func (e *Engine) Close(){
	if err := e.db.Close(); err != nil {
		Error("关闭数据库连接失败")
	}
	Info("关闭数据库连接成功了")
}

// Session 是Engine 与 Session 建立联系的方法，使Engine能间接的操作数据库。
//
// 通过 Engine 实例创建Session，进而能与数据库进行交互了。
func (e *Engine)Session()*Session {
	return New(e.db)
}

