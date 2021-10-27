const { IdService } = require("m3o/id");

// Generate a unique ID. Defaults to uuid.
async function generateAuniqueId() {
  let idService = new IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.generate({
    type: "uuid",
  });
  console.log(rsp);
}

generateAuniqueId();
