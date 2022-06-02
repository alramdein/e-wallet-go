# e-wallet-go

## Installation
1. Clone project with 
    ```
    git clone https://github.com/alramdein/e-wallet-go.git
    ```
2. On the root folder, run `go mod tidy`
3. Run the docker to serve kafka with command `docker-compose up`
4. Update di config if necessary on `config.yml` (port or etc.)
5. Run the project with 
`go run main.go` . The project should run on `localhost:<port>`

## Feature
1. GET `/wallet/:id/details`. Endpoin to see detail of a wallet
2. POST `/deposit`. Endpoint used to deposit a money to wallet  