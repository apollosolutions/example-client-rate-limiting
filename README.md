# Rate Limiting Coprocessor

This repository demonstrates how to set up a coprocessor that can rate limit requests coming from clients, connecting to an external store to track requests over time.

**The code in this repository is experimental and has been provided for reference purposes only. Community feedback is welcome but this project may not be supported in the same way that repositories in the official [Apollo GraphQL GitHub organization](https://github.com/apollographql) are. If you need help you can file an issue on this repository, [contact Apollo](https://www.apollographql.com/contact-sales) to talk to an expert, or create a ticket directly in Apollo Studio.**

## Running the Example

> Note: To run this example, you will need a GraphOS Enterprise plan and must create `/router/.env` based on `/router/.env.example` which exports `APOLLO_KEY` and `APOLLO_GRAPH_REF`.
>
> You will also need a free Upstash account to set up a Redis instance and create `/coprocessor/.env` based on `/coprocessor/.env.example` which exports `UPSTASH_REDIS_REST_URL` and `UPSTASH_REDIS_REST_TOKEN`

1. Run the subgraph from the `/subgraph` directory with `npm run dev`
2. Run the coprocessor in the `/coprocessor` directory with `npm run dev`
1. In the `/router` directory, download the router by running `./download_router.sh`
1. In the `/router` directory, compose the schema by running `./create_local_schema.sh`
1. In the `/router` directory, run the router by running `./start_router.sh`

Now if you run this code in the browser (http://127.0.0.1:4000/), you will be able to query the router and you will see the `payload` logged in the terminal by the coprocessor.

## Code Highlights

### Coprocessor Configuration

In `router/router-config.yaml`, the coprocessor is configured with the Router to be called on the `SupergraphRequest` stage.

### Coprocessor
We use the `@upstash/redis` library to connect to a free hosted Redis instance, and then use `@upstash/ratelimit` to save the request volume per client and apply a rate limit if exceeded. This also has examples on how to use the new [demand control](https://www.apollographql.com/docs/graphos/routing/security/demand-control) feature to apply a rate limit after the `@cost` is calculated in the `SupergraphResponse`.
