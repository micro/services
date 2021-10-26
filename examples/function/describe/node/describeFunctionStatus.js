const { FunctionService } = require("m3o/function");

// Get the info for a deployed function
async function describeFunctionStatus() {
  let functionService = new FunctionService(process.env.MICRO_API_TOKEN);
  let rsp = await functionService.describe({
    name: "my-first-func",
    project: "tests",
  });
  console.log(rsp);
}

describeFunctionStatus();
