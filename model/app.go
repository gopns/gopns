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
	UNKNOWN      Platform = "UNKNOWN"
)

type App struct {
	id          string
	description string
	createdAt   string
	updatedAt   string
}

func NewApp(id string, desc string) (app *App) {
	return &App{id: id, description: desc}
}

func (app App) Id() string {
	return app.id
}

func (app App) Description() string {
	return app.description
}

func (app *App) SetId(id string) {
	app.id = id
}

func (app *App) SetDescription(desc string) {
	app.description = desc
}

func NewPlatformApp(id string, appId string, platform Platform, arn string) (*PlatformApp, error) {
	papp := &PlatformApp{id: id, appId: appId, arn: arn}
	if err := papp.SetPlatform(platform); err != nil {
		return nil, err
	}
	return papp, nil
}

type PlatformApp struct {
	id              string
	name            string
	appId           string
	platform        Platform
	arn             string
	admClientId     string
	admClientSecret string
	apnsCertificate string
	apnsPrivateKey  string
	gcmApiKey       string
}

func (papp PlatformApp) Id() string {
	return papp.id
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


func (papp PlatformApp) AdmClientId() string {
	return papp.admClientId
}

func (papp PlatformApp) AdmClientSecret() string {
	return papp.admClientSecret
}

func (papp PlatformApp) ApnsCertificate() string {
	return papp.apnsCertificate
}

func (papp PlatformApp) ApnsPrivateKey() string {
	return papp.apnsPrivateKey
}

func (papp PlatformApp) GcmApiKey() string {
	return papp.gcmApiKey
}

func (papp *PlatformApp) SetId(id string) {
	papp.id = id
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

func (papp *PlatformApp) SetAdmClientId(admCid string) {
	papp.admClientId = admCid
}

func (papp *PlatformApp) SetAdmClientSecret(admSecret string) {
	papp.admClientSecret = admSecret
}

func (papp *PlatformApp) SetApnsCertificate(apnsCertificate string) {
	papp.apnsCertificate = apnsCertificate
}

func (papp *PlatformApp) SetApnsPrivateKey(apnsPrivateKey string) {
	papp.apnsPrivateKey = apnsPrivateKey
}

func (papp *PlatformApp) SetGcmApiKey(gcmApiKey string) {
	papp.gcmApiKey = gcmApiKey
}

var platformRegex = regexp.MustCompile("^(" + string(GCM) + "|" + string(ADM) + "|" + string(APNS) + "|" + string(APNS_SANDBOX) + ")$")

func ValidatePlatform(p Platform) error {
	if !platformRegex.MatchString(string(p)) {
		return errors.New("Invalid platform, valid values: " + string(GCM) + "," + string(APNS) + "," + string(APNS_SANDBOX) + " or " + string(ADM))
	}
	return nil
}
