package main

import "C"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	http "github.com/useflyent/fhttp"

	cookiejar "github.com/useflyent/fhttp/cookiejar"

	cclient "github.com/FergusJJ/Coursework_Golang_Client"
	"github.com/FergusJJ/coursework/gosrc/asos"
	"github.com/FergusJJ/coursework/gosrc/utils"
	"github.com/FergusJJ/coursework/gosrc/zalando"
)

var ProfileMap = make(map[int]*utils.Profile)

//export convertToGo
func convertToGo(_value []*C.char) {
	//slice of strings, used so that i can access the info stored by the index
	tempProfileSlice := make([]string, 0)

	//iterate through values
	for _, value := range _value {
		goValue := C.GoString(value)
		tempProfileSlice = append(tempProfileSlice, goValue)
	}
	newProfileCC := utils.PaymentDetails{
		CardNumber:  tempProfileSlice[16],
		ExpiryMonth: tempProfileSlice[17],
		ExpiryYear:  tempProfileSlice[18],
		CVC:         tempProfileSlice[19],
	}
	newProfile := utils.Profile{
		Store:        tempProfileSlice[0],
		Mode:         tempProfileSlice[1],
		ProductURL:   tempProfileSlice[2], //product
		Size:         tempProfileSlice[3], //size
		Delay:        tempProfileSlice[4],
		Proxy:        tempProfileSlice[5],
		FirstName:    tempProfileSlice[6],
		LastName:     tempProfileSlice[7],
		Email:        tempProfileSlice[8], //email
		Password:     tempProfileSlice[9], //password
		AddressLine1: tempProfileSlice[10],
		AddressLine2: tempProfileSlice[11],
		City:         tempProfileSlice[12],
		Province:     tempProfileSlice[13],
		Postcode:     tempProfileSlice[14],
		CountryCode:  tempProfileSlice[15], //country code
		CC:           newProfileCC,
		TaskID:       fmt.Sprintf("%d", len(ProfileMap)+1),
	}
	profileMapIndex := len(ProfileMap)

	ProfileMap[profileMapIndex] = &newProfile
}

//export checkProfileMap
func checkProfileMap() {
	SetupCloseHandler()
	currentTime := time.Now()
	timeStamp := fmt.Sprint(currentTime.Format("15:04:05"))
	profilesTotal := len(ProfileMap)
	fmt.Printf("[%s] %d Profiles ready...\n", timeStamp, profilesTotal)
	addRandomProxies(profilesTotal)
	startCorrectSite(profilesTotal)
}

//https://golangcode.com/handle-ctrl-c-exit-in-terminal/

func SetupCloseHandler() {
	//sets up a goroutine to monitor for a ctrl-c interrupt
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}

func addRandomProxies(profilesTotal int) {
	var Wg sync.WaitGroup
	utils.ProxiesToSlice() //appends to a map in utils#
	currentTime := time.Now()
	timeStamp := fmt.Sprint(currentTime.Format("15:04:05"))
	fmt.Printf("[%s] Assigning proxies...\n", timeStamp)
	proxiesAssigned := 0
	Wg.Add(profilesTotal)
	go func() {
		for i := 0; i < profilesTotal; i++ {
			if strings.ToLower(ProfileMap[i].Proxy) == "random" {
				proxy := utils.AddProxy() //ProfileMap[i].Proxy
				ProfileMap[i].Proxy = proxy
				proxiesAssigned++
			}
			Wg.Done() //will decrement the proxiesWg counter
		}
	}()
	Wg.Wait() //will wait for the proxiesWg counter to == 0
	currentTime = time.Now()
	timeStamp = fmt.Sprint(currentTime.Format("15:04:05"))
	fmt.Printf("[%v] %d Proxies assigned...\n", timeStamp, proxiesAssigned)
}

func readWebhookUrl() string {
	type settingsJson struct {
		WebhookUrl string `json:"discord_webhook"`
	}
	var userSettings settingsJson
	settingsFile, err := os.Open("settings.json")
	if err != nil {
		fmt.Println(err)
	}
	bytes, _ := ioutil.ReadAll(settingsFile)
	json.Unmarshal(bytes, &userSettings)
	return userSettings.WebhookUrl
}

func startCorrectSite(profilesTotal int) {
	var Wg sync.WaitGroup
	webhookUrl := readWebhookUrl()

	if strings.ToLower(ProfileMap[0].Store) == "zalando" {
		Wg.Add(profilesTotal)
		go func() {
			for i := 0; i < profilesTotal; i++ {
				zalando.InitTask(ProfileMap[i])
				Wg.Done()
			}
		}()
		Wg.Wait()
	}
	if strings.ToLower(ProfileMap[0].Store) == "asos" {
		Wg.Add(profilesTotal)
		useragent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < profilesTotal; i++ {
			go func(i int) {
				newProxy := utils.AddProxy()
				formattedProxy := utils.FormatProxy(newProxy)
				client, err := cclient.NewClient(useragent, formattedProxy)
				if err != nil {
					utils.SomethingWentWrong(fmt.Sprint(i))
				}

				formattedProxy = utils.FormatProxy(ProfileMap[i].Proxy)
				proxyUrl, err := url.Parse(formattedProxy)
				if err != nil {
					fmt.Print(utils.ColourRed(fmt.Sprintf("[%s] | %s | Error assigning proxy\n", utils.ReturnFormattedTimestamp(), ProfileMap[i])))
					Wg.Done()
					return
				}
				defaultClient := http.Client{
					Jar:       jar,
					Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
				}
				asos.InitTaskReceiver(ProfileMap[i], client, defaultClient, useragent, webhookUrl)
				Wg.Done()
			}(i)
		}

		Wg.Wait()
	} else {
		fmt.Println("No site module for that website")
	}
}

func main() {}
