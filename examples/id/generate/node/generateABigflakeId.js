import * as id from "@m3o/services/id";

// Generate a unique ID. Defaults to uuid.
async function GenerateAbigflakeId() {
  let idService = new id.IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.generate({
    type: "bigflake",
  });
  console.log(rsp);
}

await GenerateAbigflakeId();
