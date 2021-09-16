import * as address from "@m3o/services/address";

// Lookup a list of UK addresses by postcode
async function LookupPostcode() {
  let addressService = new address.AddressService(process.env.MICRO_API_TOKEN);
  let rsp = await addressService.lookupPostcode({
    postcode: "SW1A 2AA",
  });
  console.log(rsp);
}

await LookupPostcode();
