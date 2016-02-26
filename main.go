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
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"

	"github.com/ereyes01/firebase"
)

// Optional secret passed from the command line passed to the firebase client.
var watch = flag.String("firebase", "", "firebase url to watch")

// Optional secret passed from the command line passed to the firebase client.
var secret = flag.String("secret", "", "secret for firebase authentication")

// Optional secret passed from the command line passed to the firebase client.
var hook = flag.String("hook", "", "slack webhook to use")

// Required argument for which channel the slackbot should post to.
var channel = flag.String("channel", "", "channel to publish inquiries to")

// Optional cpu-profile flag for performance debugging.
var cpuprofile = flag.String("cpu-profile", "", "where to write cpu profile")

// Counter for the amount of events received. Used to ignore the first event
// firebase always sends; so that the slackbot does not duplicate inquiries.
var count uint

// Flag to determine if the cpu profiler is enabled.
var profiling bool

func main() {
	// Query the flag values from the environment before parsing the flags so
	// the flags take precedence if specified.
	*secret = os.Getenv("INQUIRYBOT_SECRET")
	*watch = os.Getenv("INQUIRYBOT_FIREBASE")
	*hook = os.Getenv("INQUIRYBOT_HOOK")
	*channel = os.Getenv("INQUIRYBOT_CHANNEL")

	// Overwrite flag values and merge with environment variables.
	flag.Parse()

	// Check required arguments.
	if *watch == "" || *channel == "" || *hook == "" {
		PrintUsage()
	}

	if *cpuprofile != "" {
		path, err := filepath.Abs(*cpuprofile)
		if err != nil {
			log.Fatalf("fatal: unable to get absolute path of: %s", *cpuprofile)
		}
		log.Printf("debug: profiling, outputting to: %s", path)

		// Remove any previously saved profile just in case.
		if _, err := os.Stat(path); err == nil {
			err = os.Remove(path)
			if err != nil {
				log.Fatalf("fatal: unable to unlink: %s", path)
			}
		}

		// Start CPU profiling.
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("fatal: unable to create profile at: %s", path)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatalf("fatal: unable to start cpu profiling: %s", err.Error())
		}
		profiling = true
	}

	// If the CPU profiler is running, flush the profile on SIGINT and exit
	// with a clean return code.
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		<-sigchan
		if profiling {
			log.Printf("debug: flushing cpu profile")
			pprof.StopCPUProfile()
		}
		os.Exit(0)
	}()

	stop := make(chan bool)
	c := firebase.NewClient(*watch, *secret, nil)
	events, err := c.Watch(requestParser, stop)
	if err != nil {
		log.Fatal(err)
	}

	// Main event loop, all valid requests from the watch are posted to slack
	// using the slack webhook.
	for event := range events {
		if event.Error != nil {
			log.Println("Stream error: ", event.Error)
			continue
		}

		if event.UnmarshallerError != nil {
			log.Println("Malformed event: ", event.UnmarshallerError)
			continue
		}

		// Skip the first event so we don't replay past requests.
		if count == 0 {
			count += 1
			continue
		}

		// Ignore nil events that can come through due to deletes.
		if event.Resource == nil {
			continue
		}

		// Post all the requests in the event to slack.
		requests := event.Resource.([]interface{})
		for _, request := range requests {
			go Post(request.(map[string]interface{}))
		}

		count += 1
	}
}

func requestParser(path string, data []byte) (interface{}, error) {
	var wrapper interface{}
	err := json.Unmarshal(data, &wrapper)
	if err != nil {
		return nil, err
	}

	if wrapper == nil {
		return nil, err
	}

	m := wrapper.(map[string]interface{})
	var requests []interface{}

	// Loop through the response. Firebase can send this down with 2 different
	// schemas that we have to check for, do not ask my why.
	for _, v := range m {
		switch vv := v.(type) {
		case map[string]interface{}:
			// If the children are maps then this is a bundle of requests.
			requests = append(requests, vv)
		case string:
			// If we're looping through strings this is a single event. Pass the
			// wrapper as it contains the request fields and return immediately
			// to prevent duplicate requests being pushed.
			requests = append(requests, m)
			return requests, err
		}
	}

	return requests, err
}
