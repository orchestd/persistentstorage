package defaultPersistentStorage

import (
	"github.com/orchestd/persistentstorage/mysqlPersistentStorage"
	"github.com/orchestd/persistentstorage/updateStampService"
	"go.uber.org/fx"
)

func GetDefaultPersistentStorage() fx.Option {
	return fx.Options(
		fx.Provide(updateStampService.NewUpdateStampService),
		fx.Provide(mysqlPersistentStorage.NewMySQLDb),
	)
}
