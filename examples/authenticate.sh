#!/bin/bash

if [ -z $APP_ID ]; then
  echo "APP_ID is empty, please set it before running this command"
  exit
fi

AUTHENTICATE_URL="https://auth.mercadolibre.com.ar/authorization?response_type=code&client_id=$APP_ID&redirect_uri=https://localhost:8080/authenticate_callback"

echo $AUTHENTICATE_URL

open -a Firefox $AUTHENTICATE_URL

echo "Retrieve the code from the URL in firefox and use it to pass it as a parameter to the auth-token script"