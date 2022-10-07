# Turbine

## Pre-requisites

### MongoDB Resource
```shell
meroxa resources create mdb --type mongodb -u "mongodb+srv://$MONGODB_USER:$MONGODB_PASSWORD@$MONGODB_HOST/myFirstDatabase?retryWrites=true&w=majority"
```
For this example I created a Shared MongoDB instance on [MongoDB Atlas](https://www.mongodb.com/atlas).

Example MongoDB Document:
```json
{
	"_id": {
		"$oid": "6318c7474d577b9917d7ff35"
	},
	"user_id": "100",
	"activity": "user logged in",
	"vip": true,
	"created_at": "1662568354",
	"updated_at": "1662568354"
}
```