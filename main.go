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
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"

	"github.com/stickmanventures/inquiry-bot/Godeps/_workspace/src/github.com/ereyes01/firebase"
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

// Find the request and pass it on to the publisher.
func Receive(event *firebase.StreamEvent) {
	if event.Event == "put" && event.Resource != nil {
		Publish(event.Resource.(map[string]interface{}))
	} else if event.Event == "patch" && event.Resource != nil {
		records := event.Resource.(map[string]interface{})

		for _, record := range records {
			if record == nil {
				continue
			}
			Publish(record.(map[string]interface{}))
		}
	}
}

// Push a request to slack as the slackbot.
func Publish(request map[string]interface{}) {
	if HasKey("email", request) && HasKey("name", request) && HasKey("phone", request) &&
		HasKey("referer", request) && HasKey("request", request) &&
		IsString(request["email"]) && IsString(request["name"]) && IsString(request["phone"]) &&
		IsString(request["referer"]) && IsString(request["request"]) {

		// Post to slack on another goroutine.
		go func() {
			Post(request["email"].(string), request["name"].(string), request["phone"].(string),
				request["referer"].(string), request["request"].(string))
		}()
	}
}

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

	api := new(firebase.Api)
	c := firebase.NewClient(*watch, *secret, *api)

	// Reconnect loop in the event of failure.
	for {
		stop := make(chan bool)
		events, err := c.Watch(nil, stop)
		if err != nil {
			log.Fatal(err)
		}

		stream(events, stop)
	}
}

// Stream over SSE with firebase until an error occurs (disconnect or whatnot).
func stream(events <-chan firebase.StreamEvent, stop chan bool) {
	for {
		// Wait for StreamEvents from the firebase SSE.
		event := <-events

		if event.Error != nil {
			close(stop)
			break
		}

		// The first StreamEvent is ignored since firebase does not allow
		// limitToLast below 1.
		if count == 0 {
			count = 1
			continue
		}

		count += 1
		Receive(&event)
	}
}
