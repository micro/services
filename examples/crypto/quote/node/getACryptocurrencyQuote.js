import * as crypto from "m3o/crypto";

// Get the last quote for a given crypto ticker
async function GetAcryptocurrencyQuote() {
  let cryptoService = new crypto.CryptoService(process.env.MICRO_API_TOKEN);
  let rsp = await cryptoService.quote({
    symbol: "BTCUSD",
  });
  console.log(rsp);
}

await GetAcryptocurrencyQuote();
