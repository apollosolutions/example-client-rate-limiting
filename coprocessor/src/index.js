import express from "express";

import { Ratelimit } from "@upstash/ratelimit";
import { Redis } from "@upstash/redis";

const ratelimit = new Ratelimit({
  redis: Redis.fromEnv(),
  limiter: Ratelimit.slidingWindow(10, "10 s"),
  analytics: true,
  enableProtection: true
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

const getOperationName = (payload) => {
  return payload?.context?.entries?.['operation_name'];
};

const processSupergraphRequestStage = async (payload) => {
  console.log(payload);

  const clientName = getClientName(payload);
  const operationName = getOperationName(payload);
  if (!clientName && operationName !== 'IntrospectionQuery') {
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

  const { success } = await ratelimit.limit(clientName);
  if (!success) {
    payload.control = { break: 429 };
    payload.body = {
      errors: [
        {
          message: `Rate limit exceeded for client ${clientName}`,
          extensions: {
            code: "RATE_LIMIT_EXCEED",
          }
        }
      ]
    };

    return payload;
  }

  return payload;
};

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

app.listen(3007, () => {
  console.log("🚀 Server running at http://localhost:3007");
});
