const { PostcodeService } = require("m3o/postcode");

// Return a random postcode and its related info
async function returnArandomPostcodeAndItsInformation() {
  let postcodeService = new PostcodeService(process.env.MICRO_API_TOKEN);
  let rsp = await postcodeService.random({});
  console.log(rsp);
}

returnArandomPostcodeAndItsInformation();
