const { QrService } = require("m3o/qr");

// Generate a QR code with a specific text and size
async function generateAqrCode() {
  let qrService = new QrService(process.env.MICRO_API_TOKEN);
  let rsp = await qrService.generate({
    size: 300,
    text: "https://m3o.com/qr",
  });
  console.log(rsp);
}

generateAqrCode();
