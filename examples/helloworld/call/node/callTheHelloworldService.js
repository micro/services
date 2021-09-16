import * as helloworld from "m3o/helloworld";

// Call returns a personalised "Hello $name" response
async function CallTheHelloworldService() {
  let helloworldService = new helloworld.HelloworldService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await helloworldService.call({
    name: "John",
  });
  console.log(rsp);
}

await CallTheHelloworldService();
