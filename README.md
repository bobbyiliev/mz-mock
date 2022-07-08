# Materialize Mock Service

This project utilizes the PostgreSQL Wire Protocol to generate a mock service for Materialize. The project uses the [pgproto3](https://pkg.go.dev/github.com/jackc/pgproto3) library as the encoder and decoder of the PostgreSQL wire protocol version 3.

## Demo

To access the demo use the following `psql` command:

```bash
psql -U materialize -h psql.bobby.sh -p 6875
```

Output:

```sql
▄▄▄             ╷╷╷
████▄          ││││    Materialize TUI
 ▀████▄     ╷│ ││││    ===============
█▄ ▀████▄ ╷│││ ││││    Commands:
███▄ ▀████▄│││ ││││     - SHOW HELP;
▀████▄ ▀████▄│ ││││     - SHOW DOCS;
▄ ▀████▄ ▀████▄││││     - SHOW GITHUB;
██▄ ▀████▄ ▀████▄││     - SHOW DEMOS;
████▄ ▀████▄ ▀████▄

materialize=> SHOW DOCS;
```

You can also access the URL directly via your browser:

https://psql.bobby.sh/

Preview:

<img width="676" alt="image" src="https://user-images.githubusercontent.com/21223421/178000533-9ebdf29e-c321-4ef8-a273-f9c687e1e228.png">

## Running the demo

To run the demo locally, you can use the following command:

```bash
docker-compose up -d
```

## Useful links

* [Materialize Cloud](https://materialize.com/cloud?utm_source=bobbyiliev)
* [Official Materialize documentation](https://materialize.com/docs?utm_source=bobbyiliev)
* [Materialize GitHub repository](https://github.com/MaterializeInc/Materialize?utm_source=bobbyiliev)
* [Materialize website](https://materialize.com?utm_source=bobbyiliev)
* [Official Materialize demos](https://materialize.com/docs/demos?utm_source=bobbyiliev)
* [Materialize white paper](https://materialize.com/resources/materialize-an-overview?utm_source=bobbyiliev)