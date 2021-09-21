import * as crypto from "m3o/crypto";

// Get news related to a currency
async function GetCryptocurrencyNews() {
  let cryptoService = new crypto.CryptoService(process.env.MICRO_API_TOKEN);
  let rsp = await cryptoService.news({
    symbol: "BTCUSD",
  });
  console.log(rsp);
}

await GetCryptocurrencyNews();
