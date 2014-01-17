package gopnsconfig

import (
	"testing"
)

func TestPlatformAppParse(t *testing.T) {
	confFile := NewConfigFile()
	confFile.AddOption("default", "platform-applications", "app1,app2")
	confFile.AddSection("app1")
	confFile.AddSection("app2")

	confFile.AddOption("app1", "arn", "SomeArn1")
	confFile.AddOption("app2", "arn", "SomeArn2")

	confFile.AddOption("app1", "type", "SomeType1")
	confFile.AddOption("app2", "type", "SomeType2")

	platformApps := parsePlatformAppConfig(confFile)
	if platformApps["app1"].Arn() != "SomeArn1" {
		t.Errorf("Expected Arn: %s, found %s", "SomeArn1", platformApps["app1"].Arn())
	}

	if platformApps["app1"].Type() != "SomeType1" {
		t.Errorf("Expected Arn: %s, found %s", "SomeArn1", platformApps["app1"].Arn())
	}

	if platformApps["app2"].Arn() != "SomeArn2" {
		t.Errorf("Expected Arn: %s, found %s", "SomeArn1", platformApps["app2"].Arn())
	}

	if platformApps["app2"].Type() != "SomeType2" {
		t.Errorf("Expected Arn: %s, found %s", "SomeArn1", platformApps["app2"].Arn())
	}

	if platformApps["app1"].ConfigMap()["Arn"] != "SomeArn1" {
		t.Errorf("Expected Arn: %s, found %s", "SomeArn1", platformApps["app1"].Arn())
	}
}

func TestPlatformAppParseFail(t *testing.T) {

	defer expectPanic(t)
	confFile := NewConfigFile()
	parsePlatformAppConfig(confFile)
}

func TestPlatformAppParseFailApp(t *testing.T) {

	defer expectPanic(t)
	confFile := NewConfigFile()
	confFile.AddOption("default", "platform-applications", "app1,app2")
	confFile.AddSection("app1")
	confFile.AddSection("app2")

	confFile.AddOption("app1", "arns", "SomeArn1")
	confFile.AddOption("app2", "arn", "SomeArn2")

	confFile.AddOption("app1", "type", "SomeType1")
	confFile.AddOption("app2", "type", "SomeType2")

	parsePlatformAppConfig(confFile)
}

func TestPlatformAppParseFailAppData(t *testing.T) {

	defer expectPanic(t)
	confFile := NewConfigFile()
	confFile.AddOption("default", "platform-applications", "app1,app2")
	parsePlatformAppConfig(confFile)
}

func expectPanic(t *testing.T) {
	if recover() == nil {
		t.Error("Panic Expected but not found")
	}
}
