#!/bin/bash
set -x

curl -X POST \
-H 'accept: application/json' \
-H 'content-type: application/x-www-form-urlencoded' \
'https://api.mercadolibre.com/oauth/token' \
-d 'grant_type=authorization_code' \
-d 'client_id='$APP_ID \
-d 'client_secret='$CLIENT_SECRET \
-d 'code='$CLIENT_CODE \
-d 'redirect_uri=https://localhost:8080/'