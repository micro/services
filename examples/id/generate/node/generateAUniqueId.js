import * as id from "@m3o/services/id";

// Generate a unique ID. Defaults to uuid.
async function GenerateAuniqueId() {
  let idService = new id.IdService(process.env.MICRO_API_TOKEN);
  let rsp = await idService.generate({
    type: "uuid",
  });
  console.log(rsp);
}

await GenerateAuniqueId();
