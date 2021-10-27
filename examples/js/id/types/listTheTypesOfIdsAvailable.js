const { IdService } = require("m3o/id");

// List the types of IDs available. No query params needed.
async function listTheTypesOfIdsAvailable() {
  let idService = new IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.types({});
  console.log(rsp);
}

listTheTypesOfIdsAvailable();
