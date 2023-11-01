//需要修改
var host_domain = 'http://192.168.0.181:8080'

var _maq = _maq || [];
_maq.push(['_setAccount', 'test']);
(function() {
	var ma = document.createElement('script');
	ma.type = 'text/javascript';
	ma.async = true;
	ma.src = host_domain + '/assets/stat.js';
	var s = document.getElementsByTagName('script')[0];
	s.parentNode.insertBefore(ma, s);
})();
