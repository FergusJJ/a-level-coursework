package utils

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	http "github.com/useflyent/fhttp"
)

//Returns true or false based on the status code of the response
func CheckResponseStatus(statusCode int) (ok bool) {
	if statusCode >= 200 && statusCode <= 299 {
		ok = true
	} else if statusCode >= 300 && statusCode <= 399 {
		fmt.Printf("Got Redirect: %d\n", statusCode)
		//will want to do something else here, not sure yet though
		//return true after something else happens
		ok = false
	} else {
		ok = false
	}
	return ok
}

func FormatProxy(proxy string) string {
	proxySlice := strings.Split(proxy, ":")
	if len(proxySlice) == 4 {
		proxyUrlStr := fmt.Sprintf("http://%s:%s@%s:%s/", proxySlice[2], proxySlice[3], proxySlice[0], proxySlice[1])
		return proxyUrlStr
	}
	if proxySlice[0] == "http" || proxySlice[0] == "https" {
		//probably already in correct format
		return proxy
	}
	return ""
}

func returnOrigin(requestURL string) (originLinkHeader string) {
	u, err := url.Parse(requestURL)
	if err != nil {
		log.Println(err)
		return requestURL
	}
	originLinkHeader = fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
	return originLinkHeader
}

func getContentType(requestURL string) (contentTypeHeader string) {
	if strings.Contains(requestURL, "identity/login") {
		contentTypeHeader = "application/x-www-form-urlencoded"
		return contentTypeHeader
	}
	if requestURL == "https://my.asos.com/my-account" {
		contentTypeHeader = "application/x-www-form-urlencoded"
		return contentTypeHeader
	}
	fmt.Printf("request url not found in if else / %s\n", requestURL)
	return ""
}

//Will assign correct request headers and header order
//Different browsers will have different headers and order, if one is not recognised then I will log it to console
//After, the default chrome header order/headers will be used
func AssignGetRequestHeaders(Browser, UserAgent string) http.Header {
	//probably a switch statement here, may need extra info such as origin, but may also be better to set those in
	for {
		switch Browser {
		case "CHROME":
			return http.Header{
				"upgrade-insecure-requests": {"1"},
				"user-agent":                {UserAgent},
				"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"sec-fetch-site":            {"none"},
				"sec-fetch-mode":            {"navigate"},
				"sec-fetch-user":            {"?1"},
				"sec-fetch-dest":            {"document"},
				"sec-ch-ua":                 {`" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`},
				"sec-ch-ua-mobile":          {"?0"},
				"sec-ch-ua-platform":        {"\"Windows\""},
				"accept-encoding":           {"gzip, deflate, br"},
				"accept-language":           {"en-GB,en;q=0.9"},
			}
		case "FIREFOX":
			return http.Header{
				"user-agent":      {UserAgent},
				"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
				"accept-language": {"en-GB,en;q=0.5"},
				"dnt":             {"1"},
			}
		default:
			//should loop and return chrome browser headers
			Browser = "CHROME"
		}
	}
}

func AssignPostRequestHeaders(Browser, UserAgent, postURL string) http.Header {
	contentType := getContentType(postURL)
	siteOrigin := returnOrigin(postURL)

	siteOrigin = strings.TrimRight(siteOrigin, "/") //removes trailing /
	for {
		switch Browser {
		case "CHROME":
			return http.Header{
				"cache-control":             {"max-age=0"},
				"sec-ch-ua":                 {`" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`},
				"sec-ch-ua-mobile":          {"?0"},
				"sec-ch-ua-platform":        {"\"Windows\""},
				"upgrade-insecure-requests": {"1"},
				"origin":                    {siteOrigin},  //just https://my.asos.com for login
				"content-type":              {contentType}, //application/x-www-form-urlencoded for login
				"user-agent":                {UserAgent},
				"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"sec-fetch-site":            {"same-origin"},
				"sec-fetch-mode":            {"navigate"},
				"sec-fetch-user":            {"?1"},
				"sec-fetch-dest":            {"document"},
				"accept-language":           {"en-GB,en;q=0.9"},
			}
		case "FIREFOX":
			return http.Header{
				"user-agent":                {UserAgent},
				"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
				"accept-language":           {"en-GB,en;q=0.5"},
				"content-type":              {contentType},
				"origin":                    {siteOrigin},
				"dnt":                       {"1"},
				"referer":                   {postURL},
				"upgrade-insecure-requests": {"1"},
				"sec-fetch-dest":            {"document"},
				"sec-fetch-mode":            {"navigate"},
				"sec-fetch-site":            {"same-origin"},
				"sec-fetch-user":            {"?1"},
				"te":                        {"trailers"},
			}
		default:
			//should loop and return chrome browser headers
			Browser = "CHROME"
		}
	}
}
