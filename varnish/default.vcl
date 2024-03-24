vcl 4.1;
import std;

backend default {
    .host = "envoy_new";
    .port = "9096";
}

sub vcl_recv {
     set req.http.Host = req.http.redirect-backend;
     set req.url = req.http.x-temp-path;
     set req.url = std.querysort(req.url);
}


sub vcl_miss {
    if (req.http.x-cluster-header == "varnish") {
        set req.http.x-cluster-header = req.http.redirect-backend;  
    } 
        return (fetch);
    
}

sub vcl_backend_fetch {
    if (bereq.http.x-temp-path) {
        set bereq.url = bereq.http.x-temp-path;
        set bereq.url = std.querysort(bereq.url);
    }
}

sub vcl_backend_response {
    # Don't cache 404 responses
    if (beresp.status == 404) {
        set beresp.uncacheable = true;
    }

    if (bereq.url ~ "/org1/") {
        set beresp.storage = storage.org1;
        set beresp.http.x-storage = "org1";
    } elsif (bereq.url ~ "/org2/") {
        set beresp.storage = storage.org2;
        set beresp.http.x-storage = "org2";
    } else {
        set beresp.storage = storage.default; // Default storage
        set beresp.http.x-storage = "default";
    }
}



sub vcl_deliver {
    if (obj.hits > 0) {
        set resp.http.X-Cached-By = "Varnish";
        set resp.http.X-Cache-Info = "Cached under host: " + req.http.Host + "; Request URI: " + req.url;
    }
}



