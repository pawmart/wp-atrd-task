from http.server import BaseHTTPRequestHandler
from db_handler import DbHandler
from json_parser import MyParser
from post_handler import PostHandler


class MyServer(BaseHTTPRequestHandler):

    def set_response(self, code, type_of, message):
        self.send_response(code)
        self.send_header("Content-type", "application/json")
        self.end_headers()
        if type_of == "message":
            to_print = MyParser.parse_message_to_json(message)
        else:
            to_print = MyParser.parse_to_json(message)
        self.wfile.write(bytes(to_print, "utf-8"))

    def do_GET(self):
        # ignoring request for icon while using browser
        if "favicon.ico" not in self.path:
            path = self.path
            if "/v1/secret/" in path:
                requested_hash = path.lstrip("/v1/secret/")
                secret = DbHandler.get_secret(requested_hash)
                if not secret:
                    self.set_response(404, "message", "Secret not found")
                else:
                    self.set_response(200, "data", secret)
            else:
                self.set_response(405, "message", "You've reached wrong endpoint.")

    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length).decode('UTF-8')
        if not PostHandler.check_validity(post_data):
            self.set_response(405, "message", "Wrong input")
        else:
            values = PostHandler.unpack_values(post_data)
            secret_data = DbHandler.post_secret(values["secret"], values["expireAfter"], values["expireAfterViews"])
            self.set_response(200, "data", secret_data)
