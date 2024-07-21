package database

// Transaction 事物操作 跟数据库有关，注入数据库的连接
type Transaction interface {
	Action(func(conn DBConn) error) error
}

// DBConn 事物的提交和回滚
type DBConn interface {
	Commit()
	Rollback()
	Begin()
}
