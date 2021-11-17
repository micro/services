const { CurrencyService } = require("micro-js-client/currency");

// Rates returns the currency rates for a given code e.g USD
async function getRatesForUsd() {
  let currencyService = new CurrencyService(process.env.MICRO_API_TOKEN);
  let rsp = await currencyService.rates({
    code: "USD",
  });
  console.log(rsp);
}

getRatesForUsd();
