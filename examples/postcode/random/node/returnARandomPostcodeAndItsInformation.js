import * as postcode from "@m3o/services/postcode";

// Return a random postcode and its related info
async function ReturnArandomPostcodeAndItsInformation() {
  let postcodeService = new postcode.PostcodeService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await postcodeService.random({});
  console.log(rsp);
}

await ReturnArandomPostcodeAndItsInformation();
