# Answers every request with 200 {} so a Backend's boot-time checks against
# other services (Univents → Payssage wallet) pass without a real service.
import http.server, sys

class H(http.server.BaseHTTPRequestHandler):
    def _ok(self):
        body = b"{}"
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)
    do_GET = do_POST = do_PUT = do_PATCH = do_DELETE = _ok
    def log_message(self, fmt, *args):
        sys.stderr.write("[stub] " + (fmt % args) + "\n")

http.server.ThreadingHTTPServer(("0.0.0.0", int(sys.argv[1])), H).serve_forever()
