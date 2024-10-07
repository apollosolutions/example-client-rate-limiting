import express from "express";

const app = express();

const processRouterRequestStage = async (payload) => {
  // This is the object sent by the Router that you can act upon to update headers, context, auth claims, etc
  // If you update the "control" property from "Continue" to something like { "break": 400 }, it will terminate the request and return the specified HTTP error
  // See: https://www.apollographql.com/docs/router/customizations/coprocessor/
  console.log(payload);

  return payload;
};

app.post("/", express.json(), async (req, res) => {
  const payload = req.body;

  let response = payload;
  switch (payload.stage) {
    case "RouterRequest":
      response = await processRouterRequestStage(payload);
      break;
  }

  res.send(response);
});

app.listen(3007, () => {
  console.log("🚀 Server running at http://localhost:3007");
});
