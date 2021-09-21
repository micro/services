import * as id from "m3o/id";

// Generate a unique ID. Defaults to uuid.
async function GenerateAsnowflakeId() {
  let idService = new id.IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.generate({
    type: "snowflake",
  });
  console.log(rsp);
}

await GenerateAsnowflakeId();
