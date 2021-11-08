const { HelloworldService } = require("m3o/helloworld");

// Stream returns a stream of "Hello $name" responses
async function streamsResponsesFromTheServerUsingWebsockets() {
  let helloworldService = new HelloworldService(process.env.MICRO_API_TOKEN);
  let rsp = await helloworldService.stream({
    messages: 10,
    name: "John",
  });
  console.log(rsp);
}

streamsResponsesFromTheServerUsingWebsockets();
