const { sendSlackMessage } = require("./alert.js");

exports.App = class App {
  sendAlert(records) {
    records.forEach((record) => {
      let payload = record.value.payload;
      sendSlackMessage(process.env.SLACK_URL, payload);
    });

    return records;
  }

  async run(turbine) {
    let source = await turbine.resources("pg");

    let records = await source.records("customer_order");

    let data = await turbine.process(records, this.sendAlert, {
      SLACK_URL: process.env.SLACK_URL,
    });

    let destination = await turbine.resources("snowflake");

    await destination.write(data, "customer_order");
  }
};
