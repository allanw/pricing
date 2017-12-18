# Pricing

## Getting the code and running it

Clone the repo and then either do `go run main.go` or: 

`go build .`
`./pricing`

or

Run `go get github.com/allanw/pricing` and cd to your $GOPATH bin dir and run `./pricing`

## Querying the API

Pass order info to the API using e.g. `curl`:

`curl -H "Content-Type: appication/json" -d '{"order": {"id": 12345, "customer": {}, "items": [{"product_id": 1, "quantity": 1}, {"product_id": 2, "quantity": 5}]}}' localhost:8080/order`

The JSON response output will be displayed in the terminal. An `orders.csv` will also be generated inside the working directory.

## Running the Python analysis script

1. Ensure you have Python 3 and `pandas` is installed
2. Run `python3 analysis.py` in the working directory
