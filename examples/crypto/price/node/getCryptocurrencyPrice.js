import * as crypto from "@m3o/services/crypto";

// Get the last price for a given crypto ticker
async function GetCryptocurrencyPrice() {
  let cryptoService = new crypto.CryptoService(process.env.MICRO_API_TOKEN);
  let rsp = await cryptoService.price({
    symbol: "BTCUSD",
  });
  console.log(rsp);
}

await GetCryptocurrencyPrice();
