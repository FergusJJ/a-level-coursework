package asos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	utils "github.com/FergusJJ/coursework/gosrc/utils"

	uuid "github.com/google/uuid"
	http "github.com/useflyent/fhttp"
)

func (s *localAsosSession) getProductPage() error {
	req, err := http.NewRequest(http.MethodGet, s.Profile.ProductURL, nil)
	if err != nil {
		fmt.Println(err)
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
			err = s.getProductPage()
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting the product page \n"))
			err = s.getProductPage()
			if err != nil {
				return err
			}

		}
		return nil
	} else {
		switch resp.StatusCode {
		case 403:
			//shouldn't happen, if it does it's due to headers, not proxies
			utils.SomethingWentWrong(s.Profile.TaskID)
			return fmt.Errorf("")
		default:
			if resp.Request.URL.String() != s.Profile.ProductURL {
				fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Invalid url - %s - HTTP Response %d\n", s.Profile.ProductURL, resp.StatusCode)))
				return fmt.Errorf("")
			}
			if resp.StatusCode < 400 {
				fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow("Got product page\n"))
			} else {
				utils.SomethingWentWrong(s.Profile.TaskID)
				return fmt.Errorf("")
			}

		}
		productPageBodyBytes, _ := io.ReadAll(resp.Body)
		productPageString := string(productPageBodyBytes)
		resp.Body.Close()
		if err = s.getATCBearer(); err != nil {
			return fmt.Errorf("")
		}
		fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow("Getting new bag...\n"))
		bagId, err := s.getBagIDForATC()
		if err != nil {
			return fmt.Errorf("")
		}
		bagATCLink := atcLink(bagId, s.AbckDevice.GetBrowser.Lang)
		var inStock bool
		for !s.IsCheckedOut {

			if s.VariantID == 0 {
				//first run set variantID, no looping
				s.VariantID, inStock = s.getVariantIdAndCheckStock(productPageString)
				//get variantID and check for stock
				if inStock && !s.IsCheckedOut {
					fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow(fmt.Sprintf("Adding item to bag - %d\n", s.VariantID)))
					//atc then  checkout
					err = s.checkoutItem(bagATCLink, bagId)
					if err != nil {
						//stop the task, don't monitor
						return fmt.Errorf("")
					}
					if s.IsCheckedOut {
						return nil
					}
				}
			}

			_, inStock = s.getVariantIdAndCheckStock(productPageString)
			if inStock && !s.IsCheckedOut {
				fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow(fmt.Sprintf("Adding item to bag - %d\n", s.VariantID)))
				err = s.checkoutItem(bagATCLink, bagId)
				if err != nil {
					//stop the task, don't monitor
					return fmt.Errorf("")
				}
				if s.IsCheckedOut {
					return nil
				}
			}
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow(fmt.Sprintf("Monitoring, out of stock - %d\n", s.VariantID)))
			utils.DelayRequest(s.Profile.Delay)
			req, err := http.NewRequest(http.MethodGet, s.Profile.ProductURL, nil)
			if err != nil {
				return fmt.Errorf("")
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
					err = s.getProductPage()
					if err != nil {
						return err
					}
				} else {
					fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst monitoring the product page \n"))
					err = s.getProductPage()
					if err != nil {
						return err
					}

				}
				return nil
			} else {
				switch resp.StatusCode {
				case 403:
					//shouldn't happen, if it does it's due to headers, not proxies
					utils.SomethingWentWrong(s.Profile.TaskID)
					return fmt.Errorf("")
				default:
					if resp.Request.URL.String() != s.Profile.ProductURL {
						fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Invalid url - %s - HTTP Response %d\n", s.Profile.ProductURL, resp.StatusCode)))
						return fmt.Errorf("")
					}
					if resp.StatusCode < 400 {
						//do nothing
					} else {
						utils.SomethingWentWrong(s.Profile.TaskID)
						return fmt.Errorf("")
					}
				}
				productPageBodyBytes, _ := io.ReadAll(resp.Body)
				productPageString = string(productPageBodyBytes)
				resp.Body.Close()
				//check stock
				//delay, then poll for stock again if not in stock => function will loop again, wthout getting the new variant as well
			}
		}
	}
	return nil
}

func (s *localAsosSession) getVariantIdAndCheckStock(responseBody string) (variantID int, inStock bool) {
	splitAtvariants := strings.Split(responseBody, "window.asos.pdp.config.product = ")
	splitAtvariants = strings.Split(splitAtvariants[1], ";") //will be the json data
	sizeSubstring := fmt.Sprintf(`,"size":"UK %s",`, s.Profile.Size)
	splitAtSizeId := strings.Split(splitAtvariants[0], sizeSubstring)
	variantIDList := strings.Split(splitAtSizeId[0], ":")
	variantID, err := strconv.Atoi(variantIDList[len(variantIDList)-1])
	if err != nil {
		fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Size out of available range - %s\n", s.Profile.Size)))
		//if the size is out of range variant id wont be found and will cause Atoi to fail as it won't be able to find a number
		fmt.Printf("[%s] | %s | size out of range\n", utils.ReturnFormattedTimestamp(), s.Profile.TaskID)
		//return some sort of size error
	}
	inStockSlice := strings.Split(splitAtSizeId[1], `,"`)
	inStockString := strings.Split(inStockSlice[0], `:`)
	inStock, _ = strconv.ParseBool(inStockString[1])
	if err != nil {
		log.Fatal(err)
	}
	return variantID, inStock
}

func (s *localAsosSession) checkoutItem(atcLink, bagId string) error {
	err := s.addVariantToCart(atcLink)
	if err != nil {
		return err
	}
	return nil
}

func (s *localAsosSession) addVariantToCart(bagATCLink string) error {
	requestBody := strings.NewReader(fmt.Sprintf(`{"variantId":%d}`, s.VariantID))
	req, err := http.NewRequest(http.MethodPost, bagATCLink, requestBody)
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return err
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`)
	req.Header.Set("asos-c-plat", "Web")
	req.Header.Set("asos-bag-origin", "EUN")
	req.Header.Set("authorization", "Bearer "+s.Bearer.Access_token)
	req.Header.Set("asos-c-name", "Asos.Commerce.Bag.Sdk")
	req.Header.Set("asos-c-store", "COM")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", s.UserAgent)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("asos-c-ismobile", "false")
	req.Header.Set("asos-c-ver", "5.5.166")
	req.Header.Set("asos-c-istablet", "false")
	req.Header.Set("asos-cid", s.Cid)
	req.Header.Set("origin", "https://www.asos.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", s.Profile.ProductURL)
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	utils.DelayRequest(s.Profile.Delay)
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return err
			}
			err = s.addVariantToCart(bagATCLink)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst adding to cart \n"))
			err = s.addVariantToCart(bagATCLink)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		//item successfull added to cart
		if !strings.Contains(string(bodyBytes), `"totalQuantity": 0`) {
			s.completeCheckout()
		} else if strings.Contains(string(bodyBytes), `"errorCode": "OutOfAdditionalStock"`) {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow(fmt.Sprintf("Product out of stock - %d - %s - monitoring for restocks...\n", s.VariantID, s.Profile.Size)))
		} else {
			utils.SomethingWentWrong(s.Profile.TaskID)
			return fmt.Errorf("")
		}
		return nil
	}
}

func atcLink(bagID, lang string) string {
	bagATCLink := fmt.Sprintf("https://www.asos.com/api/commerce/bag/v4/bags/%s/product?expand=summary,total&lang=%s", bagID, lang) // https://www.asos.com/api/commerce/bag/v4/bags/50e0eb66-e0aa-447c-aebb-186d64aff238/product?expand=summary,total&lang=en-GB
	return bagATCLink
}

func (s *localAsosSession) completeCheckout() {
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow("Getting delivery address...\n"))
	addressRequestUrl := fmt.Sprintf("https://secure.asos.com/api/commerce/bag/v4/bags/%s/checkout?expand=delivery,total,address,discount,deliveryOptions,spendLevelDiscount", s.BagId)
	addressId, err := s.getAddressId(addressRequestUrl)
	if err != nil {
		fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Something went wrong getting delivery address - %s...\n", s.Profile.Email)))
		return
		//return task, unknown error
	}
	cardRequestUrl := fmt.Sprintf("https://secure.asos.com/api/customer/paymentdetails/v2/customers/%s/paymentdetails", s.CustomerId)
	cardId, err := s.getCardId(cardRequestUrl)
	if err != nil {
		fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed(fmt.Sprintf("Something went wrong getting card id - %s...\n", s.Profile.Email)))
		return
	}
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow("Getting payment token...\n"))
	securityCodeToken, err := s.getPaymentToken()
	if err != nil {
		fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Something went wrong getting payment token...\n"))
		return
	}
	cartDataUrl := fmt.Sprintf("https://secure.asos.com/api/commerce/bag/v4/bags/%s/deliveryaddress?expand=delivery,total,address,discount,deliveryOptions,spendLevelDiscount", s.BagId)
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow("Getting cart details...\n"))
	deliveryDate, deliveryId, amount, err := s.getCartDetails(cartDataUrl, addressId)
	if err != nil {
		fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Something went wrong getting cart details...\n"))
		return
	}
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow("Submitting billing...\n"))
	paymentRef, err := s.submitBilling(addressId, cardId, securityCodeToken, deliveryDate, deliveryId, amount)
	if err != nil {
		return
	}
	fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow("Submitting order...\n"))
	err = s.submitOrder(paymentRef, amount)
	if err != nil {
		return
	}

	//fmt.Println(cardId, securityCodeToken, deliveryDate, deliveryId, amount)
}

func (s *localAsosSession) getAddressId(requestURL string) (addressId string, err error) {
	req, err := http.NewRequest(http.MethodPut, requestURL, nil)
	if err != nil {
		return "", err
	}
	utils.SetDefaultChromeHeaders(req, s.UserAgent, "same-site", "cors", "empty")
	req.Header.Set("authorization", "Bearer "+s.Bearer.Access_token)
	req.Header.Set("asos-c-name", "asos.commerce.checkout.web")
	req.Header.Set("asos-bag-origin", "EUN")
	req.Header.Set("accept", "*/*")
	req.Header.Set("asos-c-plat", "Web")
	req.Header.Set("asos-c-ver", "2.0.1.6645")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("asos-cid", s.Cid)
	req.Header.Set("origin", "https://secure.asos.com")
	req.Header.Set("referer", "https://secure.asos.com/")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept-encoding", "gzip, deflate")
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return "", err
			}
			addressId, err = s.getAddressId(requestURL)
			if err != nil {
				return "", err
			}
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting address Id \n"))
			addressId, err = s.getAddressId(requestURL)
			if err != nil {
				return "", err
			}

		}
		return addressId, nil
	} else {
		defer resp.Body.Close()
		if 199 < resp.StatusCode && resp.StatusCode < 203 {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			s.CustomerId = strings.Split(strings.Split(string(b), `customerId": `)[1], ",")[0]
			addressId = strings.Split(strings.Split(string(b), `"customerAddressId": `)[1], ",")[0]
			s.Profile.Phone = "null"
			if !strings.Contains(string(b), `"telephoneMobile":null`) {
				innerPhone := strings.Split(string(b), `"telephoneMobile": "`)[1]
				s.Profile.Phone = strings.Split(innerPhone, `"`)[0]
			}
			//get webhook info
			s.WebhookInfo.ItemName = strings.TrimSpace(strings.Split(strings.Split(string(b), `"name": "`)[1], `"`)[0])
			s.WebhookInfo.Price = strings.TrimSpace(strings.Split(strings.Split(string(b), `"text": "`)[1], `"`)[0])
			s.WebhookInfo.Pid = strings.TrimSpace(strings.Split(strings.Split(string(b), `"productId":`)[1], `,`)[0])
			s.WebhookInfo.Size = s.Profile.Size
			s.WebhookInfo.ImageLink = strings.TrimSpace(strings.Split(strings.Split(string(b), `"url": "`)[1], `"`)[0])
			return addressId, nil
		} else if resp.StatusCode == 400 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow(fmt.Sprintf("Product out of stock - %s - %s - monitoring for restocks...\n", s.WebhookInfo.Pid, s.Profile.Size)))
			return "", fmt.Errorf("")
		} else if resp.StatusCode == 419 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Session expired\n"))
			return "", fmt.Errorf("")
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen(fmt.Sprintf("An unknown error has occurred - HTTP Response %d\n", resp.StatusCode)))
			return "", fmt.Errorf("")
		}
	}
}

func (s *localAsosSession) getCardId(requestURL string) (cardId string, err error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return "", err
	}
	utils.SetDefaultChromeHeaders(req, s.UserAgent, "same-site", "cors", "empty")
	req.Header.Set("authorization", "Bearer "+s.Bearer.Access_token)
	req.Header.Set("asos-c-name", "asos.commerce.checkout.web")
	req.Header.Set("asos-bag-origin", "EUN")
	req.Header.Set("accept", "*/*")
	req.Header.Set("asos-c-plat", "Web")
	req.Header.Set("asos-c-ver", "2.0.1.6645")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("asos-cid", s.Cid)
	req.Header.Set("origin", "https://secure.asos.com")
	req.Header.Set("referer", "https://secure.asos.com/")
	req.Header.Set("accept-encoding", "gzip, deflate")
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return "", err
			}
			cardId, err := s.getCardId(requestURL)
			if err != nil {
				return "", err
			}
			return cardId, nil
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting card Id \n"))
			cardId, err := s.getCardId(requestURL)
			if err != nil {
				return "", err
			}
			return cardId, nil

		}
	} else {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		cardId = strings.Split(strings.Split(string(b), `"id":"`)[1], `",`)[0]
		defer resp.Body.Close()
		return cardId, nil
	}
}

func (s *localAsosSession) getPaymentToken() (securityToken string, err error) {
	requestUrl := "https://api.asos.com/finance/payments-card/v5/card/tokens"
	requestBody := []byte(fmt.Sprintf(`{"securityCode":"%s"}`, s.Profile.CC.CVC))
	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	utils.SetDefaultChromeHeaders(req, s.UserAgent, "same-site", "cors", "empty")
	req.Header.Set("asos-c-name", "asos-payments-card-iframe")
	req.Header.Set("asos-c-plat", "payments")
	req.Header.Set("asos-c-ver", "20200930.1")
	req.Header.Set("asos-cid", s.Cid)
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("origin", "https://payments.asos.com")
	req.Header.Set("referer", "https://payments.asos.com/")
	req.Header.Set("accept-encoding", "gzip, deflate")
	req.Header.Set("origin", "https://payments.asos.com")
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return "", err
			}
			securityToken, err := s.getPaymentToken()
			if err != nil {
				return "", err
			}
			return securityToken, nil
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting payment token \n"))
			securityToken, err := s.getPaymentToken()
			if err != nil {
				return "", err
			}
			return securityToken, nil

		}
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		securityToken = strings.Split(strings.Split(string(bodyBytes), `"securityCodeToken":"`)[1], `"`)[0]
		return securityToken, nil
	}

}

func (s *localAsosSession) getCartDetails(requestUrl, customerAddressId string) (deliveryDate string, deliveryMethodId int, amount float64, err error) {
	requestBody := []byte(fmt.Sprintf(`{"addressLine1":"%s","addressLine2":"%s","addressLine3":null,"country":"%s","countryName":"%s","countyStateProvinceOrArea":null,"customerAddressId":%s,"emailAddress":"%s","firstName":"%s","lastName":"%s","locality":"%s","postalCode":"%s","telephoneDaytime":%s,"telephoneEvening":%s,"telephoneMobile":%s}`,
		s.Profile.AddressLine1,
		s.Profile.AddressLine2,
		strings.ToUpper(s.Profile.CountryCode),
		utils.MapCountryCodeToCountryName(strings.ToUpper(s.Profile.CountryCode)),
		customerAddressId,
		s.Profile.Email,
		s.Profile.FirstName,
		s.Profile.LastName,
		s.Profile.City,
		s.Profile.Postcode,
		utils.CheckNullPhone(s.Profile.Phone),
		utils.CheckNullPhone(s.Profile.Phone),
		utils.CheckNullPhone(s.Profile.Phone),
	))
	req, err := http.NewRequest(http.MethodPut, requestUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", 0, 0.0, err
	}
	utils.SetDefaultChromeHeaders(req, s.UserAgent, "same-site", "cors", "empty")
	req.Header.Set("authorization", "Bearer "+s.Bearer.Access_token)
	req.Header.Set("asos-c-name", "asos.commerce.checkout.web")
	req.Header.Set("asos-bag-origin", "EUN")
	req.Header.Set("accept", "*/*")
	req.Header.Set("asos-c-plat", "Web")
	req.Header.Set("asos-c-ver", "2.0.1.6645")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("asos-cid", s.Cid)
	req.Header.Set("origin", "https://secure.asos.com")
	req.Header.Set("referer", "https://secure.asos.com/")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept-encoding", "gzip, deflate")
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return "", 0, 0.0, err
			}
			deliveryDate, deliveryMethodId, amount, err := s.getCartDetails(requestUrl, customerAddressId)
			if err != nil {
				return "", 0, 0.0, err
			}
			return deliveryDate, deliveryMethodId, amount, nil
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst getting cart details \n"))
			deliveryDate, deliveryMethodId, amount, err := s.getCartDetails(requestUrl, customerAddressId)
			if err != nil {
				return "", 0, 0.0, err
			}
			return deliveryDate, deliveryMethodId, amount, nil
		}
	} else {
		defer resp.Body.Close()
		if 199 < resp.StatusCode && resp.StatusCode < 203 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen("Submitted billing\n"))
		} else if resp.StatusCode == 400 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow(fmt.Sprintf("Product out of stock - %s - %s - monitoring for restocks...\n", s.WebhookInfo.Pid, s.Profile.Size)))
			return "", 0, 0.0, fmt.Errorf("")
		} else if resp.StatusCode == 419 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Session expired\n"))
			return "", 0, 0.0, fmt.Errorf("")
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen(fmt.Sprintf("An unknown error has occurred - HTTP Response %d\n", resp.StatusCode)))
			return "", 0, 0.0, fmt.Errorf("")
		}
		body, _ := ioutil.ReadAll(resp.Body) // response body is []byte
		var cartDataResponse utils.CartValues
		if err := json.Unmarshal(body, &cartDataResponse); err != nil {
			if err != nil {
				return "", 0, 0.0, err
			}
		}
		deliveryDate = cartDataResponse.Bag.Delivery.Options[0].EstimatedDeliveryDate
		deliveryMethodId = cartDataResponse.Bag.Delivery.Options[0].DeliveryMethodID
		amount = cartDataResponse.Bag.Total.Total.Value
		return deliveryDate, deliveryMethodId, amount, nil
	}
}

func (s *localAsosSession) submitBilling(customerAddressId, cardId, securityCodeToken, deliveryDate string, deliveryId int, amount float64) (string, error) {
	paymentReference := uuid.NewString()
	paymentLink := fmt.Sprintf("https://secure.asos.com/api/finance/payments-card/v5/card-on-file/payments/%s", paymentReference)
	requestBody := fmt.Sprintf(`{"billingAddress":{"address1":"%s","address2":"%s","country":"%s","countyStateProvinceOrArea":"%s","countyStateProvinceOrAreaCode":null,"firstName":"%s","lastName":"%s","locality":"%s","postalCode":"%s","telephoneMobile":%s},"card":{"cardId":"%s","token":"%s"},"delivery":{"deliveryAddress":{"address1":"%s","address2":"%s","country":"%s","countyStateProvinceOrArea":"%s","countyStateProvinceOrAreaCode":null,"firstName":"%s","lastName":"%s","locality":"%s","postalCode":"%s","telephoneMobile":%s},"deliveryDate":"%s","deliveryMethodId":%d},"platform":{"challengeWindowFullscreen":false,"deviceInfo":{"colorDepth":24,"javaEnabled":false,"language":"%s","screenHeight":%d,"screenWidth":%d,"timeZoneOffset":%d},"name":"web"},"transaction":`,
		s.Profile.AddressLine1,
		s.Profile.AddressLine2,
		strings.ToUpper(s.Profile.CountryCode),
		s.Profile.Province,
		s.Profile.FirstName,
		s.Profile.LastName,
		s.Profile.City,
		s.Profile.Postcode,
		utils.CheckNullPhone(s.Profile.Phone),
		cardId, securityCodeToken,
		s.Profile.AddressLine1,
		s.Profile.AddressLine2,
		strings.ToUpper(s.Profile.CountryCode),
		s.Profile.Province,
		s.Profile.FirstName,
		s.Profile.LastName,
		s.Profile.City,
		s.Profile.Postcode,
		utils.CheckNullPhone(s.Profile.Phone),
		deliveryDate,
		deliveryId,
		s.AbckDevice.GetBrowser.Lang,
		s.AbckDevice.Gd.Height,
		s.AbckDevice.Gd.Width,
		s.AbckDevice.Fingerprint.TimzoneOffset,
	)
	//get rid of unnecessary 0s
	if amount == float64(int64(amount)) {
		requestBody = requestBody + fmt.Sprintf(`{"amount":%d,"currency":"GBP"}}`, int(amount))
	} else {
		requestBody = requestBody + fmt.Sprintf(`{"amount":%.2f,"currency":"GBP"}}`, amount)
	}
	req, err := http.NewRequest(http.MethodPut, paymentLink, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		utils.SomethingWentWrong(s.Profile.TaskID)
		return "", err
	}
	utils.SetDefaultChromeHeaders(req, s.UserAgent, "same-origin", "cors", "empty")
	req.Header.Set("authorization", "Bearer "+s.Bearer.Access_token)
	req.Header.Set("asos-c-name", "asos.commerce.checkout.web")
	req.Header.Set("asos-bag-origin", "EUN")
	req.Header.Set("accept", "*/*")
	req.Header.Set("asos-c-plat", "Web")
	req.Header.Set("asos-c-ver", "2.0.1.6645")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("asos-cid", s.Cid)
	req.Header.Set("origin", "https://secure.asos.com")
	req.Header.Set("referer", "https://secure.asos.com/")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept-encoding", "gzip, deflate")
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return "", err
			}
			paymentReference, err := s.submitBilling(customerAddressId, cardId, securityCodeToken, deliveryDate, deliveryId, amount)
			if err != nil {
				return "", err
			}
			return paymentReference, nil
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst submitting billing \n"))
			paymentReference, err := s.submitBilling(customerAddressId, cardId, securityCodeToken, deliveryDate, deliveryId, amount)
			if err != nil {
				return "", err
			}
			return paymentReference, nil
		}
	} else {
		defer resp.Body.Close()
		if 199 < resp.StatusCode && resp.StatusCode < 203 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen("Submitted billing\n"))
		} else if resp.StatusCode == 400 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow(fmt.Sprintf("Product out of stock - %s - %s - monitoring for restocks...\n", s.WebhookInfo.Pid, s.Profile.Size)))
			return "", fmt.Errorf("")
		} else if resp.StatusCode == 419 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Session expired\n"))
			return "", fmt.Errorf("")
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen(fmt.Sprintf("An unknown error has occurred - HTTP Response %d\n", resp.StatusCode)))
			return "", fmt.Errorf("")
		}
		//check response status
		return paymentReference, nil
	}
}

func (s *localAsosSession) submitOrder(paymentRef string, amount float64) error {
	paymentLink := "https://secure.asos.com/api/commerce/order/v3/orders/createorder?lang=en-GB&expand=bag,customer"
	var requestBody string
	if amount == float64(int64(amount)) {
		requestBody = fmt.Sprintf(`{"bagId":"%s","concurrencyKey":"%d","customer":{"address":{"address1":"%s","address2":"%s","countryCode":"%s","countyStateProvinceOrArea":"%s","countyStateProvinceOrAreaCode":null,"locality":"%s","postalCode":"%s"},"customerGuid":null,"customerId":%s,"firstName":"%s","isFirstTimeBuyer":false,"lastName":"%s"},"guaranteeStockAllocation":false,"paymentReference":"%s"}`,
			s.BagId, int(amount), s.Profile.AddressLine1, s.Profile.AddressLine2, strings.ToUpper(s.Profile.CountryCode), s.Profile.Province, s.Profile.City, s.Profile.Postcode, s.CustomerId, s.Profile.FirstName, s.Profile.LastName, paymentRef)
	} else {
		requestBody = fmt.Sprintf(`{"bagId":"%s","concurrencyKey":"%.2f","customer":{"address":{"address1":"%s","address2":"%s","countryCode":"%s","countyStateProvinceOrArea":"%s","countyStateProvinceOrAreaCode":null,"locality":"%s","postalCode":"%s"},"customerGuid":null,"customerId":%s,"firstName":"%s","isFirstTimeBuyer":false,"lastName":"%s"},"guaranteeStockAllocation":false,"paymentReference":"%s"}`,
			s.BagId, amount, s.Profile.AddressLine1, s.Profile.AddressLine2, strings.ToUpper(s.Profile.CountryCode), s.Profile.Province, s.Profile.City, s.Profile.Postcode, s.CustomerId, s.Profile.FirstName, s.Profile.LastName, paymentRef)
	}
	req, err := http.NewRequest(http.MethodPost, paymentLink, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	utils.SetDefaultChromeHeaders(req, s.UserAgent, "same-origin", "cors", "empty")
	req.Header.Set("authorization", "Bearer "+s.Bearer.Access_token)
	req.Header.Set("asos-c-name", "asos.commerce.checkout.web")
	req.Header.Set("asos-bag-origin", "EUN")
	req.Header.Set("accept", "*/*")
	req.Header.Set("asos-c-plat", "Web")
	req.Header.Set("asos-c-ver", "2.0.1.6645")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("asos-cid", s.Cid)
	req.Header.Set("origin", "https://secure.asos.com")
	req.Header.Set("referer", "https://secure.asos.com/")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept-encoding", "gzip, deflate")
	resp, err := s.AsosClient.Do(req)
	if err != nil {
		if s.Profile.Proxy != "localhost" {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Proxy error, switching...\n"))
			err = s.switchProxy()
			if err != nil {
				return err
			}
			err := s.submitOrder(paymentRef, amount)
			if err != nil {
				return err
			}
			return nil
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Connection error occurred whilst submitting order \n"))
			err := s.submitOrder(paymentRef, amount)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		defer resp.Body.Close()
		if 199 < resp.StatusCode && resp.StatusCode < 203 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen(fmt.Sprintf("Successful checkout - %s - %s\n", s.WebhookInfo.Pid, s.Profile.Size)))
			s.IsCheckedOut = true
			return nil
		} else if resp.StatusCode == 400 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourYellow(fmt.Sprintf("Product out of stock - %s - %s - monitoring for restocks...\n", s.WebhookInfo.Pid, s.Profile.Size)))
			return nil
		} else if resp.StatusCode == 419 {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourRed("Session expired\n"))
			return fmt.Errorf("")
		} else {
			fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), s.Profile.TaskID, utils.ColourGreen(fmt.Sprintf("An unknown error has occurred - HTTP Response %d\n", resp.StatusCode)))
			return fmt.Errorf("")
		}
	}
	//check response status

}
