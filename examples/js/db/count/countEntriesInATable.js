const { DbService } = require("m3o/db");

// Count records in a table
async function countEntriesInAtable() {
  let dbService = new DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.count({
    table: "users",
  });
  console.log(rsp);
}

countEntriesInAtable();
