vcl 4.1;

import directors;
import std;

backend backend1 {
    .host = "envoy_new";
    .port = "9095";
}

sub vcl_init {
    new vdir = directors.round_robin();
    vdir.add_backend(backend1);
}

sub vcl_recv {
    set req.backend_hint = vdir.backend();
    if (req.http.x-cluster-header) {
        # X-Cluster-Header is present with any value
        std.syslog(4, "Hashing request for caching: " + req.url);
        return (hash);
    } else {
        return (pass);
    }
}

sub vcl_backend_fetch {
    if (bereq.method == "GET") {
        unset bereq.body;
    }
    std.syslog(4, "Fetching from backend: " + bereq.url);
    return (fetch);
}

sub vcl_miss {
    if (req.http.x-cluster-header) {
        std.syslog(4, "Cache miss for: " + req.url);
        set req.http.X-Cache-Miss = "1";  # Set custom header for cache miss
        return (pass);  # Directly deliver the response
    } else {
        return (fetch);  # Proceed with fetching the response from the backend
    }
}

sub vcl_deliver {
    std.syslog(4, "Delivering response for: " + req.url);
    if (obj.hits > 0) {
        set resp.http.X-Cache-Host = req.http.Host;
        set resp.http.X-Cache-Info = "Cached under host: " + req.http.Host + "; Request URI: " + req.url;
    }

    if (resp.http.X-Cache-Miss == "1") {
        std.syslog(4, "Setting X-Cluster-Header to actual_backend");
        # Change the X-Cluster-Header to actual_backend
        set resp.http.x-cluster-header = "actual_backend";
    }
}
