import requests
import json


def main():
    url = "http://localhost:7777"

    # Example  method
    payload = {
        "method": "Server.PullRequest",
        "params": [{"Comment": "Hello Comments"}],
        "jsonrpc": "2.0",
        "id": 0,
    }
    response = requests.post(url, json=payload).json()

    print(response["result"])
    #assert response["jsonrpc"]
    assert response["id"] == 0

if __name__ == "__main__":
    main()