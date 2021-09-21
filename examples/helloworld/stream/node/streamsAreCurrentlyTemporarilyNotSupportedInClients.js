import * as helloworld from "m3o/helloworld";

// Stream returns a stream of "Hello $name" responses
async function StreamsAreCurrentlyTemporarilyNotSupportedInClients() {
  let helloworldService = new helloworld.HelloworldService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await helloworldService.stream({
    name: "not supported",
  });
  console.log(rsp);
}

await StreamsAreCurrentlyTemporarilyNotSupportedInClients();
