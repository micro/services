import * as evchargers from "m3o/evchargers";

// Retrieve reference data as used by this API
async function GetReferenceData() {
  let evchargersService = new evchargers.EvchargersService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await evchargersService.referenceData({});
  console.log(rsp);
}

await GetReferenceData();
