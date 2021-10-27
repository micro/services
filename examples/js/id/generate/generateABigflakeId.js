const { IdService } = require("m3o/id");

// Generate a unique ID. Defaults to uuid.
async function generateAbigflakeId() {
  let idService = new IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.generate({
    type: "bigflake",
  });
  console.log(rsp);
}

generateAbigflakeId();
