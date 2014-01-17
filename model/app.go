package model

import (
	"errors"
	"regexp"
)

type Platform string

const (
	GCM          Platform = "GCM"
	APNS         Platform = "APNS"
	APNS_SANDBOX Platform = "APNS_SANDBOX"
	ADM          Platform = "ADM"
)

type App struct {
	id string
}

func NewApp(id string) (app *App) {
	return &App{id}
}

func (app App) Id() string{
	return app.id
}

func (app *App) SetId(id string) {
	app.id = id
}

func NewPlatformApp(appId string, platform Platform, arn string) (papp *PlatformApp, error) {
	papp = &PlatformApp{appId: appId, arn: arn}
	papp.SetPlatform(platform)
	return papp
}

type PlatformApp struct {
	appId    string
	platform Platform
	arn      string
}

func (papp PlatformApp) AppId() string {
	return papp.appId
}

func (papp PlatformApp) Platform() Platform {
	return papp.platform
}

func (papp PlatformApp) Arn() string {
	return papp.arn
}

func (papp *PlatformApp) SetAppId(appId string) {
	papp.appId = appId
}

func (papp *PlatformApp) SetPlatform(platform Platform) error {
	if err := ValidatePlatform(platform); err != nil {
		return err
	}else{
		app.platform = platform
	}
}

func (papp *PlatformApp) SetArn(arn String) {
	papp.arn = arn
}

func (app *App) SetId(id string) {
	app.id = id
}

var platformRegex = regexp.MustCompile("^(" + GCM + "|" + ADM + "|" APNS + "|" + APNS_SANDBOX ")$")


func ValidatePlatform(p Platform) error {
	if !platformRegex.MatchString(p) {
		return errors.New("Invalid platform, valid values: " + GCM + "," + APNS + "," + APNS_SANDBOX + "or " + ADM)
	}
	return nil
}
