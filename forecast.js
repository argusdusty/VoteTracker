var gammaln = function gammaln(x) {
	if(x < 0) return NaN;
	if(x == 0) return Infinity;
	if(!isFinite(x)) return x;

	var lnSqrt2PI = 0.91893853320467274178;
	var gamma_series = [76.18009172947146, -86.50532032941677, 24.01409824083091, -1.231739572450155, 0.1208650973866179e-2, -0.5395239384953e-5];
	var denom;
	var x1;
	var series;

	// Lanczos method
	denom = x+1;
	x1 = x + 5.5;
	series = 1.000000000190015;
	for(var i = 0; i < 6; i++) {
		series += gamma_series[i] / denom;
		denom += 1.0;
	}
	return( lnSqrt2PI + (x+0.5)*Math.log(x1) - x1 + Math.log(series/x) );
};
var data = [];
var x_data = [];

var N = 10000;

for (var i = 0; i <= N; i++) {
	x_data.push((i*100)/N);
}

var path = window.location.pathname.split('/');
var date = path[1];
var race = path[2];

$.getJSON("/" + date + "/" + race + "/summary.json", function(summary) {
	var totalC = 0.0;
	for (var key in summary.forecast) {
		console.log(key, summary.forecast, summary.forecast[key]);
		totalC += summary.forecast[key].concentration_param;
	}
	var min = N;
	var max = 0;
	for (var key in summary.forecast) {
		if (summary.forecast[key].concentration_param/totalC < summary.graph_portion_thresh || summary.forecast[key].odds < summary.graph_odds_thresh) {
			continue;
		}
		var alpha = summary.forecast[key].concentration_param;
		var beta = totalC-summary.forecast[key].concentration_param;
		var gamma = gammaln(alpha+beta)-gammaln(alpha)-gammaln(beta);
		var tmp_data = [];
		for (var i = 0; i <= N; i++) {
			var s = Math.exp(Math.log(i/N)*(alpha-1)+Math.log((N-i)/N)*(beta-1)+gamma);
			tmp_data.push(Math.exp(Math.log(i/N)*(alpha-1)+Math.log((N-i)/N)*(beta-1)+gamma));
			if (i > 0 && tmp_data[tmp_data.length-1] > 1e-10 && tmp_data[tmp_data.length-2] < 1e-10 && i < min) {
				min = i;
			}
			if (i > 0 && tmp_data[tmp_data.length-1] < 1e-10 && tmp_data[tmp_data.length-2] > 1e-10 && i > max) {
				max = i;
			}
		}
		data.push({x:x_data, y:tmp_data, mode:'lines', name: summary.forecast[key].candidate, line: {width: 4, color: summary.colors[summary.forecast[key].candidate]}});
	}
	for (var i in data) {
		data[i].x = data[i].x.slice(min, max);
		data[i].y = data[i].y.slice(min, max);
	}
	Plotly.newPlot('forecast', data.reverse(), {margin: {r: 10, t: 30, b: 30, l: 10}, legend: {x: 0.92, y: 0.99}, title: summary.name + ' Forecast', yaxis: {showline: false, showgrid: false, showticklabels: false}, xaxis: {title: 'Percent of vote'}}, {displayModeBar: false});
});