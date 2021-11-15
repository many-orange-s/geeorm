package orm

import (
	"database/sql"
	"geeorm/dialect"
	"geeorm/log"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(diver string, source string) (e *Engine) {
	db, err := sql.Open(diver, source)
	if err != nil {
		log.Errorf("connect err %s", err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Errorf("connect err %s", err)
		return
	}

	dial, ok := dialect.GetDialect(diver)
	if !ok {
		log.Errorf("dialect %s Not Found", diver)
		return
	}

	e = &Engine{
		db:      db,
		dialect: dial,
	}

	log.Infof("connect success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Errorf("close err %s", err)
	}
}

func (e *Engine) NewSession() *Session {
	return New(e.db, e.dialect)
}
