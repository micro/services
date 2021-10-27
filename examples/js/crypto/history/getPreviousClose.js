const { CryptoService } = require("m3o/crypto");

// Returns the history for the previous close
async function getPreviousClose() {
  let cryptoService = new CryptoService(process.env.MICRO_API_TOKEN);
  let rsp = await cryptoService.history({
    symbol: "BTCUSD",
  });
  console.log(rsp);
}

getPreviousClose();
