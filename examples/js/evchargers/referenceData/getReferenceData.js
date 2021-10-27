const { EvchargersService } = require("m3o/evchargers");

// Retrieve reference data as used by this API and in conjunction with the Search endpoint
async function getReferenceData() {
  let evchargersService = new EvchargersService(process.env.MICRO_API_TOKEN);
  let rsp = await evchargersService.referenceData({});
  console.log(rsp);
}

getReferenceData();
