import logging

from turbine.runtime import RecordList, Runtime

from enrich import enrich_user_email


def enrich_data(records: RecordList) -> RecordList:
    for record in records:
        try:
            logging.info(f"Got email: {record.value['payload']['email']}")

            enrichment = enrich_user_email(record.value["payload"]["email"])

            if enrichment:
                record.value["payload"]["full_name"] = enrichment.full_name
                record.value["payload"]["company"] = enrichment.company
                record.value["payload"]["location"] = enrichment.location
                record.value["payload"]["role"] = enrichment.role
                record.value["payload"]["seniority"] = enrichment.seniority
        except (KeyError, TypeError) as e:
            print(f"Error enriching data: {e}")
            raise

    return records


class App:
    @staticmethod
    async def run(turbine: Runtime):

        logging.basicConfig(level=logging.INFO)

        try:
            # Get remote resource
            resource = await turbine.resources("source_db")

            # Read from remote resource
            records = await resource.records("user_activity")

            # Register clearbit API Secret
            turbine.register_secrets("CLEARBIT_API_KEY")

            # Deploy function with source as input
            enriched = await turbine.process(records, enrich_data)

            # S3 Destination to warehouse our records
            destination = await turbine.resources("webhook-desk")

            # Write results out
            await destination.write(enriched, "user_activity_enriched")

        except ChildProcessError as cpe:
            print(cpe)
        except Exception as e:
            print(e)
