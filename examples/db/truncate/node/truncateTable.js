import * as db from "@m3o/services/db";

// Truncate the records in a table
async function TruncateTable() {
  let dbService = new db.DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.truncate({
    table: "users",
  });
  console.log(rsp);
}

await TruncateTable();
