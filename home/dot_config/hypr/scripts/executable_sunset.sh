#!/bin/bash
curl -s 'https://api.sunrisesunset.io/json?lat=50.860130&lng=17.470869&time_format=24' | jq -r '.results.sunset' | sed 's/:[0-9][0-9]$//' | tr -d '\n'
# TODO: Update coordinates for your location
# Get your coordinates from: https://www.latlong.net/
