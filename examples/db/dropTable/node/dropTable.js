const { DbService } = require("micro-js-client/db");

// Drop a table in the DB
async function dropTable() {
  let dbService = new DbService(process.env.MICRO_API_TOKEN);
  let rsp = await dbService.dropTable({
    table: "users",
  });
  console.log(rsp);
}

dropTable();
