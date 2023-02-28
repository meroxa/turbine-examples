
# Delta Table Example 

## Overview 

Data App Architecture: 
![Data App Architecture](/delta-table/img/arch.png)

## Requirments 
### Resources
- Postgres (or any data source) https://docs.meroxa.com/platform/resources/postgresql/setup
- S3 Resource(s) set up on the meroxa platform: https://docs.meroxa.com/platform/resources/amazon-s3
- Snowflake set up on the meroxa platform: https://docs.meroxa.com/platform/resources/snowflake


### Environement Variables
#### AWS Creds
| Env Var        | Description          
| ------------- |-------------|
| AWS_ACCESS_KEY_ID      | AWS Access Key for user accessing buckets   |
| AWS_SECRET_ACCESS_KEY     | AWS Secret for user accessing buckets      |  
| AWS_REGION | Region the bucket was created in     |  
| AWS_URI | The actual URI of the bucket (e.g.: ```S3://bucket-name/key-name``` )

#### Sentry (for logging )
| Env Var        | Description          
| ------------- |-------------|
| SENTRY_DSN      | Sentry Data Source Name (DSN) to upload logs and errors   |


#### Google API Creds (for data enrichement )
| Env Var        | Description          
| ------------- |-------------|
| GOOGLE_API_KEY      | API Key to access Google Location API   |

## Library
Uses [delta-rs](https://github.com/delta-io/delta-rs) which provides an interface with various delta lake providers. 


## Setup 
Included in this example is [a data generation script](/delta-table/scripts/data_generation.py) that uses [Faker](https://github.com/joke2k/faker) to populate some example data. This script requires the postgres connection URL to connect. 

## Deployment
With resources set up on the meroxa platform and in this application's directory:

``` bash
$ meroxa apps deploy 
```