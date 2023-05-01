#!/bin/bash

curl -X GET -H 'Authorization: Bearer $ACCESS_TOKEN' https://api.mercadolibre.com/classified_locations/countries/$COUNTRY_ID

printf "\n\n Search the state you want and get the corresponding ID. Export it to an env variable called STATE_ID"