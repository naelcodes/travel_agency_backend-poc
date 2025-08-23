package managers

import (
	"xorm.io/xorm"
)

type TransactionManager struct {
	session *xorm.Session
}

func NewTransactionManager(engine *xorm.Engine) *TransactionManager {
	return &TransactionManager{
		session: engine.NewSession(),
	}
}

func (tm *TransactionManager) Begin() error {
	return tm.session.Begin()
}

func (tm *TransactionManager) Commit() {
	_ = tm.session.Commit()
	tm.session.Close()
}

func (tm *TransactionManager) Rollback() error {
	if tm.session == nil {
		return nil
	}
	return tm.session.Rollback()
}

func (tm *TransactionManager) GetTransaction() *xorm.Session {
	return tm.session
}

// func (tm *TransactionManager) CatchError() error {
// 	utils.Logger.Info(fmt.Sprintf("[TransactionManager - CatchError] Transaction: %v", tm.session))
// 	if err := recover(); err != nil {
// 		rollBackErr := tm.Rollback()
// 		err, _ := err.(error)
// 		if rollBackErr != nil {
// 			return CustomErrors.UnknownError(errors.Join(rollBackErr, err))
// 		}
// 		return err
// 	}
// 	return nil
// }
