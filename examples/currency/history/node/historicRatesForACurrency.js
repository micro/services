const { CurrencyService } = require("m3o/currency");

// Returns the historic rates for a currency on a given date
async function historicRatesForAcurrency() {
  let currencyService = new CurrencyService(process.env.MICRO_API_TOKEN);
  let rsp = await currencyService.history({
    code: "USD",
    date: "2021-05-30",
  });
  console.log(rsp);
}

historicRatesForAcurrency();
