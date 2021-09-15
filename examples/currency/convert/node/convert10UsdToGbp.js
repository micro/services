import * as currency from "@m3o/services/currency";

// Convert returns the currency conversion rate between two pairs e.g USD/GBP
async function Convert10usdToGbp() {
  let currencyService = new currency.CurrencyService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await currencyService.convert({
    amount: 10,
    from: "USD",
    to: "GBP",
  });
  console.log(rsp);
}

await Convert10usdToGbp();
