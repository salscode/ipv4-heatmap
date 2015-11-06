$(document).ready(function() {
	var map = L.map('map').setView([0, 0], 2);
	L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
		attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery ? <a href="http://mapbox.com">Mapbox</a>',
		minZoom: 2,
		id: 'skytoastersal.o27gaj6k',
		accessToken: 'pk.eyJ1Ijoic2t5dG9hc3RlcnNhbCIsImEiOiJhYWJjNDNhNzc0MjJkZmY4ZjU0NzliOTU4Y2IwZWY5OSJ9.YZfsXmBNb_C1OAQ7JPOFRg'
	}).addTo(map);
	
	var heat = L.heatLayer([], heatOpts).addTo(map);
	var heatOpts = { minOpacity: 0.1, radius: 20, blur: 15, gradient: {0.1: 'blue', 0.2: 'lime', 0.6: 'red'} };
		
	function onMapChange(e) {
		var bounds = map.getBounds().pad(1.01);
		jQuery.getScript("http://heatmap.salscode.com/locations/" + bounds.getNorth() + "/" + bounds.getEast() + "/" + bounds.getSouth() + "/" + bounds.getWest(), function()
		{
			heat.setLatLngs(locations);
			heat.redraw();
		});
	}
	
	function onZoomChange(e) {
		if (map.getZoom() < 4) {
			heatOpts.radius = 20;
			heatOpts.blur = 15;
		} else if (map.getZoom() < 9) {
			heatOpts.radius = 25;
			heatOpts.blur = 20;
		} else {
			heatOpts.radius = 30;
			heatOpts.blur = 25;
		}	
		
		heat.setOptions(heatOpts);
	}
	
	map.on('load', onMapChange);
	map.on('moveend', onMapChange);
	map.on('zoomend', onZoomChange);
});