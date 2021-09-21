import * as holidays from "m3o/holidays";

//
async function ListCountries() {
  let holidaysService = new holidays.HolidaysService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await holidaysService.countries({});
  console.log(rsp);
}

await ListCountries();
