# e-wallet-go

## Installation
1. Clone project with 
    ```
    git clone https://github.com/alramdein/e-wallet-go.git
    ```
2. On the root folder, run `go mod tidy`
3. Run the docker-compose to serve the kafka with command `docker-compose up`
4. Update the config if needed on `config.yml` (like port or etc.)
5. Run the project with 
`go run main.go server` . It'll run both of the web server and the kafka consumer (more about this on the **Note** section below). The project should run on `http://localhost:<port>`. Default port is `http://localhost:1323` 

## Feature
1. GET `/wallet/:id/details`. Endpoin to get detail of a wallet.
2. POST `/deposit`. Endpoint used to deposit a money to wallet. **Use form URL encoded body request**. Ex: 
    ```
    wallet_id: 1 
    amount: 10
    ```

## Test Purpose
I have initialize the wallet hard-codedly for test purpose. You can change it on `cmd/server.go` function `initTestwallet()`

## Note 
Initially, I want to seperate the web server and the consumer process. Thus I use [Cobra Command](https://github.com/spf13/cobra) as you can see my on my project structre. But in the middile of the process I changed my mind and feel using only 1 process for web server and consumer is much more simpler. But unfortunatelly I'm running out of time to tidying up the project structure back to normal again. So that's why the project structure still looks like I'm using Cobra for running multiple process seperately.
