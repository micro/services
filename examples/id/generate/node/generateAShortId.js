const { IdService } = require("m3o/id");

// Generate a unique ID. Defaults to uuid.
async function generateAshortId() {
  let idService = new IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.generate({
    type: "shortid",
  });
  console.log(rsp);
}

generateAshortId();
