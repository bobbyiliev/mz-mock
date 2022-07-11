# Quick introduction to the Postgres wire protocol

# Introduction

PostgreSQL is a relational database management system (RDBMS) that is designed to handle large amounts of data efficiently. It is a free and open-source software product that is developed by the PostgreSQL Global Development Group.

PostgreSQL uses a custom message-based protocol for communication between the client and the server. The protocol is supported over TCP/IP and also over Unix-domain sockets. The current version of the protocol is 3.0.

In the following, we will describe the wire protocol in detail. We will be using the `psql` command-line utility that allows you to interact with PostgreSQL databases. All that `psql` does is to give you a prompt where you can enter your queries. It then takes your query and sends it over the network to the PostgreSQL server and displays the results.

# Frontend/Backend Protocol

The protocol that the client and server use to communicate is called the `Frontend/Backend Protocol`.

A typical message in the PostgreSQL messaged-based protocol is a sequence of bytes. The first byte is the message type followed by the length message and the payload. A Regular packet would look as follows:

```
char tag | int32 len | payload
```

The very first byte of the message is the message type. The message type is a single byte that indicates the type of message. Some of the message types are:

- 'R': Authentication
- 'Q': Query
- 'X': Terminate
- 'Z': Ready for Query
- 'P': Parse
- 'B': Bind
- 'E': Error
- 'C': Command Complete
- 'D': Data Row
- 'I': Empty Query Response
- 'T': Row Description
- 'N': Notice Response
- 'A': Notification Response
- 'W': Copy Both Response
- 'K': Backend Key Data
- 'P': Parameter Description

For a detailed description of the message types, see the [PostgreSQL wire protocol](https://www.postgresql.org/docs/current/protocol-flow.html).

## Authentication

The authentication flow of data between the client and the server is as follows:

![Postgres protocol - auth flow](https://user-images.githubusercontent.com/21223421/178289606-4a0d4601-b14d-410c-887c-efa235755d55.png)

If a connection requires authentication, the backend will send a `AuthenticationRequest`, there are several authentication types that can be demanded like plain-text or MD5 password. It’s up to the server to require plain text or encrypted.

Once the authentication is complete, or if no auth is necessary, the server sends a `AuthenticationOK`.

The authentication request packet consists of the following fields:

| | | | |
|-|-|-|-|
| 'R' | int32 len | int32 method | optional other |

As Wireshark 'speaks' the PostgreSQL wire protocol, we will use it to inspect the traffic between the client and server. Let's use WireShark to inspect the `AuthenticationRequest` packet to make it easier to understand.

![Authentication request inspected with WireShark](https://user-images.githubusercontent.com/21223421/178292228-ff6ab5cd-db4b-42c7-b4cb-055d70e463af.png)

To mimic the authentication flow in NodeJS, we can re-create the AuthenticationRequest packet using the following function:

```javascript
function authOk() {
  let buf = Buffer.from("R");
  buf = Buffer.concat([buf, Buffer.from([0, 0, 0, 8])]);
  buf = Buffer.concat([buf, Buffer.from([0, 0, 0, 0])]);
  buf = readyForQuery(buf);
  return buf;
}
```

For more information on how to handle different authentication types, see the [PostgreSQL wire protocol](https://www.postgresql.org/docs/current/protocol-flow.html#id-1.10.5.7.3).

Once we have the authentication established, let's proceed to the query flow.

## Simple Query Flow

A standard query cycle is usually initiated by the frontend/client sending a `Query` message to the backend/server. The query flow is as follows:

- The client sends an SQL command, which starts with 'Q' and is followed by the SQL command.
- The server replies with `RowDescription` ('T') detailing the structure
- Each column has a name, the type OID, length, and modifier (like char(16))
- Each column is marked as containing binary or text output
- After that, a `DataRow` ('D') message is sent for every row
- Finally, the server sends `CommandComplete` ('C') and `ReadyForQuery` ('Z') which informs the client that the query is complete and the server is ready for another query.

To visualize the query flow, here is a diagram representing the complete query flow:

![Postgres protocol - simple query flow](https://user-images.githubusercontent.com/21223421/178294158-ac9d8591-7224-4480-8a3d-024e2cd80782.png)

Next, let's take a closer look at each of the messages and see how they are structured and how we can try and re-create them using NodeJS.

### The Query Message

The query frames for the above example would look like this:

```
'Q' | int32 len | char[len] query
```

Rundown of the query message packet:
- `Q`: The first byte is the message type ('Q') representing the query message.
- `int32 len`: The length of the message.
- A string containing the query string itself.

Inspected with WireShark, the query frame looks like this:

![Postgres protocol simple query packet](https://user-images.githubusercontent.com/21223421/178294620-8f46d06f-9791-4e68-b78b-7d6ae26c7394.png)

This packet sent from the client is fairly simple, it just consists of the query string. On the backend, the query is parsed and the results are sent back to the client. To do that with NodeJS, we can use the following function:

```javascript
function QueryParser(query) {
  // Check if the query packet starts with 'Q'
  if (chunk[0] === "Q".charCodeAt(0)) {
    // Remove the bynary prefix 'Q' and the bytes representing the lenght from the query:
    query = query.toString("utf8").substring(5);
    // Remove the last byte from the query:
    query = query.substring(0, query.length - 1);
    return query.trim().toUpperCase();
  }
}
```

With that, you can now use a simple `case` statement to handle the different query types:

```javascript
function HandleQuery(query) {
  switch (query) {
    case "SELECT 1;":
      values = ["1"];
      return values;
    case "SHOW DOCS;":
      values = ["https://materialize.com/docs"];
      return values;
    default:
      values = ["Hello, world!"];
      return values;
  }
}
```

Now that we have the query, we can send it to the backend and get the results.

### The `RowDescription` Message

The `RowDescription` message is sent by the server to the client to describe the columns of the result set.

The `RowDescription` message consists of the following fields:

```
'T' | int32 len | int16 numfields | str col | int32 tableoid | int16 colno | int32 typeoid | int16 typelen | int32 typmod | int16 format
```

Rundown of the `RowDescription` message packet:
- `T`: The first byte is the message type ('T') representing the `RowDescription` message.
- `int32 len`: The length of the message.
- `int16 numfields`: The number of fields in the result set.

Then, for each field, there is the following:
- `str col`: The name of the column.
- `int32 tableoid`: The object ID of the table containing the column.
- `int16 colno`: The column number of the column.
- `int32 typeoid`: The object ID of the type of the column.
- `int16 typelen`: The data type size.
- `int32 typmod`: The type modifier.
- `int16 format`: The format of the data.

Here is an example of a `RowDescription` packet inspected with WireShark:

![Postgres RowDescription packet inspected with Wireshark](https://user-images.githubusercontent.com/21223421/178294947-1adab49f-75bb-436d-9cc8-7ae17aca210b.png)

To re-create the `RowDescription` message, we can use the following function:

```javascript
function RowDescription(fields) {
  rowDescriptionType = Buffer.from("T");
  emptyLenght = Buffer.from([-1, -1, -1, -1]); // -1 means that the length of the message is not known yet.

  buf = Buffer.from([0, fields.length]);
  for (let i = 0; i < fields.length; i++) {
    buf = Buffer.concat([buf, Buffer.from(fields[i].name)]);
    buf = Buffer.concat([buf, Buffer.from([0])]);
    buf = Buffer.concat([buf, Buffer.from([0, 0, 0, fields[i].tableOID])]);
    buf = Buffer.concat([
      buf,
      Buffer.from([0, fields[i].tableAttributeNumber]),
    ]);
    buf = Buffer.concat([buf, Buffer.from([0, 0, 0, fields[i].dataTypeOID])]);
    buf = Buffer.concat([
      buf,
      Buffer.from([fields[i].dataTypeSize, fields[i].dataTypeSize]),
    ]);
    buf = Buffer.concat([
      buf,
      Buffer.from([-1, -1, -1, fields[i].typeModifier]),
    ]);
    buf = Buffer.concat([buf, Buffer.from([0, fields[i].format])]);
  }
  final = Buffer.concat([
    rowDescriptionType,
    Buffer.from([0, 0, 0, buf.length + rowDescriptionType.length + 3]),
  ]);
  final = Buffer.concat([final, buf]);

  return final;
}
```

Next, let's look into the `DataRow` message.

### The `DataRow` Message

Right after the `RowDescription` message, the server sends the `DataRow` message, which contains the actual data of the result set.

```
'D' | int32 len | int16 numfields | int32 fieldlen | char[fieldlen] data ...
```

Rundown of the `DataRow` message packet:
- `D`: The first byte is the message type ('D') representing the `DataRow` message.
- `int32 len`: The length of the message.
- `int16 numfields`: The number of fields in the result set.
- `int32 fieldlen`: The length of the column.
- `char[fieldlen] data`: The data of the column.

Inspected with WireShark, the data row frame looks as follows:

![Postgres protocol data row packet inspected](https://user-images.githubusercontent.com/21223421/178295136-cbe1223b-8ca4-4e6b-9316-0375d12ce456.png)

Depending on how many rows are in the result set, the `DataRow` message can be sent multiple times.

### The `CommandComplete` Message

The `CommandComplete` message is sent by the server to the client to indicate that the query is complete.

```
'C' | int32 len | str tag
```

Rundown of the `CommandComplete` message packet:
- `C`: The first byte is the message type ('C') representing the `CommandComplete` message.
- `int32 len`: The length of the message.
- `str tag`: The tag of the command. There can be different tags for different commands, you can find a complete list of tags in the [Postgres documentation](https://www.postgresql.org/docs/current/protocol-message-formats.html).

When inspecting the packet with WireShark, the command complete frame looks as follows:

![Postgres protocol command complete packet inspected](https://user-images.githubusercontent.com/21223421/178295697-e627a9cd-074c-47b7-abb5-dbbb70d0442f.png)

To re-create the `CommandComplete` message, we can use the following function:

```javascript
function CommandComplete(buf) {
  buf = Buffer.concat([buf, Buffer.from("C")]);
  buf = Buffer.concat([buf, Buffer.from([0, 0, 0, 13])]);
  buf = Buffer.concat([buf, Buffer.from("SELECT 1")]);
  buf = Buffer.concat([buf, Buffer.from([0])]);
  return buf;
}
```

Finally, let's look into the `ReadyForQuery` message.

### The `ReadyForQuery` Message

The `ReadyForQuery` message is sent by the server to the client to indicate the server is ready for another query. This is essentially the end of the query simple query flow.

The `ReadyForQuery` message consists of the following fields:

```
'Z' | int32 len | 'I' or 'T' or 'E'
```

Rundown of the `ReadyForQuery` message packet:
- `Z`: The first byte is the message type ('Z') representing the `ReadyForQuery` message.
- `int32 len`: The length of the message.
- `'I' or 'T' or 'E'`: The state of the connection. The state can be one of the following:
  - `I`: Idle. The connection is waiting for a query.
  - `T`: In transaction. The connection is in a transaction.
  - `E`: Failed. The connection is in a failed state.

Inspected with WireShark, the ready for query frame looks as follows:

![Postgres protocol ready for query packet inspected](https://user-images.githubusercontent.com/21223421/178295926-3ef33f78-796f-4850-b866-dea4d3b66b45.png)

To re-create the `ReadyForQuery` message, we can use the following function:

```javascript
function readyForQuery(ready) {
  ready = Buffer.concat([ready, Buffer.from("Z")]);
  ready = Buffer.concat([ready, Buffer.from([0, 0, 0, 5])]);
  ready = Buffer.concat([ready, Buffer.from("I")]);
  return ready;
}
```

# Putting it all together

To put this all together, we can use the `net` NodeJS module which provides us with an asynchronous network API for creating stream-based servers.

```javascript
const Net = require("net");
const { type } = require("os");

// The port on which the server is listening.
const port = 5432;

// Create a new TCP server.
const server = new Net.Server();

// The server listens to a socket for a client to make a connection request.
server.listen(port, function () {
  console.log(
    `Server listening for connection requests on socket localhost:${port}`
  );
});

// When a client requests a connection with the server, the server creates a new socket dedicated to that client.
server.on("connection", function (socket) {
  console.log("A new connection has been established.");
  // Send an authentication OK message to the client
  buf = authOk();
  socket.write(buf);

  // Basic query flow
  socket.on("data", function (chunk) {
    // If starts with 'Q', it's a query:
    if (chunk[0] === "Q".charCodeAt(0)) {
      // Values for the RowDescription message:
      const fields = [
        { name: "Output", tableOID: 0, tableAttributeNumber: 0, dataTypeOID: 25, dataTypeSize: -1, typeModifier: -1, format: 0 },
      ];

      // Prepare the RowDescription message:
      let buf = RowDescription(fields);

      // Parse the query to string:
      let query = QueryParser(chunk);

      // Case statement to check the query and return the correct response:
      let values = HandleQuery(query)

      // Prepare the DataRow message concatenated with the RowDescription message:
      buf = Buffer.concat([buf, DataRow(values)]);

      // Prepare the CommandComplete message:
      buf = Buffer.concat([buf, CommandComplete(buf)]);

      // readyForQuery concatenated with the message:
      buf = Buffer.concat([buf, readyForQuery(buf)]);

      socket.write(buf);
    } else {
      // Ready for Query:
      buf = readyForQuery(buf);
    }
  });

  socket.on("end", function () {
    console.log("Closing connection with the client");
  });

  socket.on("error", function (err) {
    console.log(`Error: ${err}`);
  });
});
```

You can find the complete source code of the server in the `server.js` file in the following repository:

[Materialize Mock Service - NodeJS example](https://github.com/bobbyiliev/mz-mock/tree/main/nodejs-mock)

To run the server, you can use the following command:

```bash
node server.js
```

This will start the server on port `5432`. After you run the server, you can connect to it using the following command:

```bash
psql -h localhost -p 5432 -U postgres
```

As a next step, you can try adding more queries to the `HandleQuery` function and see how the server responds to them.

A great way to get a better understanding of the protocol is to install Wireshark and inspect the packets sent by the server and the client!

# Conclusion

We've covered only the basics of the Postgres protocol. For more information, see the [Postgres wire protocol documentation](https://www.postgresql.org/docs/current/protocol-flow.html). The protocol is well documented. 

Another great resource is the following implementations of the Postgres wire protocol using Golang:

[Package pgproto3 is an encoder and decoder of the PostgreSQL wire protocol version 3](https://github.com/jackc/pgproto3/).

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


## References

- [A look at the PostgreSQL wire protocol](https://www.pgcon.org/2014/schedule/attachments/330_postgres-for-the-wire.pdf)
- [How does PostgreSQL actully work](https://www.youtube.com/watch?v=OeKbL55OyL0)
- [Postgres wire protocol](https://www.postgresql.org/docs/current/protocol-flow.html)

## Useful links

* [Materialize Cloud](https://materialize.com/cloud?utm_source=bobbyiliev)
* [Official Materialize documentation](https://materialize.com/docs?utm_source=bobbyiliev)
* [Materialize GitHub repository](https://github.com/MaterializeInc/Materialize?utm_source=bobbyiliev)
* [Materialize website](https://materialize.com?utm_source=bobbyiliev)
* [Official Materialize demos](https://materialize.com/docs/demos?utm_source=bobbyiliev)
* [Materialize white paper](https://materialize.com/resources/materialize-an-overview?utm_source=bobbyiliev)