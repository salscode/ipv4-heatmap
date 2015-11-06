# ipv4-heatmap

## Setup
Create a VM with CentOS 7 (other OSes work too). Install Maria DB using the instructions found [here](https://downloads.mariadb.org/mariadb/repositories/). Then install Go using the instructions found [here](https://golang.org/doc/install). Use the files in /etc/sysconfig and /etc/systemd to impelement main.go as the system service heatmap.

## Database
The heatmap data is stored in a single MySQL table. Only the latitude, longitude and IP count are stored for effeciency. The coordinates are stored to 4 decimals of precision which is more than enough for this use.

## Data
Originally used data from [GeoLite Data](http://dev.maxmind.com/geoip/legacy/geolite/) by MaxMind, imported via import.go. Now uses GeoLite 2 data from [GeoLite2 Data](http://dev.maxmind.com/geoip/geoip2/geolite2/), also from MaxMind.

## Importing the Data
After the database is configured, `run main.go import` to import the GeoLite data or `run main.go import2` to import the GeoLite2 data.

## Endpoint Main.go
The main.go file powers the endpoint and also delegates data importing.

## Endpoint API
The APIs return the data encoded in a protocol buffer object for the caller to decode and use.
There are 2 APIs:
* /locations - Used to fetch a list of all location/IP data.
* /locations/{latitude1}/{longitude1}/{latitude2}/{longitude2} - Used to fetch location/IP data within a Lat/Lon bounding box.

## Frontend Main.go
The main.go file in the frontend fetches the data and decodes it. It also serves the static assets needed for the map.

## Testing
Currently, you can run a basic endpoint test by visiting [test page](http://heatmap.salscode.com/test). The page triggers two endpoint calls and does a basic check on the results.

## TODO/Issues
* The heat intensities are not bright enough when zoomed in at the city level.
* When zooming in from world view to state-level the heatmap doesn't appear to update correctly until you zoom back out one step.
* Some mid-level zoom layers could be improved in terms of heatmap display.
* Create a way to cancel previous pending requests when a new request is issued. Helps when the user changes the display multiple times in a row.
* Expand testing and use the go language testing framework to perform more tests.
