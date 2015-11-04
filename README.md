# ipv4-heatmap

## Database
The heatmap data is stored in a single MySQL table. Only the latitude, longitude and IP count are stored for effeciency. The coordinates are stored to 4 decimals of precision which is more than enough for this use.

## Data
Originally used data from [GeoLite Data](http://dev.maxmind.com/geoip/legacy/geolite/) by MaxMind, imported via import.go. Now uses GeoLite 2 data from [GeoLite2 Data](http://dev.maxmind.com/geoip/geoip2/geolite2/), also from MaxMind.

## Importing the Data
After the database is configured, `run main.go import` to import the GeoLite data or `run main.go import2` to import the GeoLite2 data.

## Main Server File
The main.go file powers the endpoint and also serves the main map map along with the static assets needed for the map.

## API
There are 2 APIs:
* /locations - Used to fetch a list of all location/IP data.
* /locations/{latitude1}/{longitude1}/{latitude2}/{longitude2} - Used to fetch location/IP data within a Lat/Lon bounding box.
