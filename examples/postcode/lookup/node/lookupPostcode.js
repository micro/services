import * as postcode from "m3o/postcode";

// Lookup a postcode to retrieve the related region, county, etc
async function LookupPostcode() {
  let postcodeService = new postcode.PostcodeService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await postcodeService.lookup({
    postcode: "SW1A 2AA",
  });
  console.log(rsp);
}

await LookupPostcode();
