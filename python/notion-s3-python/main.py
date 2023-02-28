import hashlib
import logging
import sys

from turbine.runtime import RecordList
from turbine.runtime import Runtime

logging.basicConfig(level=logging.INFO)

def character_count(records: RecordList) -> RecordList:
    for record in records:
        try:
            payload = record.value["payload"]["after"]
            metadata = payload["metadata"]

            plaintext_len = len(payload.get('plaintext'))

            # Add the character count in the metadata JSON
            metadata.update({"characterCount": plaintext_len})
        except Exception as e:
            print("Error occurred while parsing records: " + str(e))
    return records

class App:
    @staticmethod
    async def run(turbine: Runtime):
        try:
             # Source Notion Resource
            source = await turbine.resources("company-notion")

            # Currently, we ignore the collection name passed in 
            # for Notion resources. Your turbine app will pull all 
            # pages and databases the Notion integration has access to. 
            records = await source.records("")

            # Our function which performs character counts
            counted_records = await turbine.process(records, character_count)

            # S3 destination resource
            destination = await turbine.resources("s3")
            # Writes to specific folder in S3 bucket
            await destination.write(counted_records, "notion-pages-with-character-count", {})
        except Exception as e:
            print(e, file=sys.stderr)
