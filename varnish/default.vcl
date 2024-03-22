vcl 4.1;
import std;

backend default {
    .host = "envoy_new";
    .port = "9096";
}


sub vcl_miss {
    if (req.http.x-cluster-header == "varnish") {
        set req.http.x-cluster-header = req.http.redirect-backend;  
    } 
        return (fetch);
    
}

sub vcl_backend_response {
    # Don't cache 404 responses
    if (beresp.status == 404) {
        set beresp.uncacheable = true;
    }
}

sub vcl_deliver {
    if (obj.hits > 0) {
        set resp.http.X-Cached-By = "Varnish";
        set resp.http.X-Cache-Info = "Cached under host: " + req.http.Host + "; Request URI: " + req.url;
    }
}



