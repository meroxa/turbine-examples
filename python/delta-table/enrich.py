from __future__ import annotations

import os

import requests

VALIDATION_API_URL = "https://maps.googleapis.com/maps/api/geocode/json"
GOOGLE_KEY = os.getenv("GOOGLE_API_KEY")

NULL_ISLAND = 0.0


class GeoLocation:
    def __init__(self, lat, lon):
        self.latitude = lat
        self.longitude = lon


def get_geo_location_from_postcode(postcode: str) -> GeoLocation:

    params = {"address": postcode, "key": GOOGLE_KEY}
    response = requests.post(url=VALIDATION_API_URL, params=params)

    maybe_geocode = response.json().get("results")
    response_is_empty = len(maybe_geocode) == 0

    lat = (
        maybe_geocode[0].get("geometry").get("location").get("lat")
        if not response_is_empty
        else NULL_ISLAND
    )
    lon = (
        maybe_geocode[0].get("geometry").get("location").get("lng")
        if not response_is_empty
        else NULL_ISLAND
    )

    return GeoLocation(
        lat=lat,
        lon=lon,
    )
