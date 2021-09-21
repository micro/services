import * as file from "m3o/file";

// Delete a file by project name/path
async function DeleteFile() {
  let fileService = new file.FileService(process.env.MICRO_API_TOKEN);
  let rsp = await fileService.delete({
    path: "/document/text-files/file.txt",
    project: "examples",
  });
  console.log(rsp);
}

await DeleteFile();
