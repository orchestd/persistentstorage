package persistentstorage

import (
	"context"
)

type PersistentStorage interface {
	QueryOne(c context.Context, target QueryGetter, params map[string]interface{}) error
	QueryMany(c context.Context, target QueryGetter, params map[string]interface{}) error

	QueryInt(c context.Context, query QueryGetter, params map[string]interface{}) (int64, error)
	QueryString(c context.Context, query QueryGetter, params map[string]interface{}) (string, error)

	GetOne(c context.Context, target interface{}, params interface{}) error
	GetMany(c context.Context, target interface{}, params interface{}) error

	Insert(c context.Context, target interface{}) error
	Update(c context.Context, update interface{}, query interface{}) error
	Delete(c context.Context, model interface{}, params interface{}) error

	Exec(c context.Context, target QueryGetter, params map[string]interface{}) error
}

type QueryGetter interface {
	GetQuery() string
}

type UpdateStampGetter interface {
	GetUpdateStamp(c context.Context, description string) (int64, error)
}
