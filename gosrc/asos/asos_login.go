package asos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	http "github.com/useflyent/fhttp"

	fhttp "github.com/AlienRecall/fhttp"
	utils "github.com/FergusJJ/coursework/gosrc/utils"
	"github.com/FergusJJ/coursework/gosrc/webhook"
	akamaisensor "github.com/FergusJJ/go-sensor"
	akamaisensorutils "github.com/FergusJJ/go-sensor/utils"
)

type localAsosBearerResponse utils.AsosBearerResponse
type localAsosSession utils.AsosSession

func InitTaskReceiver(profile *utils.Profile, akamaiClient fhttp.Client, DefaultHttpClient http.Client, ua, webhookUrl string) {
	//for loop could be used to wrap this code in order to loop until the task has checked out
	TaskData := &localAsosSession{
		Profile:             profile,
		AkamaiClient:        &akamaiClient,
		AsosClient:          &DefaultHttpClient,
		UserAgent:           ua,
		SecuredTouchSession: &utils.SecuredTouchSession{},
		WebhookInfo:         &utils.WebhookInfo{},
		VariantID:           0,
		IsCheckedOut:        false,
		IsLoggedIn:          false,
	}
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), TaskData.Profile.TaskID, utils.ColourYellow("Getting Akamai cookies...\n"))
	TaskData.GetAbck(LOGIN_URL, "https://my.asos.com/QJUhtQ2Q/L-H1QGA/90A827E/I5/7pEQfhLbEQaY/YwpHIloC/ckFcNg/djU1k")
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), TaskData.Profile.TaskID, utils.ColourGreen("Akamai handled successfully\n"))
	err := TaskData.sendReceiverHomepageRequest()
	if err != nil {
		return
	}
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), TaskData.Profile.TaskID, utils.ColourYellow("Getting login page...\n"))
	loginURL, err := TaskData.getReceiverLoginPageAsos()
	if err != nil {
		return
	}
	TaskData.LoginURL = loginURL
	TaskData.SecuredTouchSession.Index = 0
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), TaskData.Profile.TaskID, utils.ColourYellow("Submitting Secured Touch...\n"))
	if err = TaskData.stSendStarter(); err != nil {
		return
	}
	if err = TaskData.sendInteractions(); err != nil {
		return
	}

	TaskData.loginAsos(TaskData.LoginURL)
	TaskData.getProductPage()
	if len(webhookUrl) > 0 && TaskData.IsCheckedOut {
		webhook.SendWebhook(TaskData.WebhookInfo.ItemName,
			TaskData.Profile.ProductURL,
			TaskData.WebhookInfo.ImageLink,
			TaskData.Profile.Store,
			TaskData.Profile.Proxy,
			TaskData.Profile.Size,
			TaskData.WebhookInfo.Pid,
			TaskData.WebhookInfo.Price,
			webhookUrl,
			profile.TaskID,
		)
	}
}

func (s *localAsosSession) getReceiverLoginPageAsos() (loginURL string, err error) {

	req, err := http.NewRequest(http.MethodGet, LOGIN_URL, nil)
	if err != nil {
		fmt.Println(err)
	}
	requestHeaders := utils.AssignGetRequestHeaders(string(s.AbckDevice.SelectedBrowser), s.UserAgent) //s.AbckDevice.SelectedBrowser
	req.Header = requestHeaders
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return "", err
			}
			loginURL, err = s.getReceiverLoginPageAsos()
			if err != nil {
				return "", err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting the homepage \n"))
			loginURL, err = s.getReceiverLoginPageAsos()
			if err != nil {
				return "", err
			}

		}
		return loginURL, nil
	} else {
		defer resp.Body.Close()
		statusOk := utils.CheckResponseStatus(resp.StatusCode)
		if !statusOk {
			//forbidden, will occur if a proxy has been banned & now need a new proxy
			if resp.StatusCode == 403 {
				if s.Profile.Proxy != "localhost" {
					fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
					newProxy := utils.AddProxy()
					s.Profile.Proxy = newProxy
					formattedProxy := utils.FormatProxy(s.Profile.Proxy)
					proxyUrl, err := url.Parse(formattedProxy)
					if err != nil {
						utils.SomethingWentWrong(s.Profile.TaskID)
					}
					s.AsosClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
					utils.DelayRequest(s.Profile.Delay)
					s.getReceiverLoginPageAsos()
				} else {
					fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Something went wrong getting login - HTTP Response %d\n", resp.StatusCode)))
					utils.DelayRequest(s.Profile.Delay)
					s.getReceiverLoginPageAsos()
				}
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen("Got login page\n"))
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				utils.SomethingWentWrong(s.Profile.TaskID)
				return "", err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			xsrfToken, err := parseFormInputValues(resp.Body, "idsrv_xsrf", "value", "id")
			if err != nil {
				utils.SomethingWentWrong(s.Profile.TaskID)
				return "", err
			}

			resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			postURL, err := parseFormInputValues(resp.Body, "signInForm", "action", "id")
			if err != nil {
				utils.SomethingWentWrong(s.Profile.TaskID)
				return "", err
			}

			s.Idxsrf = xsrfToken
			//new
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				utils.SomethingWentWrong(s.Profile.TaskID)
				return "", err
			}
			s.SecuredTouchSession.StSessionId = getStSessionIdFromScript(string(bodyBytes))
			s.StLoginToken, err = s.getSecuredTouchToken()
			if err != nil {
				return "", err
			}
			return postURL, nil
		}
		return "", err
	}
}

func (s *localAsosSession) loginAsos(loginURL string) error {
	//may need a way to check for captcha on login.
	securedTouchJWT := s.StLoginToken
	postBodyData := formatLoginBody(s.Idxsrf, securedTouchJWT, s.Profile.Email, s.Profile.Password)
	req, _ := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(postBodyData.Encode()))
	requestHeaders := utils.AssignPostRequestHeaders(string(s.AbckDevice.SelectedBrowser), s.UserAgent, loginURL)
	req.Header = requestHeaders

	utils.DelayRequest(s.Profile.Delay)
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return err
			}
			err = s.loginAsos(loginURL)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst sending login \n"))
			err = s.loginAsos(loginURL)
			if err != nil {
				return err
			}

		}
		return nil
	} else {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.SomethingWentWrong(s.Profile.TaskID)
		}
		parsedUrl, _ := url.Parse(loginURL)
		s.setCookiesToClient(parsedUrl, s.AkamaiClient.Jar.Cookies(parsedUrl), []string{"*"})
		redirectForm, accessToken, err := getSecondLoginFormValues(data)
		s.LoginAccessToken = accessToken
		if err != nil {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Couldn't log in, trying guest checkout - %s\n", s.Profile.Email)))
		} else {
			ok := s.postSecondLoginForm(redirectForm)
			if !ok {
				utils.SomethingWentWrong(s.Profile.TaskID)
				return nil
			} else {
				fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen(fmt.Sprintf("Successfully logged in - %s\n", s.Profile.Email)))
				s.IsLoggedIn = true
			}
		}
	}
	return nil

}

func (s *localAsosSession) sendReceiverHomepageRequest() error {
	req, err := http.NewRequest(http.MethodGet, HOME_URL, nil)
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return err
	}
	requestHeaders := utils.AssignGetRequestHeaders(string(s.AbckDevice.SelectedBrowser), s.UserAgent)

	req.Header = requestHeaders
	resp, err := s.AsosClient.Do(req)

	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return err
			}
			err = s.sendReceiverHomepageRequest()
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting the homepage \n"))
			err = s.sendReceiverHomepageRequest()
			if err != nil {
				return err
			}
		}
	} else {
		defer resp.Body.Close()
		statusOk := utils.CheckResponseStatus(resp.StatusCode)
		if !statusOk {
			//forbidden, will occur if a proxy has been banned & now need a new proxy
			if resp.StatusCode == 403 {
				if s.Profile.Proxy != "localhost" {
					fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
					newProxy := utils.AddProxy()
					s.Profile.Proxy = newProxy
					formattedProxy := utils.FormatProxy(s.Profile.Proxy)
					proxyUrl, err := url.Parse(formattedProxy)
					if err != nil {
						utils.SomethingWentWrong(s.Profile.TaskID)
						return err
					}
					s.AsosClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
					utils.DelayRequest(s.Profile.Delay)
					s.sendReceiverHomepageRequest()
				} else {
					fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Something went wrong getting the homepage - HTTP Response %d\n", resp.StatusCode)))
					utils.DelayRequest(s.Profile.Delay)
					s.sendReceiverHomepageRequest()
				}
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen("Got homepage\n"))
			return nil
		}
	}
	return nil
}

func (s *localAsosSession) postSecondLoginForm(formData url.Values) (ok bool) {
	postData := formData.Encode()
	req, err := http.NewRequest(http.MethodPost, LOGIN_URL, strings.NewReader(postData))
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return false
	}

	req.Header = utils.AssignPostRequestHeaders(string(s.AbckDevice.SelectedBrowser), s.UserAgent, LOGIN_URL)
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return false
			}
			err = s.getProductPage()
			if err != nil {
				return false
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst submitting second login form \n"))
			err = s.getProductPage()
			if err != nil {
				return false
			}

		}
		return false
	} else {
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.SomethingWentWrong(s.Profile.TaskID)
			return false
		}
		bodyString := string(bodyBytes)
		//parse cid
		s.Cid = parseCid(bodyString)
		parsedUrl, _ := url.Parse(LOGIN_URL)
		s.setCookiesToClient(parsedUrl, s.AkamaiClient.Jar.Cookies(parsedUrl), []string{"*"})
		return utils.CheckResponseStatus(resp.StatusCode)
	}
}

func (s *localAsosSession) GetAbck(currentTaskURL, akamaiUrl string) {
	s.AbckDevice = akamaisensor.SensorData(
		s.UserAgent,
		currentTaskURL,
		&akamaisensorutils.Device{},
		*s.AkamaiClient,
		akamaiUrl,
	)
	u, _ := url.Parse(currentTaskURL)

	s.setCookiesToClient(u, s.AkamaiClient.Jar.Cookies(u), []string{"_abck"})

}

func (s *localAsosSession) setCookiesToClient(u *url.URL, ckie []*fhttp.Cookie, include []string) {
	var newCookies []*http.Cookie
	for _, ckie := range ckie {
		for _, includeCookie := range include {
			if ckie.Name == includeCookie || includeCookie == "*" {

				setCkie := &http.Cookie{
					Name:       ckie.Name,
					Value:      ckie.Value,
					Path:       ckie.Path,
					Domain:     ckie.Domain,
					Expires:    ckie.Expires,
					RawExpires: ckie.RawExpires,
					MaxAge:     ckie.MaxAge,
					Secure:     ckie.Secure,
					HttpOnly:   ckie.HttpOnly,
					SameSite:   http.SameSite(ckie.SameSite),
					Raw:        ckie.Raw,
					Unparsed:   ckie.Unparsed,
				}
				newCookies = append(newCookies, setCkie)
				ckie.Expires = time.Unix(0, 0)
				ckie.MaxAge = -1
			}
		}
	}
	s.AsosClient.Jar.SetCookies(u, newCookies)

}
