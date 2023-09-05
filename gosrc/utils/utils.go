package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	fhttp "github.com/AlienRecall/fhttp"
	colour "github.com/gookit/color"
	http "github.com/useflyent/fhttp"
)

//will be able to access this whilst tasks are ran
var ProxySlice []string
var ColourGreen = colour.FgGreen.Render
var ColourRed = colour.FgRed.Render
var ColourYellow = colour.FgYellow.Render
var ColourGrey = colour.FgDarkGray.Render

func readIndividualFile(path string) []string {

	var proxies []string
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close() // makes sure that file is closed when function ends
	//scanner allows me to read lines of a file
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		proxies = append(proxies, fileScanner.Text())
	}
	return proxies
}

func ProxiesToSlice() []string {

	proxyFiles, err := ioutil.ReadDir("proxies/")
	if err != nil {
		fmt.Println(err)
	}
	for _, _file := range proxyFiles {
		var proxies []string
		relativePath := "proxies/" + _file.Name()
		proxies = readIndividualFile(relativePath)
		ProxySlice = append(ProxySlice, proxies...)
	}
	return ProxySlice
}
func AddProxy() (proxy string) {

	//get a proxy from the slice and add make it the proxy field
	indexOfProxy := rand.Intn(len(ProxySlice))
	_proxy := ProxySlice[indexOfProxy]
	ProxySlice[indexOfProxy] = ProxySlice[len(ProxySlice)-1] //last value in slice is used to replace current proxy
	ProxySlice[len(ProxySlice)-1] = ""                       //gets rid of last value
	ProxySlice = ProxySlice[:len(ProxySlice)-1]              //shortens slice

	return _proxy

}

func ReturnRandomUA() string {
	userAgent := ""

	indexOfUA := rand.Int63n(2)
	userAgent = userAgentsSlice[indexOfUA]
	return userAgent
}

func ReturnFormattedTimestamp() string {
	currentTime := time.Now()
	timestamp := fmt.Sprint(currentTime.Format("15:04:05"))
	return ColourGrey(fmt.Sprintf("[%s]", timestamp))
}

func DelayRequest(delay string) {
	intDelay, _ := strconv.Atoi(delay)
	time.Sleep(time.Millisecond * time.Duration(intDelay))
}

func GetJsonValueFromString(keyName, jsonString string) (value string) {
	splitAtSub := strings.Split(jsonString, keyName)
	value = (strings.Split(splitAtSub[1], "\""))[2]
	return value
}

var userAgentsSlice = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.53",
}

func SetDefaultChromeHeaders(req *http.Request, userAgent, secSite, secMode, secDest string) {

	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`)
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-fetch-site", secSite)
	req.Header.Set("sec-fetch-mode", secMode)
	req.Header.Set("sec-fetch-dest", secDest)
	req.Header.Set("user-agent", userAgent)
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
}

func FSetDefaultChromeHeaders(req *fhttp.Request, userAgent, secSite, secMode, secDest string) {

	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`)
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-fetch-site", secSite)
	req.Header.Set("sec-fetch-mode", secMode)
	req.Header.Set("sec-fetch-dest", secDest)
	req.Header.Set("user-agent", userAgent)
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
}

func MapCountryCodeToCountryName(CountryCode string) (countryName string) {
	for _, val := range CountryCodeNameMap {
		if val.CountryCode == CountryCode {
			return val.Name
		}
	}
	log.Printf("Country code %s not in map", CountryCode)
	return CountryCode
}

func CheckNullPhone(phone string) string {
	if phone == "null" {
		return phone
	}
	return fmt.Sprintf(`"%s"`, phone)
}

func SomethingWentWrong(taskId string) {
	fmt.Printf("%s | %s | %s", ReturnFormattedTimestamp(), taskId, ColourRed("An unknown error occurred...\n"))
}
