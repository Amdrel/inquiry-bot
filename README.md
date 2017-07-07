# Inquiry Bot

Inquiry Bot is a simple slack bot using webhooks to report business form
submissions to slack. Data is stored in firebase and sent to slack after
being invoked by a [Firebase Realtime
Trigger](https://firebase.google.com/docs/functions/database-events)

## Deploying to Firebase

After cloning the repo, install dependencies under the [functions](functions)
directory. We have a `yarn.lock` file for yarn, but you can use npm if you
want to as well.

Next you need to add a config value for your slack webhook.

```firebase functions:config:set slack.webhook=https://hooks.slack.com/services/WEBHOOK_STRING```

Once complete you can deploy the project just like any other firebase project.

```firebase deploy --only functions```

And that's all there is to it!

## License

Copyright 2016 Stickman Ventures

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
