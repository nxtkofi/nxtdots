#!/bin/bash
curl -s 'https://api.sunrisesunset.io/json?lat=52.2297&lng=21.0122&time_format=24' | jq -r '.results.dusk' | sed 's/:[0-9][0-9]$//' | tr -d '\n'
# TODO: Update coordinates for your location
# Get your coordinates from: https://www.latlong.net/
