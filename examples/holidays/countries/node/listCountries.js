const { HolidaysService } = require("m3o/holidays");

// Get the list of countries that are supported by this API
async function listCountries() {
  let holidaysService = new HolidaysService(process.env.MICRO_API_TOKEN);
  let rsp = await holidaysService.countries({});
  console.log(rsp);
}

listCountries();
