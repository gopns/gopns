package sns

import (
	"github.com/usmanismail/gpns/com/techtraits/gpns/aws"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RegistrarStruct struct {
	awsConfig gpnsconfig.AWSConfig
}

func Initilize(awsConfig gpnsconfig.AWSConfig) Registrar {
	return RegistrarStruct{awsConfig}
}

type Registrar interface {
	RegisterDevice() (arn string, err error)
}

func (this RegistrarStruct) RegisterDevice() (arn string, err error) {
	values := url.Values{}
	values.Set("Action", "CreatePlatformEndpoint")
	values.Set("CustomUserData", "ENG_GB")
	values.Set("Token", "APA91bF1felnMnAgGtJm4NWcp2Zv4zpeKDko742sSdhBfK9uFtYREcoFQnLBuGockhSxMHMqTf2t5y_HwYe32PYVJNg0rwvGpdMbJwedgOZVdQ2lcQl6yB6CCp1xw2SosQcxU5JvGJLiO3aPuh53Qu_3Gzz-zpUgja2ZgLe31TAtHpY3Kgo3Fmc")
	values.Set("PlatformApplicationArn", "arn:aws:sns:us-east-1:238699486533:app/GCM/Test")
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	url_, err := url.Parse("http://sns.us-east-1.amazonaws.com/")
	if err != nil {
		return "", err
	}

	aws.SignRequest(this.awsConfig, "POST", "/", values, url_.Host)

	r, err := http.PostForm(url_.String(), values)
	if err != nil {
		log.Println(err.Error())
		return "", err
	} else {
		log.Println("No error")
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {

		contents, _ := ioutil.ReadAll(r.Body)
		log.Printf("Not 200 %d %s", r.StatusCode, contents)
		return "", nil
	} else {
		contents, _ := ioutil.ReadAll(r.Body)
		log.Printf("200 %d %s", r.StatusCode, contents)
		//err = xml.Unmarshal(r.Body, resp)
	}
	return "", err
}

func multiMap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}
