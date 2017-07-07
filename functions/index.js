// Copyright 2016 Stickman Ventures
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const functions = require('firebase-functions');
const IncomingWebhook = require('@slack/client').IncomingWebhook;

// Slack webhook client used to post messages to a channel.
const webhook = new IncomingWebhook(functions.config().slack.webhook);

// A list of quotes that are randomly chosen from when an inquiry comes through.
const quotes = [
  "It appears someone wants to do business with us.",
  "The word is getting out, we got a potential client.",
  "Someone with good tastes wants to do some business.",
  "A new proposition has come in.",
  "Some new business came in.",
  "A new client, hopefully not another peasant.",
];

// Posts a comment to a slack channel with a random quote and contact
// information when a request is submitted to the firebase database.
exports.postInquiry = functions.database.ref('/requests/{requestId}').onWrite((event) => {
  const data = event.data.val();
  const quote = quotes[Math.floor(Math.random() * quotes.length)];
  const message = `${quote}

Email: ${data.email}
Name: ${data.name}
Phone: ${data.phone}
Referer: ${data.referer}
Request: ${data.request}`;

  return new Promise((resolve, reject) => {
    webhook.send(message, (err, header, statusCode, body) => {
      if (err) {
        reject(err);
        return;
      }
      resolve();
    });
  });
});
