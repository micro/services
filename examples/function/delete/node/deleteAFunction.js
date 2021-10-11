import * as fx from "m3o/function";

// Delete a function by name
async function DeleteAfunction() {
  let functionService = new fx.FunctionService(process.env.MICRO_API_TOKEN);
  let rsp = await functionService.delete({
    name: "my-first-func",
    project: "tests",
  });
  console.log(rsp);
}

await DeleteAfunction();
