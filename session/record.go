package session

import (
	"github.com/astaxie/beego/context"
)

type Record struct {
}

func (r *Record) Match(ctx context.Context) error {
	return nil
}

func (r *Record) Expired() bool {
	return false
}

func NewRecordFromContext(ctx context.Context) Record {
	var r Record
	return r
}
