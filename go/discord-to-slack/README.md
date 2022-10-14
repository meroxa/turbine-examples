# Discord -> Slack Turbine App

## Prerequisites

**Create Discord App and Bot**
This is beyond the scope of this README. You can find more details in the [Discord Developer Portal](https://discord.com/developers/docs/intro).

**Create HTTP Source Resource**
In order to integrate with the Discord API we will be using the HTTP Source Connector. This connector is currently only
available as part of a developer preview. You can request access via Meroxa Support.

```shell
meroxa resource create discord \
        --type url \
        --url https://discord.com/api/v10
```