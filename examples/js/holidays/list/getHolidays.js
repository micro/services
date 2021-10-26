const { HolidaysService } = require("m3o/holidays");

// List the holiday dates for a given country and year
async function getHolidays() {
  let holidaysService = new HolidaysService(process.env.MICRO_API_TOKEN);
  let rsp = await holidaysService.list({
    country_code: "GB",
    year: 2022,
  });
  console.log(rsp);
}

getHolidays();
