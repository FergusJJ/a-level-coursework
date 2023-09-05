package asos

//
import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/FergusJJ/coursework/gosrc/utils"
	http "github.com/useflyent/fhttp"

	fhttp "github.com/AlienRecall/fhttp"
)

func getStSessionIdFromScript(body string) (value string) {

	splitAtJsKey := strings.Split(body, "sessionId:")
	value = (strings.Split(splitAtJsKey[1], "'"))[1]
	return value
}

func (s *localAsosSession) getSecuredTouchToken() (string, error) {
	unencodedQuery := fmt.Sprintf(`{"pingVersion":"1.3.0p","appId":"asos","appSessionId":"%s"}`, s.SecuredTouchSession.StSessionId)

	encodedQuery := getB64URLQuery(unencodedQuery)
	requestURL := fmt.Sprintf("https://st-static.asos.com/sdk/pong.js?body=%s", encodedQuery)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header = http.Header{
		"sec-ch-ua":          {`" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {s.UserAgent},
		"sec-ch-ua-platform": {"\"Windows\""},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"same-site"},
		"sec-fetch-mode":     {"no-cors"},
		"sec-fetch-dest":     {"script"},
		"referer":            {"https://my.asos.com/"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en-GB,en;q=0.9"},
	}
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return "", err
			}
			stToken, err := s.getSecuredTouchToken()
			if err != nil {
				return "", err
			}
			return stToken, nil
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting Secured Touch bearer\n"))
			stToken, err := s.getSecuredTouchToken()
			if err != nil {
				return "", err
			}
			return stToken, nil
		}
	} else {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		responseBodyString := string(bodyBytes)
		splitted := strings.Split(responseBodyString, "'")
		return splitted[1], nil
	}

}

func getB64URLQuery(unencoded string) (encoded string) {
	encoded = b64.StdEncoding.EncodeToString([]byte(unencoded))
	encoded = url.QueryEscape(encoded)
	return encoded
	// encoded, _ = b64.StdEncoding.Encode(unencoded)
}

func (s *localAsosSession) stSendStarter() error {
	type stStarterObject struct {
		DeviceId      string `json:"device_id"`
		ClientVersion string `json:"clientVersion"`
		DeviceType    string `json:"deviceType"`
		AuthToken     string `json:"authToken"`
	}
	appSecretBase64 := "YjIxMzVjdDIxSnVsVnlP"
	postLink := "https://st.asos.com/SecuredTouch/rest/services/v2/starter/asos"
	if s.SecuredTouchSession.StDeviceId == "" {
		s.SecuredTouchSession.StDeviceId = "Id-" + helperStGuid()
	}
	stStarter := stStarterObject{
		DeviceId:      s.SecuredTouchSession.StDeviceId,
		ClientVersion: "3.13.2w",
		DeviceType:    s.getDeviceType(),
		AuthToken:     "",
	}
	requestBody, err := json.Marshal(stStarter)
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return err
	}
	req, err := fhttp.NewRequest(fhttp.MethodPost, postLink, bytes.NewBuffer(requestBody))
	if err != nil {

		utils.SomethingWentWrong(s.Profile.TaskID)
		return err
	}
	req.Header = fhttp.Header{
		"content-length":     {strconv.Itoa(len(requestBody))},
		"sec-ch-ua":          {`" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`},
		"content-encoding":   {"gzip"},
		"attempt":            {"0"},
		"sec-ch-ua-mobile":   {"?0"},
		"authorization":      {appSecretBase64},
		"clientepoch":        {fmt.Sprintf("%d", time.Now().UnixMilli())},
		"content-type":       {"application/json"},
		"accept":             {"application/json"},
		"instanceuuid":       {helperStGuid()},
		"user-agent":         {s.UserAgent},
		"clientversion":      {"3.13.2w"},
		"sec-ch-ua-platform": {`"Windows"`},
		"origin":             {"https://my.asos.com"},
		"sec-fetch-site":     {"same-site"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://my.asos.com/"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en-GB,en;q=0.9"},
	}
	resp, err := s.AkamaiClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return err
			}
			err = s.stSendStarter()
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst submitting Secured Touch \n"))
			err = s.stSendStarter()
			if err != nil {
				return err
			}

		}
		return nil
	} else {
		defer resp.Body.Close()
		if resp.StatusCode < 300 {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				utils.SomethingWentWrong(s.Profile.TaskID)
				return err
			}
			bodyStr := string(data)
			splitAtSpeech := strings.Split(bodyStr, `"`)
			s.SecuredTouchSession.StToken = splitAtSpeech[3]
			s.SecuredTouchSession.StDeviceId = splitAtSpeech[11]
			return nil
		} else {
			fmt.Print(utils.ColourRed(fmt.Sprintf("[%s] | %s | Something went wrong submitting Secured Touch - HTTP Response %d\n", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, resp.StatusCode)))
			return err
			//handle error, try again with new data
		}
	}

}
