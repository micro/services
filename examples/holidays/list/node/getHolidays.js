import * as holidays from "m3o/holidays";

// List the holiday dates for a given country and year
async function GetHolidays() {
  let holidaysService = new holidays.HolidaysService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await holidaysService.list({
    country_code: "GB",
    year: 2022,
  });
  console.log(rsp);
}

await GetHolidays();
