package repository

import (
	"errors"

	"test.com/common/errs"
	"test.com/project-project/internal/repository/database"
)

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
	var bErr *errs.BError
	if errors.Is(err, bErr) {
		bErr = err.(*errs.BError)
		if bErr != nil {
			t.conn.Rollback()
			return bErr
		} else {
			t.conn.Commit()
			return nil
		}
	}
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}
