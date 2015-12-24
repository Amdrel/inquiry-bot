package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Add some personality to inquiry_bot.
var quotes = []string{
	"It appears someone wants to do business with us.",
	"The word is getting out, we got a potential client.",
	"Someone with good tastes wants to do some business.",
	"A new proposition has come in.",
	"Some new business came in.",
	"A new client, hopefully not another peasant.",
}

// Post the inquiry details to the channel.
func Post(email string, name string, phone string, referer string, request string) {
	message := fmt.Sprintf(""+
		"%s\n\n"+

		"Email: %s\n"+
		"Name: %s\n"+
		"Phone: %s\n"+
		"Referer: %s\n"+
		"Request: %s\n",

		quotes[RandRange(0, len(quotes))], email, name, phone, referer, request)

	jsonPayload, err := json.Marshal(map[string]interface{}{
		"text":    message,
		"channel": *channel,
	})
	if err != nil {
		log.Print(err)
		return
	}

	_, err = http.PostForm(*hook, url.Values{
		"payload": {string(jsonPayload)},
	})
	if err != nil {
		log.Print(err)
		return
	}
}
