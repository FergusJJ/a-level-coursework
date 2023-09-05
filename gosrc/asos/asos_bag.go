package asos

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	http "github.com/useflyent/fhttp"

	fhttp "github.com/AlienRecall/fhttp"
	utils "github.com/FergusJJ/coursework/gosrc/utils"
)

//returns the bag id used when adding items to cart via their variant id's
func (s *localAsosSession) getBagIDForATC() (string, error) {
	createNewBagLink := s.getExistingBagLink()
	req, err := http.NewRequest(http.MethodGet, createNewBagLink, nil)
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return "", err
	}
	req.Header = formatNewBagHeaders(string(s.AbckDevice.SelectedBrowser), s.UserAgent, createNewBagLink, s.Bearer.Access_token, s.Cid, s.Profile.ProductURL)
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return "", err
			}
			s.BagId, err = s.getBagIDForATC()
			if err != nil {
				return "", err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting a new bag \n"))
			s.BagId, err = s.getBagIDForATC()
			if err != nil {
				return "", err
			}

		}
		return s.BagId, nil
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow("Requesting new bag\n"))
			//bag needs to be created
			s.BagId, err = s.createNewBag()
			if err != nil {
				return "", err
			}
			return s.BagId, nil
		} else {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			s.BagId = utils.GetJsonValueFromString("bagId", string(bodyBytes))
			//
			return s.BagId, nil
		}
	}
}

//formats the link used to fetch an existing bag that is already connected to an account
func (s *localAsosSession) getExistingBagLink() string {
	tokenParts := strings.Split(s.Bearer.Access_token, ".")
	payloadB64 := tokenParts[1]
	sDec, _ := b64.StdEncoding.DecodeString(payloadB64)
	sub := utils.GetJsonValueFromString("sub", string(sDec))
	newBagLink := fmt.Sprintf("https://www.asos.com/api/commerce/bag/v4/customers/%s/countries/GB/bag/total?keyStoreDataversion=%s&lang=en-GB&cb=%d", sub, KEY_STORE_DATA_VERSION, time.Now().UnixMilli())
	//^this will get an existing bag for an account. If the bag doesn't exist then 404 is returned and the json response is the following:
	//newBagLink := fmt.Sprintf("https://www.asos.com/api/commerce/bag/v4/bags/%s?expand=summary,total&lang=en-GB&cb=%d", sub, time.Now().UnixMilli())
	//sub may need to be replaced with bagid once fetched. It looks like this link is for getting an existing bag and the one above returns a new bag -> returns 404 like in the browser as is
	// fmt.Println(newBagLink)
	return newBagLink
}

//request sent in the event that a bag does not exist for the account being used
func (s *localAsosSession) createNewBag() (string, error) {

	bagLink := s.postBagLink()
	var data = strings.NewReader(`{"currency":"GBP","lang":"en-GB","sizeSchema":"UK","country":"GB","originCountry":"GB"}`)
	bagId, err := s.createBag_new(bagLink, data)
	if err != nil {
		return "", err
	}
	return bagId, nil

}

//returns the link needed to send POST request in order to generate a new bag
func (s *localAsosSession) postBagLink() string {
	guid := ""
	u, _ := url.Parse(LOGIN_URL)
	for _, c := range s.AsosClient.Jar.Cookies(u) {
		if strings.Contains(c.Name, "cgd") {
			guid = c.Value
		}
	}
	var postBagLink string
	if len(guid) == 0 {
		tokenParts := strings.Split(s.Bearer.Access_token, ".")
		payloadB64 := tokenParts[1]
		sDec, _ := b64.StdEncoding.DecodeString(payloadB64)
		sub := utils.GetJsonValueFromString("sub", string(sDec))
		postBagLink = fmt.Sprintf("https://www.asos.com/api/commerce/bag/v4/customers/%s/bags/getbag?expand=summary,total&lang=en-GB&keyStoreDataversion=hgk0y12-29", sub)
	} else {
		postBagLink = fmt.Sprintf("https://www.asos.com/api/commerce/bag/v4/customers/%s/bags/getbag?expand=summary,total&lang=en-GB&keyStoreDataversion=hgk0y12-29", guid)
	}
	return postBagLink
}

//formats the request headers for the GET request used to return existing bag id's
func formatNewBagHeaders(Browser, UserAgent, URL, bearer, cid, productLink string) http.Header {
	bearerHeader := "Bearer " + bearer
	for {
		switch Browser {
		case "CHROME":
			return http.Header{

				"sec-ch-ua":          {`" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`},
				"x-requested-with":   {"XMLHttpRequest"},
				"asos-c-plat":        {"Web"},
				"authorization":      {bearerHeader},
				"asos-c-name":        {"Asos.Commerce.Bag.Sdk"},
				"content-type":       {"application/json"},
				"accept":             {"application/json, text/javascript, */*; q=0.01"},
				"asos-c-store":       {"COM"},
				"sec-ch-ua-mobile":   {"?0"},
				"asos-c-ver":         {"5.5.166"},
				"asos-c-istablet":    {"false"},
				"user-agent":         {UserAgent},
				"asos-cid":           {cid},
				"asos-c-ismobile":    {"false"},
				"sec-ch-ua-platform": {"\"Windows\""},
				//origin
				"sec-fetch-site":  {"same-origin"},
				"sec-fetch-mode":  {"cors"},
				"sec-fetch-dest":  {"empty"},
				"referer":         {productLink},
				"accept-encoding": {"gzip, deflate, br"},
				"accept-language": {"en-GB,en;q=0.9"},
			}
		default:
			//should loop and return chrome browser headers
			Browser = "CHROME"
		}
	}
}

//gets the bearer JWT (access token) that is used in the authorization header in the requests used to atc & get bag id's
func (s *localAsosSession) getATCBearer() error {
	bearerLink := s.generateATCAuthLink()
	req, _ := http.NewRequest(http.MethodGet, bearerLink, nil)
	req.Header = http.Header{
		"sec-ch-ua":                 {`" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"Windows\""},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br"},
		"accept-language":           {"en-GB,en;q=0.9"},
		fhttp.HeaderOrderKey: {
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
			"accept-language",
		},
		fhttp.PHeaderOrderKey: {
			":method", ":authority", ":scheme", ":path",
		},
	}
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return err
			}
			err = s.getATCBearer()
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting a bearer token \n"))
			err = s.getATCBearer()
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		bearerResponse := &localAsosBearerResponse{}
		json.Unmarshal(b, bearerResponse)
		s.Bearer = (*utils.AsosBearerResponse)(bearerResponse)
		parsedBearerLink, _ := url.Parse(bearerLink)
		s.setCookiesToClient(parsedBearerLink, s.AkamaiClient.Jar.Cookies(parsedBearerLink), []string{"*"})
		return nil
	}
}

func (s *localAsosSession) createBag_new(bagLink string, data *strings.Reader) (string, error) {
	if len(s.Cid) == 0 {
		s.Cid = helperStGuid()
	}

	req, err := http.NewRequest("POST", bagLink, data)
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return "", err
	}
	req.Header = http.Header{
		"sec-ch-ua":          {`" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`},
		"x-requested-with":   {"XMLHttpRequest"},
		"asos-c-plat":        {"Web"},
		"authorization":      {"Bearer " + s.Bearer.Access_token},
		"asos-c-name":        {"Asos.Commerce.Bag.Sdk"},
		"content-type":       {"application/json"},
		"accept":             {"application/json, text/javascript, */*; q=0.01"},
		"asos-c-store":       {"COM"},
		"sec-ch-ua-mobile":   {"?0"},
		"asos-c-ver":         {"5.5.166"},
		"asos-c-istablet":    {"false"},
		"user-agent":         {s.UserAgent},
		"asos-cid":           {s.Cid},
		"asos-c-ismobile":    {"false"},
		"sec-ch-ua-platform": {"\"Windows\""},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {s.Profile.ProductURL},
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
			bagIdSlice, err := s.getBagIDForATC()
			if err != nil {
				return "", err
			}
			return bagIdSlice, nil
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Something went wrong getting a new bag - HTTP Response %d\n", resp.StatusCode)))
			bagIdSlice, err := s.getBagIDForATC()
			if err != nil {
				return "", err
			}
			return bagIdSlice, nil

		}
	} else {
		defer resp.Body.Close()
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.SomethingWentWrong(s.Profile.TaskID)
			return "", err
		}
		splitAtId := strings.Split(string(bodyText), `"id": "`)
		bagIdSlice := strings.Split(splitAtId[1], `"`)
		return bagIdSlice[0], nil
	}

}
