const { IpService } = require("m3o/ip");

// Lookup the geolocation information for an IP address
async function lookupIpInfo() {
  let ipService = new IpService(process.env.MICRO_API_TOKEN);
  let rsp = await ipService.lookup({
    ip: "93.148.214.31",
  });
  console.log(rsp);
}

lookupIpInfo();
