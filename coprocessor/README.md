# Javascript Coprocessor

## Description

In `js-coprocessor/src/index.js`, the coprocessor is setup with `express` to listen to the `/` POST endpoint for various Router stages.

## Running the coprocessor

1. Setup a redis instance in the environment. We are using [@upstash/ratelimit-js](https://github.com/upstash/ratelimit-js).
   2.  You can also setup a free serverless Redis instance on https://console.upstash.com/
1. Run `npm install` to install dependencies
1. Run `npm run dev` to start the service
