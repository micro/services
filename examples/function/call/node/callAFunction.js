const { FunctionService } = require("micro-js-client/function");

// Call a function by name
async function callAfunction() {
  let functionService = new FunctionService(process.env.MICRO_API_TOKEN);
  let rsp = await functionService.call({
    name: "my-first-func",
    request: {},
  });
  console.log(rsp);
}

callAfunction();
