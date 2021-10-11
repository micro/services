import * as fx from "m3o/function";

// Call a function
async function CallAfunction() {
  let functionService = new fx.FunctionService(process.env.MICRO_API_TOKEN);
  let rsp = await functionService.call({
    name: "my-first-func",
    request: {},
  });
  console.log(rsp);
}

await CallAfunction();
