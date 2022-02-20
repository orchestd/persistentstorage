package defaultPersistentStorage

import (
	"bitbucket.org/HeilaSystems/persistentstorage/mysqlPersistentStorage"
	"bitbucket.org/HeilaSystems/persistentstorage/updateStampService"
	"go.uber.org/fx"
)

func GetDefaultPersistentStorage() fx.Option {
	return fx.Options(
		fx.Provide(updateStampService.NewUpdateStampService),
		fx.Provide(mysqlPersistentStorage.NewMySQLDb),
	)
}
