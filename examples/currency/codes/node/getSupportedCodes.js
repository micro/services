import * as currency from "m3o/currency";

// Codes returns the supported currency codes for the API
async function GetSupportedCodes() {
  let currencyService = new currency.CurrencyService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await currencyService.codes({});
  console.log(rsp);
}

await GetSupportedCodes();
