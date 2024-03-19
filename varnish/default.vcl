vcl 4.1;
import std;

backend default {
    .host = "envoy_new";
    .port = "9095";
}


sub vcl_miss {
    if (req.http.x-cluster-header == "varnish") {
        set req.http.x-cluster-header = "actual_backend";  
    } 
        return (fetch);
    
}

sub vcl_deliver {
   
        set resp.http.X-Cached-By = "Varnish";
    
}


