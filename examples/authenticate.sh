#!/bin/bash

if [ -z $APP_ID ]; then
  echo "APP_ID is empty, please set it before running this command"
  exit
fi

echo "https://auth.mercadolibre.com.ar/authorization?response_type=code&client_id=$APP_ID&redirect_uri=https://localhost:8080/"