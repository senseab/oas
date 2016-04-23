// api.go
package gooas

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	OasDefaultContentType    = "application/octet-stream"
	OasDefaultSendBufferSize = 8192
	OasDefaultGetBufferSize  = 10 * 1024 * 1024
	OasDefaultProvider       = "OAS"

	OasHttpPort  = 80
	OasHttpsPort = 443

	OasUseHttps = true
	OasNoHttps  = false

	OasUserAgent = "gooas-OAS Go SDK"
)

type OasClient struct {
	host      string
	apiKey    string
	apiSecret string
}

func NewOasClient(host, apikey, secret string, port int, security bool) *OasClient {
	o := new(OasClient)
	o.apiKey = apikey
	o.apiSecret = secret
	if security || port == 443 {
		o.host = fmt.Sprintf("https://%s:%s", host, port)
	} else {
		o.host = fmt.Sprintf("http://%s:%s", host, port)
	}
	return o
}

func (o *OasClient) getResource(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	tmpHeaders := make(map[string]interface{})
	for k, v := range params {
		tmpK := strings.TrimSpace(strings.ToLower(k))
		tmpHeaders[tmpK] = v
	}

	overrideResponseList := []string{
		"limit", "marker", "response-content-type", "response-content-language",
		"response-cache-control", "logging", "response-content-encoding",
		"acl", "uploadId", "uploads", "partNumber", "group",
		"delete", "website", "location", "objectInfo",
		"response-expires", "response-content-disposition"}
	sort.Strings(overrideResponseList)

	resource := ""
	separator := "?"
	for _, i := range overrideResponseList {
		if _, ok := tmpHeaders[strings.ToLower(i)]; ok {
			resource = fmt.Sprintf("%s%s%s", resource, separator, i)
			tmpKey := tmpHeaders[strings.ToLower(i)]
			if tmpKey != "" {
				resource = fmt.Sprintf("%s=%v", resource, tmpKey)
			}
			separator = "&"
		}
	}
	return resource
}
func (o *OasClient) httpRequest(method string, headers http.Header,
	body string, params [string]interface{}) (http.Response, error) {
	headers.Set("User-Agent", OasUserAgent)
	headers.Set("Host", o.host0)
	headers.Set("Date", time.Now().UTC().Format(time.RFC1123))
	headers.Set("x-oas-version", "0.2.5")
	if len(body) > 0 {
		headers.Set("Content-Length", fmt.Sprint(len(body)))
	}
	resource := fmt.Sprinf("%s%s", o.host, o.getResource(params))
	if len(params) != 0 {
		url = appendParam(url, params)
	}
	headers.Set("Authorization", o.createSignForNormalAuth(method, headers, resource))
	// TODO
}

func (o *OasClient) createSignForNormalAuth(method string, headers http.Header,
	resource string) string {
	res := make([]string)
	authValue := fmt.Printf("%s %s:%s", OasDefaultProvider, o.apiKey,
		getAssign(o.apiSecret, method, headers, resource, &res))
	return authValue
}

func (o *OasClient) CreateVault(name string) {

}
