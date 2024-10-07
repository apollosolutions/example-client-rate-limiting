# Coprocessor "Hello World"

This repository demonstrates how to setup a skeleton coprocessor that simply logs the `payload` from the Router.

## Running the Example

> Note: To run this example, you will need a GraphOS Enterprise plan and must create `/router/.env` based on `/router/.env.example` which exports `APOLLO_KEY` and `APOLLO_GRAPH_REF`.

1. Run the subgraph from the `/subgraph` directory with `npm run dev`
1. Run the coprocessor based on your language of choice by following the README from the appropriate `/*-coprocessor` directory ([javascript](./js-coprocessor/README.md), [Java](./java-coprocessor/README.md), [golang](./golang-coprocessor/README.md)).
1. In the `/router` directory, download the router by running `./download_router.sh`
1. In the `/router` directory, compose the schema by running `./create_local_schema.sh`
1. In the `/router` directory, run the router by running `./start_router.sh`

Now if you run this code in the browser (http://127.0.0.1:4000/), you will be able to query the router and you will see the `payload` logged in the terminal by the coprocessor.

## Code Highlights

### Coprocessor Configuration

In `router/router-config.yaml`, the coprocessor is configured with the Router to be called on the `router` `request` stage.

### Coprocessor
