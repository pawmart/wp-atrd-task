from http.server import HTTPServer
from server import MyServer
from os import path
from db_handler import DbHandler


if __name__ == "__main__":
    if not path.exists("task_db.db"):
        DbHandler.initialize()
    server_address = ("localhost", 8000)
    server = HTTPServer(server_address, MyServer)
    print("Starting http server on Localhost:8000")
    server.serve_forever()
