package asos

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/FergusJJ/coursework/gosrc/utils"
	http "github.com/useflyent/fhttp"
)

type stInteractions struct {
	ApplicationId           string
	DeviceId                string
	DeviceType              string
	AppSessionId            string
	StToken                 string
	KeyboardInteractions    []string
	MouseInteractions       []string
	IndirectEventsPayload   []string
	IndirectEventsCounters  map[string]string
	Gestures                []string
	MetricsData             map[string]string
	AccelerometerData       []string
	GyroscopeData           []string
	LinearAccelerometerData []string
	RotationData            []string
	Index                   int
	PayloadId               string
	Tags                    stTags
	Environment             stEnviroment
	IsMobile                bool
	UsernameTs              int
	Username                string
}

type stTags struct {
	Name      string
	EpochTs   int
	Timestamp int
}

type stEnviroment struct {
	Ops              int
	WebGl            string
	DevicePixelRatio int
	ScreenWidth      int
	ScreenHeight     int
}

func (s *localAsosSession) sendInteractions() error {

	stInteractions := s.createStInteractionsPayload()

	stringifiedData := fmt.Sprintf(
		`{"applicationId":"%s","deviceId":"%s","deviceType":"%s","appSessionId":"%s","stToken":"%s","keyboardInteractionPayloads":[],"mouseInteractionPayloads":[],"indirectEventsPayload":[],"indirectEventsCounters":{},"gestures":[],"metricsData":{},"accelerometerData":[],"gyroscopeData":[],"linearAccelerometerData":[],"rotationData":[],"index":%d,"payloadId":"%s","tags":[{"name":"location:%s","epochTs":%d,"timestamp":%d}],"environment":{"ops":%d,"webGl":"","devicePixelRatio":%d,"screenWidth":%d,"screenHeight":%d},"isMobile":%t,"usernameTs":%d,"username":"%s"}`,
		stInteractions.ApplicationId,
		stInteractions.DeviceId,
		stInteractions.DeviceType,
		stInteractions.AppSessionId,
		s.StLoginToken,
		//payloads going to be static and empty so skip those
		stInteractions.Index,
		stInteractions.PayloadId,
		stInteractions.Tags.Name,
		stInteractions.Tags.EpochTs,
		stInteractions.Tags.Timestamp,
		stInteractions.Environment.Ops,
		stInteractions.Environment.DevicePixelRatio,
		stInteractions.Environment.ScreenWidth,
		stInteractions.Environment.ScreenHeight,
		stInteractions.IsMobile,
		stInteractions.UsernameTs,
		stInteractions.Username,
	)
	stInteractions.Index++
	postBody := formatPayloadData(stringifiedData)
	req, err := http.NewRequest(http.MethodPost, "https://st.asos.com/SecuredTouch/rest/services/v2/interactions/asos", bytes.NewBuffer(postBody))
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return err
	}
	req.Header = http.Header{
		"content-length":     {strconv.Itoa(len(postBody))},
		"sec-ch-ua":          {`" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`},
		"content-encoding":   {"gzip"},
		"attempt":            {"0"},
		"sec-ch-ua-mobile":   {"?0"},
		"authorization":      {s.SecuredTouchSession.StToken},
		"clientepoch":        {fmt.Sprintf("%d", time.Now().UnixMilli())},
		"content-type":       {"application/json"},
		"accept":             {"application/json"},
		"instanceuuid":       {helperStGuid()},
		"user-agent":         {s.UserAgent},
		"encrypted":          {"1"},
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
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return err
			}
			err = s.sendInteractions()
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst submitting Secured Touch \n"))
			err = s.sendInteractions()
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
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen("Secured Touch successfully\n"))
			return nil
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Something went wrong submitting Secured Touch - HTTP Response %d\n", resp.StatusCode)))
			return err
		}
	}
}

func (s *localAsosSession) createStInteractionsPayload() *stInteractions {
	newInteractions := &stInteractions{}
	newInteractions.ApplicationId = "asos"
	newInteractions.DeviceId = "Id-" + helperStGuid()
	newInteractions.DeviceType = s.interactionsGetDeviceType()
	newInteractions.AppSessionId = s.SecuredTouchSession.StSessionId
	newInteractions.StToken = s.SecuredTouchSession.StToken
	newInteractions.KeyboardInteractions = []string{}
	newInteractions.MouseInteractions = []string{}
	newInteractions.IndirectEventsPayload = []string{}
	newInteractions.IndirectEventsCounters = map[string]string{}
	newInteractions.Gestures = []string{}
	newInteractions.MetricsData = map[string]string{}
	newInteractions.AccelerometerData = []string{}
	newInteractions.GyroscopeData = []string{}
	newInteractions.LinearAccelerometerData = []string{}
	newInteractions.RotationData = []string{}
	newInteractions.Index = s.SecuredTouchSession.Index
	newInteractions.PayloadId = helperStGuid()
	newInteractions.Tags = s.interactionsStGetTags()
	newInteractions.Environment = s.interactionsGetEnviromentData()
	newInteractions.IsMobile = false
	newInteractions.UsernameTs = newInteractions.Tags.Timestamp - rand.Intn(8)
	newInteractions.Username = s.SecuredTouchSession.Username
	if len(s.SecuredTouchSession.StDeviceId) == 0 {
		newInteractions.DeviceId = "Id-" + helperStGuid()
		s.SecuredTouchSession.StDeviceId = "Id-" + helperStGuid()
	} else {
		newInteractions.DeviceId = s.SecuredTouchSession.StDeviceId
	}
	if s.SecuredTouchSession.Username == "" {
		newInteractions.Username = newInteractions.DeviceId
		s.SecuredTouchSession.StDeviceId = newInteractions.DeviceId
	}

	return newInteractions
}

func (s *localAsosSession) interactionsStGetTags() (tags stTags) {
	nameTag := s.LoginURL
	if strings.Contains(s.LoginURL, "&") {
		splitAtQuery := strings.Split(s.LoginURL, "&")
		nameTag = splitAtQuery[0]
	}
	tags = stTags{
		Name:      nameTag,
		EpochTs:   int(time.Now().UnixMilli()),
		Timestamp: int(time.Now().UnixMilli()),
	}
	return tags
}

func formatPayloadData(stringifiedData string) []byte {

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	_, err := gzWriter.Write([]byte(stringifiedData))
	if err != nil {
		log.Fatal(err)
	}
	if err := gzWriter.Close(); err != nil {
		log.Fatal(err)
	}
	//uint8 array
	bytesArray := buf.Bytes()
	encryptionBytesString := "eG9yLWVuY3J5cHRpb24"
	postBody := encryptBytes(bytesArray, encryptionBytesString)
	return postBody
}

func encryptBytes(byteSlice []byte, encryptionString string) []byte {
	newByteSlice := []byte{}
	for counter := 0; counter < len(byteSlice); counter++ {
		newByte := int(byteSlice[counter]) ^ int([]rune(encryptionString)[counter%len(encryptionString)]) //  encryptionString.charCodeAt(counter % encryptionString.length); js equivalent
		newByteSlice = append(newByteSlice, byte(newByte))
	}
	return newByteSlice

}

func (s *localAsosSession) interactionsGetDeviceType() string {
	os := ""
	browser := ""
	verSubString := ""
	osVer1 := ""

	if strings.Contains(s.UserAgent, "Edge") {
		browser = "Edge"
		split1 := strings.Split(s.UserAgent, "Edge/")
		versionarr := strings.Split(split1[1], " ")
		verSubString = versionarr[0]
	} else if strings.Contains(s.UserAgent, "Firefox") {
		browser = "Firefox"
		split1 := strings.Split(s.UserAgent, "Firefox/")
		versionarr := strings.Split(split1[1], " ")
		verSubString = versionarr[0]
	} else if strings.Contains(s.UserAgent, "Safari") && !strings.Contains(s.UserAgent, "Chrome") {
		browser = "Safari"
		split1 := strings.Split(s.UserAgent, "Safari/")
		versionarr := strings.Split(split1[1], " ")
		verSubString = versionarr[0]
	} else {
		browser = "Chrome"
		split1 := strings.Split(s.UserAgent, "Chrome/")
		versionarr := strings.Split(split1[1], " ")
		verSubString = versionarr[0]
	}

	if strings.Contains(s.UserAgent, "Macintosh") {
		os = "Mac OS"
		os_ver := strings.Split(s.UserAgent, "X ")
		os_ver = strings.Split(os_ver[1], ")")
		osVer1 = strings.Replace(os_ver[0], "_", ".", 1)
	} else if strings.Contains(s.UserAgent, "Linux") {
		os = "Linux"
		os_ver := strings.Split(s.UserAgent, "Linux ")
		os_ver = strings.Split(os_ver[1], ")")
		osVer1 = os_ver[0]
	} else {
		os = "Windows"
		os_ver := strings.Split(s.UserAgent, "NT ")
		os_ver = strings.Split(os_ver[1], ";")
		osVer1 = os_ver[0]
	}
	stDeviceType := fmt.Sprintf("%s(%s)-%s(%s)", browser, verSubString, os, osVer1)
	return stDeviceType
}

func (s *localAsosSession) interactionsGetEnviromentData() (data stEnviroment) {

	data = stEnviroment{
		Ops:              0,
		WebGl:            "",
		DevicePixelRatio: 1,
		ScreenWidth:      int(s.AbckDevice.Gd.Width),
		ScreenHeight:     int(s.AbckDevice.Gd.Height),
	}

	return data
}
