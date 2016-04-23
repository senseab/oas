package gooas

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	OasDefineHeaderPrefix = "x-oas-"
)

func appendParam(url string, params map[string]interface{}) string {
	l := make([]string)
	for k, v := range params {
		k = strings.Replace(k, "_", "-", -1)
		if k == "maxkeys" {
			k = "max-keys"
		}
		nv := fmt.Sprint(v)
		if nv != "" {
			l = append(l, fmt.Sprintf("%s=%s",
				url.QueryEscape(k), url.QueryEscape(nv)))
		} else if k == "acl" {
			l = append(l, url.QueryEscape(k))
		} else if nv == "" {
			l = append(l, url.QueryEscape(k))
		}
	}
	if len(l) != 0 {
		url = fmt.Sprintf("%s?%s", url, strings.Join(l, "&"))
	}
	return url
}

func safeGetElement(name string, contianer map[string]interface{}) interface{} {
	for k, v := range contianer {
		if strings.ToLower(strings.TrimSpace(k)) ==
			strings.ToLower(strings.TrimSpace(name)) {
			return v
		}
	}
	return nil
}

func formatHeader(headers http.Header) http.Header {
	tmpHeaders = make(http.Header)
	for k, _ := range headers {
		if strings.HasPrefix(strings.ToLower(k), OasDefineHeaderPrefix) {
			kLower := strings.ToLower(k)
			tmpHeaders.Set(kLower, headers.Get(k))
		} else {
			tmpHeaders.Set(k) = headers.Get(k)
		}
	}
	return tmpHeaders
}

func getAssign(secret, method string, headers http.Header,
	resource string, result *[]string) string {
	canonicalizedBcHeaders := ""
	canonicalizedResource := resource

	date := headers.Get("Date")
	tmpHeader := formatHeader(headers)
	if len(*tmpHeader) > 0 {
		xHeaderList := make([]string)
		for k, _ := range *tmpHeader {
			xHeaderList = append(xHeaderList, k)
		}
		sort.Strings(xHeaderList)
		for _, k := range xHeaderList {
			if strings.HasPrefix(k, OasDefineHeaderPrefix) {
				canonicalizedBcHeaders = fmt.Sprintf("%s%s:%v\n",
					canonicalizedBcHeaders, k, tmpHeader.Get(k))
			}
		}
	}
	stringToSign := fmt.Sprintf("%s\n%s\n%s%s", method, date,
		canonicalizedBcHeaders, canonicalizedResource)
	*result = append(*result, stringToSign)

	h := sha1.New()
	h.Write([]byte(stringToSign))
	h := hmac.New(h, []byte(secret))
	b := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return strings.TrimSpace(b)
}
