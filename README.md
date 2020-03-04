# Go Redcoins API

This API simulates a bitcoin exchange platform in which it is possible to create an account and make bitcoin transactions (buy and sell) with it.

## Usage
### Creating an user
Make a **POST** request to `/users` with the following payload format:

```json
{
    "name": "Victor Moura",
    "email": "victor@email.com",
    "password": "password",
    "birth_date": "2020-03-04T12:47:29.001196Z" // Optional parameter
}
```

### Authentication
All of the endpoints, except for user creation, require JWT authentication. To authenticate, make a **POST** request to `/login` with the following payload format:

```json
{
    "email": "victor@email.com",
    "password": "password",
}
```
If the password and email are correct, you should receive an authentication token:

```
200 ok
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE1ODMzMzAwNjcsInVzZXJfaWQiOjN9.YAtBvRmAllAYN_92GerTuaUeRU_FwxnIJEPk3B3i6fg"
```

Insert this authentication token in your requests header as:

**Authorization**: `<your token goes here>`

In cURL, it would look like:

```bash
curl -i -H "Authorization: <your token goes here>" -H "Content-Type: application/json" http://localhost:8080/users
```

### Buying and Selling Bitcoins
To buy or sell bitcoins, make a **POST** request to `/transactions` with the following payload format:

```json
{
    "bc_value": 1, // bitcoin amount or -1
    "usd_value": -1, // usd amount or -1
    "type": "sell" // sell or buy
}
```

**Examples:**

- Buying 100USD in bitcoins:
```json
{
    "bc_value": -1, // bitcoin amount or -1
    "usd_value": 100, // usd amount or -1
    "type": "buy" // sell or buy
}
```

- Buying 1 bitcoin:
```json
{
    "bc_value": 1, // bitcoin amount or -1
    "usd_value": -1, // usd amount or -1
    "type": "buy" // sell or buy
}
```

- Selling 1 bitcoin:
```json
{
    "bc_value": 1, // bitcoin amount or -1
    "usd_value": -1, // usd amount or -1
    "type": "sell" // sell or buy
}
```

- Selling 100USD in bitcoins:
```json
{
    "bc_value": -1, // bitcoin amount or -1
    "usd_value": 100, // usd amount or -1
    "type": "sell" // sell or buy
}
```

Keep in mind that only one of `bc_value` or `usd_value` should be set as -1. If not, an error should pop.

### Listing all transactions from a user
To list all transactions from a user, you should first get the user ID. A **GET** `/users` will list all users. Look for the one that interests you and use its ID to make a **GET** request to `/users/{user id}/transactions`. You should get something as:

```json
// 200 ok
// http://localhost:8080/users/3/transactions

[
  {
    "id": 3,
    "type": "sell",
    "bc_value": 1,
    "usd_value": 8752.97327347,
    "owner_id": 3,
    "created_at": "2020-03-04T12:55:00.919413Z",
    "updated_at": "2020-03-04T12:55:00.919413Z"
  },
  {
    "id": 4,
    "type": "sell",
    "bc_value": 1,
    "usd_value": 8752.97327347,
    "owner_id": 3,
    "created_at": "2020-03-04T12:55:02.785985Z",
    "updated_at": "2020-03-04T12:55:02.785985Z"
  }
]
```

### Listing all transactions in a day
To list all transactions that happened in a day, make a **POST** request to `/transactions/by_day` with the following payload format:

```json
{
    "date": "2020-03-04" // Date in YYYY-MM-DD format
}
```

With the payload above, a possible result is:

```json
[
  {
    "id": 1,
    "type": "buy",
    "bc_value": 1,
    "usd_value": 8752.97327347,
    "owner_id": 1,
    "created_at": "2020-03-04T12:47:29.001196Z",
    "updated_at": "2020-03-04T12:47:29.001196Z"
  },
  {
    "id": 2,
    "type": "buy",
    "bc_value": 1.1424689288506913,
    "usd_value": 10000,
    "owner_id": 2,
    "created_at": "2020-03-04T12:47:29.17732Z",
    "updated_at": "2020-03-04T12:47:29.17732Z"
  }
]
```

## Running the application
This application was built to work within a Docker environment. To execute the whole stack, run:

```bash
sudo docker-compose up
```

Keep in mind that there are some environment variables, as `COIN_MARKET_CAP_API_TOKEN`, that must be defined in `etc/env/api.env`, located in the root of this project.

## Bitcoin Price Update Policy
This application is configured to fetch the latest Bitcoin value from [Coin Market Cap](https://coinmarketcap.com/) every 10 minutes through a cronjob (And on application startup). If, due to multiple failures, the latest price is outdated by more than 60 minutes, no transactions will be created until a more updated price is fetched.