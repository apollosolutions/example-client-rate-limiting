# Javascript Coprocessor

## Description

In `js-coprocessor/src/index.js`, the coprocessor is setup with `express` to listen to the `/` POST endpoint and respond to the `RouterRequest` stage.

In the `processRouterRequestStage` function, the payload is logged.

## Running the coprocessor

1. Run `npm install` to install dependencies
1. Run `npm run dev` to start the service
