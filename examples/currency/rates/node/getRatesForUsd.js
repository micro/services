import * as currency from "m3o/currency";

// Rates returns the currency rates for a given code e.g USD
async function GetRatesForUsd() {
  let currencyService = new currency.CurrencyService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await currencyService.rates({
    code: "USD",
  });
  console.log(rsp);
}

await GetRatesForUsd();
