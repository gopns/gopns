package aws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/usmanismail/gpns/com/techtraits/gpns/gpnsconfig"

	"net/url"
	"sort"
	"strings"
)

var b64 = base64.StdEncoding

func SignRequest(awsConfig gpnsconfig.AWSConfig, method, path string, params url.Values, host string) {
	params.Set("AWSAccessKeyId", awsConfig.UserID())
	params.Set("SignatureVersion", "2")
	params.Set("SignatureMethod", "HmacSHA256")
	params.Set("Version", "2010-03-31")

	paramsStr := params.Encode()
	var sarray []string
	sarray = append(sarray, strings.Replace(paramsStr, "+", "%20", -1))
	sort.StringSlice(sarray).Sort()
	joined := strings.Join(sarray, "&")
	payload := method + "\n" + host + "\n" + path + "\n" + joined

	hash := hmac.New(sha256.New, []byte(awsConfig.UserSecret()))

	hash.Write([]byte(payload))
	signature := make([]byte, b64.EncodedLen(hash.Size()))
	b64.Encode(signature, hash.Sum(nil))

	params.Set("Signature", string(signature))

}
