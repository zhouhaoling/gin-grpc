package repository

import "test.com/project-user/internal/repository/database"

type TransactionImpl struct {
	conn database.DBConn
}

func NewTransaction() *TransactionImpl {
	return &TransactionImpl{
		conn: database.NewTran(),
	}
}

func (t TransactionImpl) Action(f func(conn database.DBConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}
