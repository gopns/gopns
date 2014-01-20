package access

import (
	"github.com/gopns/gopns/model"
)

type AppManager interface {
	ListApps(cursor string) (apps *[]model.App, newCursor string, err error)
	GetApp(id string) (model.App, error)
	CreateOrUpdateApp(app model.App) error
	GetPlatformApps(appId string, cursor string) (papps *[]model.PlatformApp, newCursor string, err error)
	CreateOrUpdatePlatformApp(papp model.PlatformApp) error
	// whenever platform app is created it will check if the app already exists with SNS otherwise will create it
}
