const { AddressService } = require("m3o/address");

// Lookup a list of UK addresses by postcode
async function lookupPostcode() {
  let addressService = new AddressService(process.env.MICRO_API_TOKEN);
  let rsp = await addressService.lookupPostcode({
    postcode: "SW1A 2AA",
  });
  console.log(rsp);
}

lookupPostcode();
