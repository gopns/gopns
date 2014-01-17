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

func (app App) Id() string {
	return app.id
}

func (app *App) SetId(id string) {
	app.id = id
}

func NewPlatformApp(appId string, platform Platform, arn string) (*PlatformApp, error) {
	papp := &PlatformApp{appId: appId, arn: arn}
	if err := papp.SetPlatform(platform); err != nil {
		return nil, err
	}
	return papp, nil
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

func (papp *PlatformApp) SetPlatform(p Platform) error {
	if err := ValidatePlatform(p); err != nil {
		return err
	}

	papp.platform = p
	return nil
}

func (papp *PlatformApp) SetArn(arn string) {
	papp.arn = arn
}

var platformRegex = regexp.MustCompile("^(" + string(GCM) + "|" + string(ADM) + "|" + string(APNS) + "|" + string(APNS_SANDBOX) + ")$")

func ValidatePlatform(p Platform) error {
	if !platformRegex.MatchString(string(p)) {
		return errors.New("Invalid platform, valid values: " + string(GCM) + "," + string(APNS) + "," + string(APNS_SANDBOX) + " or " + string(ADM))
	}
	return nil
}
