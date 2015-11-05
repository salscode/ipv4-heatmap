$(document).ready(function() {
	var map = L.map('map').setView([35.889599, -78.924400], 5);
	L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
		attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery ? <a href="http://mapbox.com">Mapbox</a>',
		minZoom: 2,
		id: 'skytoastersal.o27gaj6k',
		accessToken: 'pk.eyJ1Ijoic2t5dG9hc3RlcnNhbCIsImEiOiJhYWJjNDNhNzc0MjJkZmY4ZjU0NzliOTU4Y2IwZWY5OSJ9.YZfsXmBNb_C1OAQ7JPOFRg'
	}).addTo(map);
	
	var heat;
	jQuery.getScript("http://heatmap.salscode.com/locations", function( data, textStatus, jqxhr ) {
		heat = L.heatLayer(locations, { radius: 40, blur: 25, gradient: {0.1: 'blue', 0.2: 'lime', 0.5: 'red'} }).addTo(map);
	});
	
	function onMapChange(e) {
		var bounds = map.getBounds().pad(1.01);
		jQuery.getScript("http://heatmap.salscode.com/locations/" + bounds.getNorth() + "/" + bounds.getEast() + "/" + bounds.getSouth() + "/" + bounds.getWest());
		heat.setLatLngs(locations);
		heat.redraw();
	}
	
	map.on('load', onMapChange);
	map.on('moveend', onMapChange);
});