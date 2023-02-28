#!/usr/bin/python
import os

from deltalake import DeltaTable

STORAGE_OPTIONS = {
    "AWS_ACCESS_KEY_ID": os.getenv("AWS_ACCESS_KEY_ID"),
    "AWS_SECRET_ACCESS_KEY": os.getenv("AWS_SECRET_ACCESS_KEY"),
    "AWS_REGION": os.getenv("AWS_REGION"),
    "AWS_S3_ALLOW_UNSAFE_RENAME": "true",
}

S3_URI = os.getenv("AWS_URI")

dt = DeltaTable(table_uri=S3_URI, storage_options=STORAGE_OPTIONS)

print("*** DELTA TABLE SCHEMA ***")
print(dt.schema().to_pyarrow())

print("\n\n*** TABLE ***")
print(dt.to_pandas())
