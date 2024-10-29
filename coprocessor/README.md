# Javascript Coprocessor

## Description

In `coprocessor/src/index.js`, the coprocessor is set up with `express` to listen to the `/` POST endpoint for various Router stages.

## Running the coprocessor

- Setup a redis instance in the environment. We are using the rate limiting library [@upstash/ratelimit-js](https://github.com/upstash/ratelimit-js).
  - You can also set up a free serverless Redis instance on https://console.upstash.com/
- Run `npm install` to install dependencies
- Run `npm run dev` to start the service
