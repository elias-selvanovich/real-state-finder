#!/bin/bash
if [ -z $1 ]; then
  echo "Please provide the client code"
  exit
fi

if [ -z $CLIENT_SECRET ]; then
  echo "Client secret can't be empty"
  exit
fi

CLIENT_CODE=$1

ACCESS_TOKEN=$(curl -s -X POST \
-H 'accept: application/json' \
-H 'content-type: application/x-www-form-urlencoded' \
'https://api.mercadolibre.com/oauth/token' \
-d 'grant_type=authorization_code' \
-d 'client_id='$APP_ID \
-d 'client_secret='$CLIENT_SECRET \
-d 'code='$CLIENT_CODE \
-d 'redirect_uri=https://localhost:8080/authenticate_callback' | jq -r '.access_token')

echo "Execute the following command"
echo "export ACCESS_TOKEN=$ACCESS_TOKEN"