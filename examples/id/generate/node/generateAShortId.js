import * as id from "m3o/id";

// Generate a unique ID. Defaults to uuid.
async function GenerateAshortId() {
  let idService = new id.IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.generate({
    type: "shortid",
  });
  console.log(rsp);
}

await GenerateAshortId();
