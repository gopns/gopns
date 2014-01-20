package access

import (
	"github.com/gopns/gopns/model"
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sns"
)

type AppManager interface {
	ListApps(cursor string) (apps *[]model.App, newCursor string, err error)
	GetApp(id string) (*model.App, error)
	CreateOrUpdateApp(app model.App) error
	GetPlatformApps(appId string, cursor string) (papps *[]model.PlatformApp, newCursor string, err error)
	CreateOrUpdatePlatformApp(papp model.PlatformApp) error
	// whenever platform app is created it will check if the app already exists with SNS otherwise will create it
}


type DefaultAppManager struct {
	SnsClient    sns.SNSClient
	DynamoClient dynamodb.DynamoClient
	AppTable  string
	PlatformAppTable string
}

func NewAppManager(snsClient sns.SNSClient, dynamoClient dynamodb.DynamoClient, appTable string, pappTable string) AppManager {
	appManagerInstance := &DefaultAppManager{
		snsClient,
		dynamoClient,
		appTable,
		pappTable}
	return appManagerInstance
}

func (dm *DefaultAppManager) ListApps(cursor string) (apps *[]model.App, newCursor string, err error) {
	return nil, "", nil
}

func (dm *DefaultAppManager) GetApp(id string) (*model.App, error) {
	return nil, nil
}

func (dm *DefaultAppManager) CreateOrUpdateApp(app model.App) error {
	return nil
}

func (dm *DefaultAppManager) GetPlatformApps(appId string, cursor string) (papps *[]model.PlatformApp, newCursor string, err error) {
	return nil, "", nil
}

func (dm *DefaultAppManager) CreateOrUpdatePlatformApp(papp model.PlatformApp) error {
	return nil
}