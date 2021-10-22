const { CryptoService } = require("m3o/crypto");

// Get the last quote for a given crypto ticker
async function getAcryptocurrencyQuote() {
  let cryptoService = new CryptoService(process.env.MICRO_API_TOKEN);
  let rsp = await cryptoService.quote({
    symbol: "BTCUSD",
  });
  console.log(rsp);
}

getAcryptocurrencyQuote();
