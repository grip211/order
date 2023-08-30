### Orders

#### Запустить сервис

```shell
$ go run ./cmd/order/main.go --config-file config.yml
```

#### Запускаем nats

в проекте лежит экзешник

```shell
$ ./nats-server

> [14172] 2023/08/30 19:43:55.604859 [INF] Starting nats-server
> [14172] 2023/08/30 19:43:55.647513 [INF]   Version:  2.9.21
> [14172] 2023/08/30 19:43:55.647513 [INF]   Git:      [b2e7725]
> [14172] 2023/08/30 19:43:55.647513 [INF]   Name:     NBE6HFAF2SAVPRIAP7PGZUWSLNM5YARVRO5WP3AIWYGOL5M4KVIGQJVU
> [14172] 2023/08/30 19:43:55.647513 [INF]   ID:       NBE6HFAF2SAVPRIAP7PGZUWSLNM5YARVRO5WP3AIWYGOL5M4KVIGQJVU
> [14172] 2023/08/30 19:43:55.652155 [INF] Listening for client connections on 0.0.0.0:4222
> [14172] 2023/08/30 19:43:55.677002 [INF] Server is ready
```

#### Записать в очередь в Nats

```shell
curl --location --request POST 'http://localhost:3001/api/v1/publish' \
--header 'Content-Type: application/json' \
--data-raw '{
  "order_uid": "b563feb7b2b84b6tess",
  "track_number": "WBILMTESTTRACK2",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6tess",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934931,
      "track_number": "WBILMTESTTRACK2",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    },
    {
      "chrt_id": 9934932,
      "track_number": "WBILMTESTTRACK2",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}'
```

#### Записать на прямую в базу

```shell
curl --location --request POST 'http://localhost:3001/api/v1/save' \
--header 'Content-Type: application/json' \
--data-raw '{
  "order_uid": "b563feb7b2b84b6tess",
  "track_number": "WBILMTESTTRACK2",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6tess",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934931,
      "track_number": "WBILMTESTTRACK2",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    },
    {
      "chrt_id": 9934932,
      "track_number": "WBILMTESTTRACK2",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}'
```

#### Получить запись по ID

```shell
curl --location --request GET 'http://localhost:3001/api/v1/get?id=b563feb7b2b84b6tess' \
  --header 'Content-Type: application/json'
```

#### Получить все записи

```shell
curl --location --request GET 'http://localhost:3001/api/v1/all' \
  --header 'Content-Type: application/json'
```

#### Открыть в бразуере простейший интерфейс

http://localhost:3001/all