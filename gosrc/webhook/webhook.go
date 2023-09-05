package webhook

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/FergusJJ/coursework/gosrc/utils"
)

func SendWebhook(shoeName, productLink, productImage, siteName, proxyUrl, shoeSize, pid, price, webhookUrl, TaskID string) {

	timeStamp := fmt.Sprint(time.Now().Format("2006-01-02T15:04:05.000Z"))
	webhookBody := fmt.Sprintf(`
	  {
		"username":"Success",
		"timestamp": "%s",
		"embeds": [{
			"color": 5814783,
				"author": {
						"name": "%s",
						"url": "%s"
				},
				"thumbnail": {
						"url":"https://img.yeet.mx/proxy/?url=https://%s"
				  },
				"fields": [
				  {
						"name": "Site",
						"value": "%s",
						"inline": true
				},
				{
						"name": "Proxy",
						"value": "||%s||",
						"inline": true
				  },
				  {
						"name": "Size",
						"value": "%s",
						"inline": true
				},
				{
						"name": "PID",
						"value": "%s",
						"inline": true
				},
				{
						"name": "Price",
						"value": "%s",
						"inline": true
				}
			]
				
		  }
		]
	} 
	  `,
		timeStamp,
		shoeName,
		productLink,
		productImage,
		siteName,
		proxyUrl,
		shoeSize,
		pid,
		price,
	)
	bufBody := bytes.NewBuffer([]byte(webhookBody))
	resp, err := http.Post(webhookUrl, "application/json", bufBody)
	if err != nil {
		fmt.Printf("%s Error sending webhook, check the webhook url in settings...", utils.ReturnFormattedTimestamp())
	}
	_, _ = io.ReadAll(resp.Body)
	if resp.StatusCode < 205 {
		return
	} else {
		fmt.Printf("%s | %s | %s", utils.ReturnFormattedTimestamp(), TaskID, utils.ColourRed(fmt.Sprintf("Something went wrong sending webhook - HTTP Response %d\n", resp.StatusCode)))
	}

}
