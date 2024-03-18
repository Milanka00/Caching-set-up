vcl 4.1;

import directors;
import std;

backend backend1 {
    .host = "envoy_new";
    .port = "10000";
}

sub vcl_init {
    new vdir = directors.round_robin();
    vdir.add_backend(backend1);
}

sub vcl_recv {
    set req.backend_hint = backend1;
    # std.log("Host header received: " + req.http.Host);
}

sub vcl_backend_fetch {
    if (bereq.method == "GET") {
        unset bereq.body;
    }
    return (fetch);
}

sub vcl_miss {
    set req.http.X-Cache-Miss = "1";
    return (pass);
}

sub vcl_deliver {
    if (obj.hits > 0) {
        set resp.http.X-Cache-Host = req.http.Host;
        set resp.http.X-Cache-Info = "Cached under host: " + req.http.Host + "; Request URI: " + req.url;
    }
}
