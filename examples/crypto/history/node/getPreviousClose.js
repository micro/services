import * as crypto from "m3o/crypto";

// Returns the history for the previous close
async function GetPreviousClose() {
  let cryptoService = new crypto.CryptoService(process.env.MICRO_API_TOKEN);
  let rsp = await cryptoService.history({
    symbol: "BTCUSD",
  });
  console.log(rsp);
}

await GetPreviousClose();
