package session

import "github.com/go-examples-with-tests/database/v3/log"

func (s *Session) Begin() (err error) {
	log.Info("transactioin begin")
	if s.transaction, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.transaction.Commit(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.transaction.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
