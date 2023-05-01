#!/bin/bash

curl -X GET -H 'Authorization: Bearer $ACCESS_TOKEN' https://api.mercadolibre.com/classified_locations/countries

printf "\r\nLocate your desired country and export the env variable COUNTRY_ID with the desired ID"