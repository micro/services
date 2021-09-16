import * as db from "@m3o/services/db";

// Read data from a table. Lookup can be by ID or via querying any field in the record.
async function ReadRecords() {
  let dbService = new db.DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.read({
    query: "age == 43",
    table: "users",
  });
  console.log(rsp);
}

await ReadRecords();
