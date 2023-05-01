#!/bin/bash

curl -X GET -H 'Authorization: Bearer $ACCESS_TOKEN' https://api.mercadolibre.com/classified_locations/states/$STATE_ID

printf "\n\n%s" "Export an env variable called CITY_ID with the desired city id"