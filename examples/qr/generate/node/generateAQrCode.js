import * as qr from "m3o/qr";

//
async function GenerateAqrCode() {
  let qrService = new qr.QrService(process.env.MICRO_API_TOKEN);
  let rsp = await qrService.generate({
    size: 300,
    text: "https://m3o.com/qr",
  });
  console.log(rsp);
}

await GenerateAqrCode();
