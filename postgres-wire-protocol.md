# Quick introduction to the Postgres wire protocol

## Introduction

PostgreSQL is a relational database management system (RDBMS) that is designed to handle large amounts of data efficiently. It is a free and open source software product that is developed by the PostgreSQL Global Development Group.

PostgreSQL uses a custom message-based protocol for communication between the client and the server. The protocol is supported over TCP/IP and also over Unix-domain sockets. The current version of the protocol is 3.0.

In the following, we will describe the wire protocol in detail. We will be using the `psql` is a command line utility that allows you to interact with PostgreSQL databases. All that `psql` does is to give you a prompt where you can enter your queries. It then takes your query and sends it over the network to the PostgreSQL server and displays the results.

## Frontend/Backend Protocol

The protocol that the client and server use to communicate is called the `Frontend/Backend Protocol`.

A typical message in the PostgreSQL messaged-based protocol is a sequence of bytes. The first byte is the message type followed by the lenght message and the payload. A Regular packet would look as follows:

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

### Authentication

The authentication flow of data between the client and the server is as follows:

![](https://imgur.com/6vGdq40.png)

If a connection requires authentication, the backend will send a `AuthenticationRequest`, there are several authentication types that can be demanded like plain-text or MD5 password. It’s up to the server to require plain text or encrypted.

Once the authentication is complete, or if no auth is necessary, the server sends a `AuthenticationOK`.

### Query

A standard query cycle is usually initiated by the frontened/client sending a `Query` message to the backend/server. The query flow is as follows:

- The client sends an SQL command, which starts with 'Q' and is followed by the SQL command.
- The server replies with `RowDescription` ('T') detailing the structure
- Each column has a name, the type OID, length and modifier (like char(16))
- Each column is marked as containing binary or text output
- After that a `DataRow` ('D') message is sent for every row
- Finally, the server sends `CommandComplete` ('C') and `ReadyForQuery` ('Z') which informs the client that the query is complete and the server is ready for another query.

Here is a diagram of the query flow:

![](https://imgur.com/2s7zJ5l.png)

The query frames for the above example would look like:

- Query:

```
'Q' | int32 len | char[len] query
```

- RowDescription:

```
'T' | int32 len | int16 numfields | str col | int32 tableoid | int16 colno | int32 typeoid | int16 typelen | int32 typmod | int16 format
```

- DataRow:

```
'D' | int32 len | int16 numfields | int32 fieldlen | char[fieldlen] data ...
```

- CommandComplete:

```
'C' | int32 len | str tag
```

- ReadyForQuery:

```
'Z' | int32 len | 'I' or 'T' or 'E'
```

This is all well documented in the [PostgreSQL wire protocol](https://www.postgresql.org/docs/current/protocol-flow.html).

## Conclusion

For more information, about the Postgres wire protocol, see the [Postgres wire protocol documentation](https://www.postgresql.org/docs/current/protocol-flow.html).

You can take a look at the following implementations of the Postgres wire protocol using Golang:

[Package pgproto3 is a encoder and decoder of the PostgreSQL wire protocol version 3.](https://github.com/jackc/pgproto3/).

## References

- [A look at the PostgreSQL wire protocol](https://www.pgcon.org/2014/schedule/attachments/330_postgres-for-the-wire.pdf)
- [How does PostgreSQL actully work](https://www.youtube.com/watch?v=OeKbL55OyL0)