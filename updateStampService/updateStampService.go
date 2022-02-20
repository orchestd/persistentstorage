package updateStampService

import (
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/configuration"
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/transport"
	"bitbucket.org/HeilaSystems/persistentstorage"
	"context"
	"fmt"
)

func NewUpdateStampService(client transport.HttpClient, config configuration.Config) persistentstorage.UpdateStampGetter {
	return updateStampService{client: client, config: config}
}

type updateStampService struct {
	client transport.HttpClient
	config configuration.Config
}

func (us updateStampService) GetUpdateStamp(c context.Context, description string) (int64, error) {
	updateStampResult := make(map[string]int64)
	payload := map[string]interface{}{"Description": description, "Increment": 1}
	err := us.client.Post(c, payload, "updateStampBase", "getUpdateStamp", &updateStampResult, nil)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}
	return updateStampResult["updateStamp"], nil
}
