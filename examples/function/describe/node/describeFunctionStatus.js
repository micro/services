import * as fx from "m3o/function";

//
async function DescribeFunctionStatus() {
  let functionService = new fx.FunctionService(process.env.MICRO_API_TOKEN);
  let rsp = await functionService.describe({
    name: "my-first-func",
    project: "tests",
  });
  console.log(rsp);
}

await DescribeFunctionStatus();
