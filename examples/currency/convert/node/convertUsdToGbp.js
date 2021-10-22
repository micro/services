const { CurrencyService } = require("m3o/currency");

// Convert returns the currency conversion rate between two pairs e.g USD/GBP
async function convertUsdToGbp() {
  let currencyService = new CurrencyService(process.env.MICRO_API_TOKEN);
  let rsp = await currencyService.convert({
    from: "USD",
    to: "GBP",
  });
  console.log(rsp);
}

convertUsdToGbp();
