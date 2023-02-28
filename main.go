package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kantapit123/fungjai-webhook/producer"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineMessage struct {
	Destination string `json:"destination"`
	Events      []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Timestamp  int64  `json:"timestamp"`
		Source     struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		Message struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

type ReplyMessageGo struct {
	ReplyToken string `json:"replyToken"`
	Messages   []Text `json:"messages"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ProFile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

type producedMessage struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

const (
	ChannelSecret = "0e3cfdba64d09eb82bbd1d2a9661c56d"
	ChannelToken  = "++mfeiTr/aYgS2psxjE7RIqrYfZ/GkoHs09NezY8gMp7WHo28rL9ij9A0d6/DkGBHMbpZ5zOBvZfaxJhyqFoxZr04iEVAhXdzbEzzpYzNDmN5BdSz8xwwfqZRVfT3Nn5Ug7DLTQAdagA0vvYi4XGJAdB04t89/1O/w1cDnyilFU="
)

func main() {
	// loads .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize linebot client
	client := &http.Client{}
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"), linebot.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Line bot client ERROR: ", err)
	}
	// Initialize kafka producer
	err = producer.InitKafka()
	if err != nil {
		log.Fatal("Kafka producer ERROR: ", err)
	}

	// Initilaze Echo web server
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/webhook", func(c echo.Context) error {

		// event := new(LineMessage)
		// if err := c.Bind(event); err != nil {
		// 	log.Println("err")
		// 	return c.String(http.StatusOK, "error")
		// }

		// log.Println(event)
		// fullname := getProfile(event.Events[0].Source.UserID)

		// if event.Events[0].Message.Text == "test" {
		// 	json_flex1 := []byte(`
		// 	{
		// 		"type": "bubble",
		// 		"header": {
		// 		  "type": "box",
		// 		  "layout": "vertical",
		// 		  "contents": [
		// 			{
		// 			  "type": "image",
		// 			  "url": "https://sv1.picz.in.th/images/2023/02/24/LMiinI.png",
		// 			  "size": "full",
		// 			  "aspectMode": "cover",
		// 			  "position": "relative"
		// 			}
		// 		  ],
		// 		  "background": {
		// 			"type": "linearGradient",
		// 			"angle": "0deg",
		// 			"startColor": "#E8E0D5",
		// 			"endColor": "#ffffff"
		// 		  },
		// 		  "position": "relative",
		// 		  "paddingAll": "0%"
		// 		},
		// 		"body": {
		// 		  "type": "box",
		// 		  "layout": "vertical",
		// 		  "contents": [
		// 			{
		// 			  "type": "box",
		// 			  "layout": "horizontal",
		// 			  "contents": [
		// 				{
		// 				  "type": "box",
		// 				  "layout": "vertical",
		// 				  "contents": [
		// 					{
		// 					  "type": "image",
		// 					  "url": "https://sv1.picz.in.th/images/2023/02/24/LMigQZ.png",
		// 					  "size": "full",
		// 					  "aspectMode": "cover",
		// 					  "aspectRatio": "21:16",
		// 					  "gravity": "center",
		// 					  "action": {
		// 						"type": "message",
		// 						"label": "action",
		// 						"text": "4"
		// 					  }
		// 					},
		// 					{
		// 					  "type": "image",
		// 					  "url": "https://sv1.picz.in.th/images/2023/02/24/LMiTmz.png",
		// 					  "size": "full",
		// 					  "aspectMode": "cover",
		// 					  "aspectRatio": "21:16",
		// 					  "gravity": "center",
		// 					  "action": {
		// 						"type": "message",
		// 						"label": "action",
		// 						"text": "2"
		// 					  }
		// 					}
		// 				  ],
		// 				  "flex": 1
		// 				},
		// 				{
		// 				  "type": "box",
		// 				  "layout": "vertical",
		// 				  "contents": [
		// 					{
		// 					  "type": "image",
		// 					  "url": "https://sv1.picz.in.th/images/2023/02/24/LMxaAa.png",
		// 					  "gravity": "center",
		// 					  "size": "full",
		// 					  "aspectRatio": "21:16",
		// 					  "aspectMode": "cover",
		// 					  "action": {
		// 						"type": "message",
		// 						"label": "action",
		// 						"text": "3"
		// 					  }
		// 					},
		// 					{
		// 					  "type": "image",
		// 					  "url": "https://sv1.picz.in.th/images/2023/02/24/LMi1bR.png",
		// 					  "gravity": "center",
		// 					  "size": "full",
		// 					  "aspectRatio": "21:16",
		// 					  "aspectMode": "cover",
		// 					  "action": {
		// 						"type": "message",
		// 						"label": "action",
		// 						"text": "1"
		// 					  }
		// 					}
		// 				  ]
		// 				}
		// 			  ]
		// 			}
		// 		  ],
		// 		  "paddingAll": "0px"
		// 		}
		// 	  }
		// 	`)

		// 	reply(event.Events[0].ReplyToken, json_flex1, bot)

		// } else {
		// 	text := Text{
		// 		Type: "text",
		// 		Text: fullname + " mood score is " + event.Events[0].Message.Text,
		// 	}

		// 	message := ReplyMessageGo{
		// 		ReplyToken: event.Events[0].ReplyToken,
		// 		Messages: []Text{
		// 			text,
		// 		},
		// 	}
		// 	replyMessageLine(message)
		// }

		events_bot, err := bot.ParseRequest(c.Request())
		if err != nil {
			// Do something when something bad happened.
		}
		topics := "user-messages"

		for _, event := range events_bot {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					messageJson, _ := json.Marshal(&producedMessage{
						Id:      message.ID,
						Message: message.Text,
					})
					producerErr := producer.Produce(topics, string(messageJson))
					if producerErr != nil {
						log.Print(err)
					} else {
						messageResponse := fmt.Sprintf("Produced [%s] successfully", message.Text)
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageResponse)).Do(); err != nil {
							log.Print(err)
						}
					}
				}
			}
		}

		log.Println("%% message success")
		return c.String(http.StatusOK, "ok!!")

	})

	e.Logger.Fatal(e.Start(":1323"))
}

func reply(replyToken string, jsonData []byte, bot *linebot.Client) error {

	container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
	if err != nil {
		log.Printf("failed to unmarshal Flex Message: %v", err)
	}

	if _, err = bot.ReplyMessage(
		replyToken,
		linebot.NewFlexMessage("You have new flex message", container),
	).Do(); err != nil {
		return err
	}
	return nil
}

func replyMessageLine(Message ReplyMessageGo) error {
	value, _ := json.Marshal(Message)

	url := "https://api.line.me/v2/bot/message/reply"

	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return err
}

func getProfile(userId string) string {

	url := "https://api.line.me/v2/bot/profile/" + userId

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var profile ProFile
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Println("%% err \n")
	}
	log.Println(profile.DisplayName)
	return profile.DisplayName

}
