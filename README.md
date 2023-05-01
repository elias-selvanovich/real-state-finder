## REAL STATE FINDER
### HOW TO AUTHENTICATE IN MERCADOLIBRE API

1. Login into https://developers.mercadolibre.com.ar
2. Get the APP_ID and CLIENT_SECRET values from your test application
3. Export these values to be accesible from the example scripts
```BASH
    export APP_ID=<YOUR_APP_ID>
    export CLIENT_SECRET=<YOUR_CLIENT_SECRET>
```
4. Execute the example `examples/authenticate.sh`. This will print the proper URL. Paste it in a browser and it'll ask for you to login to mercado libre. After login is succesfull it'll redirect you to a URL like `https://localhost:8080/?code=TG-...`. Get the `code` in the UrlParam and export it
```BASH
    export CLIENT_CODE=<YOUR_CLIENT_CODE>
``` 
5. You'll get a response like this one
```JSON
{"access_token":"APP_USR<REDACTED>","token_type":"Bearer","expires_in":21600,"scope":"offline_access read","user_id":45305128,"refresh_token":"TG-<REDACTED>"}
```
get the `access_token` field and export it
```BASH
    export ACCESS_TOKEN=<YOUR_ACCESS_TOKEN>
```

Now you can start running the examples in sequence.

### Running the examples
#### Getting the Country lists
**examples/get-countries.sh**

once you execute this script it'll retrieve the list of countries available. Search the one you need and export the COUNTRY_ID env variable with it

```SH
~ ./get-countries.sh
[{"id":"AR","name":"Argentina","locale":"es_AR","currency_id":"ARS"},...,{"id":"COL","name":"newCOL","locale":"es_COL","currency_id":"COLS"}]

Locate your desired country and export the env variable COUNTRY_ID with the desired ID

~ export COUNTRY_ID=AR
```

#### Getting Country details
**examples/get-country-details.sh**

This script retrieves the list of states from a country. Search the state you want and export the STATE_ID env variable with it's ID

```SH
~ ./get-country-details.sh

{"id":"AR","name":"Argentina","locale":"es_AR","currency_id":"ARS","decimal_separator":",","thousands_separator":".","time_zone":"GMT-03:00","geo_information":{"location":{"latitude":-38.416096,"longitude":-63.616673}},"states":[{"id":"TUxBUEJSQWwyMzA1","name":"Brasil"},{"id":"TUxBUENPU2ExMmFkMw","name":"Bs.As. Costa Atlántica"},...,{"id":"TUxBUFVTQWl1cXdlMg","name":"USA"}]}

 Search the state you want and get the corresponding ID. Export it to an env variable called STATE_ID
 export STATE_ID=TUxBUENBUGw3M2E1
```

#### Getting State details
**examples/get-state-details.sh**


This script retrieves the details from the state, in particular the list of cities in that state and the ID for the citie we'll retrieve details in the next step. Export from this script an env variable called CITY_ID

```SH
~ ./get-state-details.sh
{"id":"TUxBUENBUGw3M2E1","name":"Capital Federal","country":{"id":"AR","name":"Argentina"},"geo_information":{"location":{"latitude":-34.6143048,"longitude":-58.4401655}},"time_zone":"GMT-03:00","cities":[{"id":"TUxBQ0NBUGZlZG1sYQ","name":"Capital Federal"}]}

Export an env variable called CITY_ID with the desired city id

~ export CITY_ID=TUxBQ0NBUGZlZG1sYQ
```

