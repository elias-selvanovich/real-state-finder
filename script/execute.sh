#! /bin/bash

AUTHENTICATE_SCRIPT="script/1_authenticate.sh"
AUTH_TOKEN_SCRIPT="script/2_auth-token.sh"

CURRENT_DIR=$(pwd)

echo $CURRENT_DIR

# validate that we have everything we need

if [ -z $APP_ID ]; then
  echo "APP_ID is empty, please set it before running this command"
  exit
fi

if [ -z $CLIENT_SECRET ]; then
  echo "Client secret can't be empty"
  exit
fi

# if ACCESS_TOKEN is not set then we need to retrieve it

if [ -z $ACCESS_TOKEN ]; then

    source "$CURRENT_DIR/$AUTHENTICATE_SCRIPT"

    read -p "Please paste the code from the response:" CLIENT_CODE

    export ACCESS_TOKEN=$(source "$CURRENT_DIR/$AUTH_TOKEN_SCRIPT" $CLIENT_CODE)
fi

go run ./cmd/real-state-finder/main.go -access-token=$ACCESS_TOKEN -search-result-limit=50 -command=search -max-offset=20

go run ./cmd/real-state-finder/main.go -access-token=$ACCESS_TOKEN -search-result-limit=50 -command=generate-html

git status
git add index.html
git commit -m "Deploy index.html"
git push origin main




