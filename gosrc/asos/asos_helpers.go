package asos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/CrimsonAIO/radix"
	utils "github.com/FergusJJ/coursework/gosrc/utils"
	http "github.com/useflyent/fhttp"
	"golang.org/x/net/html"
)

const HOME_URL = "https://www.asos.com/"

const LOGIN_URL = "https://my.asos.com/my-account"

const KEY_STORE_DATA_VERSION = "hgk0y12-29"

const CHROME_CLIENT_CONFIG = "D91F2DAA-898C-4E10-9102-D6C974AFBD59"

func parseHTML(typeOfField, fieldName string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == typeOfField && a.Val == fieldName {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = parseHTML(typeOfField, fieldName, c); ok {
			return
		}
	}
	return
}

func parseFormInputValues(respBody io.ReadCloser, elementId, elementAttribute, typeOfField string) (string, error) {
	root, err := html.Parse(respBody)
	if err != nil {
		log.Fatal(err)
	}
	idValue, ok := parseHTML(typeOfField, elementId, root)
	if !ok {
		return "", fmt.Errorf("couldn't parse %s", elementId)
	}
	for _, v := range idValue.Attr {
		if v.Key == elementAttribute {
			return v.Val, nil
		}
	}
	return "", fmt.Errorf("couldn't find value in attributes of %s", elementId)
}

func formatLoginBody(xsrf, securedTouchToken, username, password string) url.Values {
	data := url.Values{}
	data.Set("idsrv.xsrf", xsrf)
	data.Set("SecuredTouchToken", securedTouchToken)
	data.Set("Username", username)
	data.Set("Password", password)
	return data
}

func getSecondLoginFormValues(responseData []byte) (formValues url.Values, accessToken string, err error) {

	respBody := ioutil.NopCloser(bytes.NewBuffer(responseData))
	idToken, err := parseFormInputValues(respBody, "id_token", "value", "name")
	if err != nil {
		return url.Values{}, "", err
	}
	respBody = ioutil.NopCloser(bytes.NewBuffer(responseData))
	accessToken, _ = parseFormInputValues(respBody, "access_token", "value", "name")

	respBody = ioutil.NopCloser(bytes.NewBuffer(responseData))
	tokenType, _ := parseFormInputValues(respBody, "token_type", "value", "name")

	respBody = ioutil.NopCloser(bytes.NewBuffer(responseData))
	expiresIn, _ := parseFormInputValues(respBody, "expires_in", "value", "name")

	respBody = ioutil.NopCloser(bytes.NewBuffer(responseData))
	scope, _ := parseFormInputValues(respBody, "scope", "value", "name")

	respBody = ioutil.NopCloser(bytes.NewBuffer(responseData))
	state, _ := parseFormInputValues(respBody, "state", "value", "name")

	respBody = ioutil.NopCloser(bytes.NewBuffer(responseData))
	sessionState, _ := parseFormInputValues(respBody, "session_state", "value", "name")

	formData := url.Values{}
	formData.Set("id_token", idToken)
	formData.Set("access_token", accessToken)
	formData.Set("token_type", tokenType)
	formData.Set("expires_in", expiresIn)
	formData.Set("scope", scope)
	formData.Set("state", state)
	formData.Set("session_state", sessionState)
	return formData, accessToken, nil

}

func parseCid(body string) string {
	splitAtCid := strings.Split(body, "clientId")
	splitAtQuotation := strings.Split(splitAtCid[1], "\"")
	return splitAtQuotation[2]
}

//check genlinkn and convert to go
func genRand(i int64) (randNumStr string) {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano() + i))
	randNum := (float64((time.Now().UnixMilli())) + newRand.Float64()) * newRand.Float64()
	randNumStr = fmt.Sprintf("%.4f", randNum)
	randNumStr = strings.Replace(randNumStr, ".", "", 1)
	return randNumStr
}

func (s *localAsosSession) generateATCAuthLink() string {
	link := "https://my.asos.com/identity/connect/authorize?"
	queryMap1 := map[string]string{
		"client_id":     CHROME_CLIENT_CONFIG,
		"redirect_uri":  s.Profile.ProductURL,
		"response_type": "id_token token",
		"scope":         "openid sensitive profile",
	}
	queryMap2 := map[string]string{
		"ui_locales":    "en-GB",
		"acr_values":    "0",
		"response_mode": "json",
	}

	rand1 := genRand(time.Now().UnixNano())
	link = fmt.Sprintf("%sstate=%s", link, url.QueryEscape(rand1))
	rand2 := genRand(time.Now().UnixNano())
	link = fmt.Sprintf("%s&nonce=%s", link, url.QueryEscape(rand2))
	for k, v := range queryMap1 {
		link = fmt.Sprintf("%s&%s=%s", link, k, url.QueryEscape(v))

	}
	for k, v := range queryMap2 {
		link = fmt.Sprintf("%s&%s=%s", link, k, url.QueryEscape(v))
	}
	link = fmt.Sprintf("%s&store=COM&country=GB&keyStoreDataversion=%s&lang=en-GB", link, KEY_STORE_DATA_VERSION)
	return link
}

func helperStGuid() string {
	baseString := "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, currentChar := range baseString {
		if string(currentChar) != "4" && string(currentChar) != "-" {
			t := newRand.Intn(16)
			var newChar string
			if currentChar == 'x' {
				newChar = radix.ToString(float64(t), 16)
			} else {
				newChar = radix.ToString((float64(3&t | 8)), 16)
			}
			baseString = strings.Replace(baseString, string(currentChar), newChar, 1)
		}
	}
	return baseString
}

func (s *localAsosSession) getDeviceType() string {
	if strings.Contains(s.UserAgent, "Chrome") {
		split1 := strings.Split(s.UserAgent, "Chrome/")
		versionarr := strings.Split(split1[1], " ")
		verSubString := versionarr[0]
		return fmt.Sprintf("Chrome(%s)-Windows(10)", verSubString)
	} else {
		log.Fatalf("no version for user agent: %s", s.UserAgent)
		return ""
	}

}

func (s *localAsosSession) switchProxy() error {
	newProxy := utils.AddProxy()
	s.Profile.Proxy = newProxy
	formattedProxy := utils.FormatProxy(s.Profile.Proxy)
	proxyUrl, err := url.Parse(formattedProxy)
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return err
	}
	s.AsosClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	return nil
}
