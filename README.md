inquiry-bot
-----------

inquiry-bot is a simple slack bot using webhooks to report business form
submissions to slack. The bot watches a firebase endpoint and looks for
specific properties before pushing to slack.

Building the Docker image
-------------------------

The easiest way to build the docker image is to use make. It is assumed you
have a working docker environment.
> make docker

Inquiry Bot requires that some environment variables be set when running in a
docker container. The required environment variables are:
```shell
INQUIRYBOT_FIREBASE='<FIREBASE_ENDPOINT>'
INQUIRYBOT_SECRET='<FIREBASE_SECRET>'
INQUIRYBOT_HOOK='<SLACKBOT_WEBHOOK_URL>'
INQUIRYBOT_CHANNEL='#channel'
```

Here's an example of creating a container with these environment variables:
```shell
docker create \
    -e INQUIRYBOT_FIREBASE='<FIREBASE_ENDPOINT>' \
    -e INQUIRYBOT_SECRET='<FIREBASE_SECRET>' \
    -e INQUIRYBOT_HOOK='<SLACKBOT_WEBHOOK_URL>' \
    -e INQUIRYBOT_CHANNEL='#channel' \
    stickmanventures.com/inquirybot
```

Deploying to GCP Compute Engine
-------------------------------

Inquiry Bot can be deployed to GCP using the container-vm image. After making
sure you are on the right project, a container vm can be created:
```shell
gcloud compute instances create slack-inquirybot \
    --image container-vm \
    --zone us-central1-f \
    --machine-type f1-micro
```

Now the image needs to be pushed to GCP, to do so the docker image needs to be
tagged under gcr.io:
```shell
docker tag stickmanventures.com/inquirybot gcr.io/<PROJECT_NAME>/inquirybot
```

Once tagged the image should be ready to push:

```shell
gcloud docker push gcr.io/<PROJECT_NAME>/inquirybot
```

From there you should be able to find your instance in the Google Cloud Console
and have a working docker environment with the image accessible.
