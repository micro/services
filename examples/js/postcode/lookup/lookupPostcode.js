const { PostcodeService } = require("m3o/postcode");

// Lookup a postcode to retrieve the related region, county, etc
async function lookupPostcode() {
  let postcodeService = new PostcodeService(process.env.MICRO_API_TOKEN);
  let rsp = await postcodeService.lookup({
    postcode: "SW1A 2AA",
  });
  console.log(rsp);
}

lookupPostcode();
