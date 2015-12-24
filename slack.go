// Copyright 2016 Stickman Ventures
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
