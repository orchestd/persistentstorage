package persistentstorage

import (
	"context"
	"time"
)

type PersistentStorage interface {
	QueryOne(c context.Context, target QueryGetter, params map[string]interface{}) error
	QueryMany(c context.Context, target QueryGetter, params map[string]interface{}) error

	QueryInt(c context.Context, query QueryGetter, params map[string]interface{}) (int64, error)
	QueryString(c context.Context, query QueryGetter, params map[string]interface{}) (string, error)

	GetOne(c context.Context, target interface{}, params map[string]interface{}) error
	GetMany(c context.Context, target interface{}, params map[string]interface{}) error

	Insert(c context.Context, target BaseModelSetter, now time.Time) error
	Update(c context.Context, model interface{}, update map[string]interface{}, params map[string]interface{}) error
	Delete(c context.Context, model interface{}, params map[string]interface{}, now time.Time) error

	Exec(c context.Context, target QueryGetter, params map[string]interface{}) error
}

type QueryGetter interface {
	GetQuery() string
}

type UpdateStampGetter interface {
	GetUpdateStamp(c context.Context, description string) (int64, error)
}

type BaseModelSetter interface {
	SetCreateDate(date time.Time)
	SetUpdateStamp(updateStamp int64)
	SetCancelDate(date time.Time)
}
