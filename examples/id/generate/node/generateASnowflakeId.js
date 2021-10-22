const { IdService } = require("m3o/id");

// Generate a unique ID. Defaults to uuid.
async function generateAsnowflakeId() {
  let idService = new IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.generate({
    type: "snowflake",
  });
  console.log(rsp);
}

generateAsnowflakeId();
