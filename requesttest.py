import requests

if __name__ == "__main__":
    r = requests.post("http://localhost:8000/v1/secrets/", data={"secret": "MyVeryOwnsecret", "expireAfterViews": "5","expireAfter": "10"})
    print(r.content)
