import * as holidays from "m3o/holidays";

// Get the list of countries that are supported by this API
async function ListCountries() {
  let holidaysService = new holidays.HolidaysService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await holidaysService.countries({});
  console.log(rsp);
}

await ListCountries();
