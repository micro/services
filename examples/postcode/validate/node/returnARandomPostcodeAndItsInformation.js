import * as postcode from "m3o/postcode";

// Validate a postcode.
async function ReturnArandomPostcodeAndItsInformation() {
  let postcodeService = new postcode.PostcodeService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await postcodeService.validate({
    postcode: "SW1A 2AA",
  });
  console.log(rsp);
}

await ReturnArandomPostcodeAndItsInformation();
