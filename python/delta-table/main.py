import os
import logging
import sys

import sentry_sdk
from turbine.runtime import RecordList, Runtime

import delta  # Our Delta Lake writing code
import enrich  # Our Data Enrichment code

logging.basicConfig(level=logging.INFO)

sentry_sdk.init(
    dsn=os.getenv("SENTRY_DSN"),
    traces_sample_rate=0.5,
)


def format_and_enrich(records: RecordList) -> RecordList:
    data = {}

    for record in records:

        payload = record.value["payload"]

        geolocation = enrich.get_geo_location_from_postcode(payload["postcode"])
        payload["latitude"] = geolocation.latitude
        payload["longitude"] = geolocation.longitude

        for key, val in payload.items():
            if key in data:
                data[key].append(val)
            else:
                data.update({key: [val]})

    delta.write_records(data=data)

    return records


class App:
    @staticmethod
    async def run(turbine: Runtime):
        try:

            # secrets for writing to a S3 deltatable
            turbine.register_secrets("AWS_ACCESS_KEY_ID")
            turbine.register_secrets("AWS_SECRET_ACCESS_KEY")
            turbine.register_secrets("AWS_REGION")
            turbine.register_secrets("AWS_URI")

            # DSN for capturing exceptions and sending them to Sentry
            turbine.register_secrets("SENTRY_DSN")

            # API key for using Google Address validation API to enhance a record
            # with additional data
            turbine.register_secrets("GOOGLE_API_KEY")

            source = await turbine.resources("postgres_source")
            raw = await source.records("employees")

            processed = await turbine.process(raw, format_and_enrich)
            silver = await turbine.resources("url")

            await silver.write(
                processed,
                "employees_enriched",
                {"table.name.format": "employees_enriched"},
            )

        except Exception as e:
            print(e, file=sys.stderr)
