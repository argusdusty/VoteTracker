var info = L.control();

info.onAdd = function (map) {
	this._div = L.DomUtil.create('div', 'info');
	this.reset();
	return this._div;
};

info.update = function (props) {
	this._div.innerHTML = '<h4>' + props.name + '</h4>'
	if (props.total_votes == 0) {
		this._div.innerHTML += "No results yet"
	} else {
		for (var i in props.candidates) {
			if (props.total_votes != 0 && (props.candidates[i].votes) >= (summary.map_portion_thresh*props.total_votes)) {
				this._div.innerHTML += props.candidates[i].candidate + ': ' + props.candidates[i].votes + '(' + Math.round(10000*props.candidates[i].votes/props.total_votes)/100 + '%)</br>';
			}
		}
		this._div.innerHTML += 'Estimated percent complete: ' + Math.round(10000*props.portion_complete)/100 + '%</br>';
	}
};

info.reset = function () {
	this._div.innerHTML = '<h4>Area</h4>Hover over a region';
};

function getProps(fips) {
	if (!(fips in summary.regions)) {
		return {};
	}
	return summary.regions[fips];
}

var margin_scale = 1.0;

function loadColor(candidate) {
	if (candidate in colors) {
		return colors[candidate];
	}

	if ("" in colors) {
		return colors[""];
	}

	return "#CCCCCC";
}

function getColor(props) {
	if (!('exclude' in props) || props.exclude) {
		return chroma('#D4D4D4');
	}
	if ((!('candidates' in props)) || props.candidates == null || props.candidates.length <= 0) {
		return chroma('#EEEEEE');
	}

	if (props.candidates.length == 1) {
		return chroma(loadColor(props.candidates[0].candidate));
	}

	if (props.total_votes == 0 || props.candidates[0].votes == props.candidates[1].votes) {
		return chroma('#EEEEEE');
	}

	margin = (props.candidates[0].votes-props.candidates[1].votes)/props.total_votes;
	margin /= margin_scale;

	return chroma(loadColor(props.candidates[0].candidate)).alpha((margin*0.65)+0.35);
}

function style(feature) {
	var color = getColor(getProps(feature.properties.GEOID));
	return {
		weight: 1,
		opacity: 1,
		color: 'white',
		fillOpacity: color.alpha(),
		fillColor: color.hex()
	};
}

function highlightFeature(e) {
	var layer = e.target;

	layer.setStyle({
		weight: 3,
		color: '#666666',
	});

	if (!L.Browser.ie && !L.Browser.opera && !L.Browser.edge) {
		layer.bringToFront();
	}

	info.update(getProps(layer.feature.properties.GEOID));
}

function resetHighlight(e) {
	geojson.resetStyle(e.target);
	info.reset();
}
function onEachFeature(feature, layer) {
	props = getProps(feature.properties.GEOID);
	if (!('exclude' in props) || props.exclude) {
		return;
	}
	layer.on({
		mouseover: highlightFeature,
		mouseout: resetHighlight,
	});
}

var legend = L.control({position: 'bottomright'});

legend.onAdd = function (map) {
	var div = L.DomUtil.create('div', 'info legend'),
		labels = [],
		from, to;

	if (summary.candidates == null) {
		return div;
	}
	summary.candidates.sort(function(a, b) {
		return b.votes-a.votes;
	});

	for (var i in summary.candidates) {
		if ((summary.total_votes != 0 && (summary.candidates[i].votes) >= (summary.map_portion_thresh*summary.total_votes)) || (summary.total_votes == 0 && (summary.candidates[i].candidate in summary.priority))) {
			var name = summary.candidates[i].candidate;
			labels.push('<i style="background:' + loadColor(name) + '"></i> ' + name);
		}
	}

	div.innerHTML = labels.join('<br>');
	return div;
};

var path = window.location.pathname.split('/');
var date = path[1];
var race = path[2];

$.getJSON("/" + date + "/" + race + "/summary.json", function(data) {
	summary = data;
	if (!('colors' in summary)) {
		summary.colors = {};
	}
	if (!('regions' in summary)) {
		summary.regions = {};
	}
	if (!('candidates' in summary)) {
		summary.candidates = [];
	}
	if (!('priority' in summary)) {
		summary.priority = {};
	}
	colors = summary.colors;
	margin_scale = 0;
	for (var key in summary.regions) {
		var region = summary.regions[key];
		if (region.candidates != null && region.candidates.length > 1 && region.total_votes > 0) {
			region.candidates.sort(function(a, b) {
				return b.votes-a.votes;
			});
			var margin = (region.candidates[0].votes-region.candidates[1].votes)/region.total_votes;
			if (margin > margin_scale) {
				margin_scale = margin;
			}
		}
	}
	if (margin_scale == 0) {
		margin_scale = 1;
	}
	$.getJSON("/" + date + "/" + race + "/topo.json", function(topodata) {
		map = L.map('map');
		map.setView([0, 0], 0);

		for (key in topodata.objects) {
			if (key != "settings" && key != "roads" && key != "outline" && key != "cities") {
				geodata = topojson.feature(topodata, topodata.objects[key]);
				geojson = L.geoJson(geodata, {style: style, onEachFeature: onEachFeature, coordsToLatLng: function(coords) {
					if (coords[0] > 10000.0) {
						return [coords[0]/100000.0, coords[1]/100000.0];
					} else if (coords[0] < -10000.0) {
						return [coords[1]/100000.0, coords[0]/100000.0];
					}
					return [coords[1], coords[0]];
				}}).addTo(map);
				map.fitBounds(geojson.getBounds());
			}
		}

		info.addTo(map);
		legend.addTo(map);
		map.setMinZoom(map.getZoom()-1);
		map.setMaxZoom(map.getZoom()+4);
	});
});