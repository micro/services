import * as fx from "m3o/function";

//
async function ListFunctions() {
  let functionService = new fx.FunctionService(process.env.MICRO_API_TOKEN);
  let rsp = await functionService.list({});
  console.log(rsp);
}

await ListFunctions();
