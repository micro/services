import * as currency from "m3o/currency";

// Convert returns the currency conversion rate between two pairs e.g USD/GBP
async function ConvertUsdToGbp() {
  let currencyService = new currency.CurrencyService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await currencyService.convert({
    from: "USD",
    to: "GBP",
  });
  console.log(rsp);
}

await ConvertUsdToGbp();
