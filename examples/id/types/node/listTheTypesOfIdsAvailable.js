import * as id from "m3o/id";

// List the types of IDs available. No query params needed.
async function ListTheTypesOfIdsAvailable() {
  let idService = new id.IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.types({});
  console.log(rsp);
}

await ListTheTypesOfIdsAvailable();
