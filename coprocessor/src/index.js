import express from "express";
import { Ratelimit } from "@upstash/ratelimit";
import { Redis } from "@upstash/redis";

/**
 * We are using a free plan on upstash.com
 */
const redis = new Redis({
  url: process.env['UPSTASH_REDIS_REST_URL'],
  token: process.env['UPSTASH_REDIS_REST_TOKEN'],
});

/**
 * This class only accepts an Upstash Redis instance,
 * but you can see how the logic applies below
 */
const rateLimit = new Ratelimit({
  redis: redis,
  limiter: Ratelimit.slidingWindow(3, "10 s"),
  analytics: true,
  enableProtection: true,
});

// This is currently a fixed value decided by Router
const CONTEXT_CLIENT_KEY = 'apollo_telemetry::client_name';
// This value can be customized to what value you want
const HEADER_NAME = 'apollographql-client-name';

const app = express();

/**
 * Get the header name to use to identify clients.
 * Default to using the same name as Apollo GraphOS telemetry,
 * but fallback to another header if needed.
 */
const getClientName = (payload) => {
  const contextName = payload?.context?.entries?.[CONTEXT_CLIENT_KEY];

  if (contextName) {
    return contextName;
  }

  const headers = payload.headers?.[HEADER_NAME] ?? [];
  const headerName = headers[0];
  if (headerName) {
    return headerName;
  }

  return undefined;
};

/**
 * Operation name to help identify other types of requests like Introspection
 */
const getOperationName = (payload) => {
  return payload?.context?.entries?.['operation_name'];
};

/**
 * Process and validate each request to apply a rate limit
 */
const processSupergraphRequestStage = async (payload) => {
  const clientName = getClientName(payload);
  const operationName = getOperationName(payload);

  // Do not fail or rate limit introspection queries
  if (operationName === 'IntrospectionQuery') {
    return payload;
  }

  // Make sure we have a client name to apply a rate limit to
  if (!clientName) {
    payload.control = { break: 400 };
    payload.body = {
      errors: [
        {
          message: "Invalid or missing client id",
          extensions: {
            code: "INVALID_CLIENT_ID",
          }
        }
      ]
    };

    return payload;
  }

  console.log(`Request from client '${clientName}'. Applying rate limit...`);

  // To apply a custom rate, you could change the logic to pass in the `cost` or wait for the `SupergraphResponse`
  // and decide how much rate to subtract instead of just a fixed amount.
  // See more details at: https://upstash.com/docs/redis/sdks/ratelimit-ts/methods
  const rateLimitResponse = await rateLimit.limit(clientName);

  if (!rateLimitResponse.success) {
    payload.control = { break: 429 };
    let errors = payload.body.errors || [];
    errors.push({
      message: `Rate limit exceeded for client: ${clientName}`,
      extensions: {
        code: "RATE_LIMIT_EXCEED",
        rateLimit: {
          clientName,
          limit: rateLimitResponse.limit,
          remaining: rateLimitResponse.remaining,
          reset: rateLimitResponse.reset
        }
      }
    });
    payload.body.errors = errors;

    return payload;
  }

  // Do other stuff in coprocessor if needed below...

  return payload;
};

/**
 * Set up the coprocessor server endpoint
 */
app.post("/", express.json(), async (req, res) => {
  const payload = req.body;

  let response = payload;
  switch (payload.stage) {
    case "SupergraphRequest":
      response = await processSupergraphRequestStage(payload);
      break;
  }

  res.send(response);
});

/**
 * Start the coprocessor server
 */
app.listen(3007, () => {
  console.log("🚀 Coprocessor running at http://localhost:3007");
});
