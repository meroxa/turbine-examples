import json
from typing import Any

from slack_sdk.webhook import WebhookClient


def build_body(payload: dict[str, Any]) -> list:
    msg = f" :dollar: A new order has arrived :dollar: :\n ```{json.dumps(payload, indent=4)}```"
    return [
        {
            "type": "section",
            "text": {"type": "mrkdwn", "text": msg},
        }
    ]


def send_slack_alert(webhook_url: str, payload: dict[str, Any]) -> None:
    webhook = WebhookClient(url=webhook_url)
    webhook.send(blocks=build_body(payload=payload))
