import * as db from "m3o/db";

// Count records in a table
async function CountEntriesInAtable() {
  let dbService = new db.DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.count({
    table: "users",
  });
  console.log(rsp);
}

await CountEntriesInAtable();
