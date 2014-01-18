package dynamodb

import (
	"encoding/json"
	"errors"
	"github.com/gopns/gopns/aws"
	"io"
	"log"
	"strings"
	"testing"
)

func TestClientCreation(t *testing.T) {
	dynamoClient, err := New("Derp", "Derp is actually smart", "us-east-1")
	if err != nil || dynamoClient == nil {
		t.Error("Unable to create new client")
	}
}
func TestFindTableFound(t *testing.T) {

	dynamoClient, err := _new("Derp", "Derp is actually smart", "us-east-1", (&MockRequestor{t, 200, []string{"DerpTable"}, false, false}).mockMakeRequest)
	if err != nil {
		t.Errorf("Unable to create client %s", err.Error())
	}

	found, err := dynamoClient.FindTable("DerpTable")

	if err != nil {
		t.Errorf("Should have found table but did but got error %s", err.Error())
	} else if !found {
		t.Errorf("Should have found table but did not")
	}
}

func TestFindTableNoTable(t *testing.T) {
	dynamoClient, err := _new("Derp", "Derp is actually smart", "us-east-1", (&MockRequestor{t, 200, make([]string, 0, 0), false, false}).mockMakeRequest)
	if err != nil {
		t.Errorf("Unable to create client %s", err.Error())
	}

	found, err := dynamoClient.FindTable("DerpTable")

	if err != nil {
		t.Errorf("Should have found table but did but got error %s", err.Error())
	} else if found {
		t.Errorf("Should not have found table but did")
	}
}

func TestFindTableNotFound(t *testing.T) {
	dynamoClient, err := _new("Derp", "Derp is actually smart", "us-east-1", (&MockRequestor{t, 200, []string{"NoMatch"}, false, false}).mockMakeRequest)
	if err != nil {
		t.Errorf("Unable to create client %s", err.Error())
	}

	found, err := dynamoClient.FindTable("DerpTable")

	if err != nil {
		t.Errorf("Should have found table but did but got error %s", err.Error())
	} else if found {
		t.Errorf("Should not have found table but did")
	}
}

func TestFindTableHttpError(t *testing.T) {
	dynamoClient, err := _new("Derp", "Derp is actually smart", "us-east-1", (&MockRequestor{t, 200, []string{"NoMatch"}, true, false}).mockMakeRequest)
	if err != nil {
		t.Errorf("Unable to create client %s", err.Error())
	}

	_, err = dynamoClient.FindTable("DerpTable")

	if err == nil {
		t.Errorf("Simulated Http Error but got %s", err.Error())
	} else if !strings.Contains(err.Error(), "Simulated Http Error") {
		t.Errorf("Simulated Http Error but got %s", err.Error())
	}
}

func TestFindTableOKParsingError(t *testing.T) {
	dynamoClient, err := _new("Derp", "Derp is actually smart", "us-east-1", (&MockRequestor{t, 200, []string{"NoMatch"}, false, true}).mockMakeRequest)
	if err != nil {
		t.Errorf("Unable to create client %s", err.Error())
	}

	_, err = dynamoClient.FindTable("DerpTable")

	if err == nil {
		t.Errorf("Expected Parsing Error but got nil")
	} else if !strings.Contains(err.Error(), "Unable to serialize") {
		t.Errorf("Expected Parsing Error but got %s", err.Error())
	}
}

func TestFindTable500ParsingError(t *testing.T) {
	dynamoClient, err := _new("Derp", "Derp is actually smart", "us-east-1", (&MockRequestor{t, 500, nil, false, true}).mockMakeRequest)
	if err != nil {
		t.Errorf("Unable to create client %s", err.Error())
	}

	_, err = dynamoClient.FindTable("DerpTable")

	if err == nil {
		t.Errorf("Expected Parsing Error but got nil")
	} else if !strings.Contains(err.Error(), "Unable to serialize") {
		t.Errorf("Expected Parsing Error but got %s", err.Error())
	}
}

func TestFindTable500(t *testing.T) {
	dynamoClient, err := _new("Derp", "Derp is actually smart", "us-east-1", (&MockRequestor{t, 500, nil, false, false}).mockMakeRequest)
	if err != nil {
		t.Errorf("Unable to create client %s", err.Error())
	}

	_, err = dynamoClient.FindTable("DerpTable")

	if err == nil {
		t.Errorf("Expecting Unable to find table. TestType: TestMessage but got nil")
	} else if !strings.Contains(err.Error(), "Unable to find table. TestType: TestMessage") {
		t.Errorf("Expecting Unable to find table. TestType: TestMessage but found %s", err.Error())
	}
}

type MockRequestor struct {
	t            *testing.T
	status       int
	tables       []string
	httpError    bool
	parsingError bool
}

func (this *MockRequestor) mockMakeRequest(host string, query string, action string) (int, io.ReadCloser, error) {
	log.Printf("Mock Make Request Called")
	if this.httpError {
		return 0, nil, errors.New("Simulated Http Error")
	} else {
		return this.status, this, nil
	}
}

func (this *MockRequestor) Read(p []byte) (n int, err error) {
	tables := make(map[string][]string)
	tables["TableNames"] = this.tables

	if !this.parsingError && this.status == 200 {
		resp, _ := json.Marshal(tables)
		copy(p, resp)
		return len(resp), io.EOF
	} else if !this.parsingError && this.status != 200 {
		resp, _ := json.Marshal(aws.ErrorStruct{"TestType", "TestCode", "TestMessage"})
		copy(p, resp)
		return len(resp), io.EOF
	} else {
		resp := "InvalidParsingText"
		copy(p, resp)
		return len(resp), io.EOF
	}
}

func (this *MockRequestor) Close() error {
	return nil
}
