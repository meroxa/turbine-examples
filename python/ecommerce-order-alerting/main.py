import logging
import sys

import hashlib

from datetime import date

from turbine.runtime import RecordList, Runtime

from alert import send_slack_alert
from constants import WEBHOOK_URL

logging.basicConfig(level=logging.INFO)

"""
Standarized function signature that takes a list of Records 
(database rows) and returns a list of Records with desired
mutations
"""


def send_alert(records: RecordList) -> RecordList:
    for record in records:
        try:
            payload = record.value["payload"]

            # Hash the email
            payload["customer_email"] = hashlib.sha256(
                payload["customer_email"].encode("utf-8")
            ).hexdigest()

            send_slack_alert(webhook_url=WEBHOOK_URL, payload=payload)

            record.value["payload"] = payload

        except Exception as e:
            logging.warning(f"Error occurred while parsing records: {e}")

    return records


class App:
    @staticmethod
    async def run(turbine: Runtime):

        try:
            # Postgres database being monitored for changes
            # by the Meroxa Platform
            source = await turbine.resources("source_db")

            # Tell Turbine and Meroxa that this is the collection
            # that should be watched
            records = await source.records("customer_order")

            # Apply my function to the rows/records from the
            # specified collection and database
            data = await turbine.process(records, send_alert)

            # S3 Destination to warehouse our records
            destination = await turbine.resources("warehouse")

            # Use Turbine to write records to a collection called
            # `customer_orders $DATE` using Meroxa
            await destination.write(data, f"customer_orders {str(date.today())}")
        except Exception as e:
            print(e, file=sys.stderr)
