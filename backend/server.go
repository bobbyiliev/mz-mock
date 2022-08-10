package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/jackc/pgproto3"
)

type PgFortuneBackend struct {
	backend   *pgproto3.Backend
	conn      net.Conn
	responder func() ([]byte, error)
}

func NewPgFortuneBackend(conn net.Conn, responder func() ([]byte, error)) *PgFortuneBackend {
	backend := pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn)

	connHandler := &PgFortuneBackend{
		backend:   backend,
		conn:      conn,
		responder: responder,
	}

	return connHandler
}

func (p *PgFortuneBackend) Run() error {
	defer p.Close()

	err := p.handleStartup()
	if err != nil {
		return err
	}

	for {
		msg, err := p.backend.Receive()
		if err != nil {
			return fmt.Errorf("error receiving message: %w", err)
		}

		switch a := msg.(type) {
		case *pgproto3.Query:
			var buf []byte
			response, err := p.generateQueryResponse(a.String)
			if err != nil {
				return err
			}

			buf = makeResponseData(response)
			_, err = p.conn.Write(buf)
			if err != nil {
				return fmt.Errorf("error writing query response: %w", err)
			}
		case *pgproto3.Terminate:
			return nil
		default:
			return fmt.Errorf("received message other than Query from client: %#v", msg)
		}
	}
}

func (p *PgFortuneBackend) generateQueryResponse(query string) (response []byte, err error) {
	query = strings.ToLower(query)
	switch query {
	case "show docs;":
		response = []byte("https://materialize.com/docs")
	case "show demos;":
		response = []byte("https://materialize.com/demos")
	case "show github;":
		response = []byte("https://github.com/MaterializeInc/materialize")
	case "show sources;":
		response = []byte("https://materialize.com/docs/sql/create-source/")
	case "welcome;":
		response = []byte("welcome!")
	case "show help;":
		response = helpString()
	default:
		response, err = p.responder()
		if err != nil {
			return response, fmt.Errorf("error generating query response: %w", err)
		}
	}

	return response, nil
}

func helpString() []byte {
	return []byte(`Available queries:
- SHOW DOCS;
- SHOW DEMOS;
- SHOW SOURCES;
- SHOW GITHUB;`)
}

func makeResponseData(message []byte) []byte {
	buf := (&pgproto3.RowDescription{Fields: []pgproto3.FieldDescription{
		{
			Name:                 []byte("Output"),
			TableOID:             0,
			TableAttributeNumber: 0,
			DataTypeOID:          25,
			DataTypeSize:         -1,
			TypeModifier:         -1,
			Format:               0,
		},
	}}).Encode(nil)
	buf = (&pgproto3.DataRow{Values: [][]byte{message}}).Encode(buf)
	buf = (&pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")}).Encode(buf)
	buf = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)

	return buf
}

func (p *PgFortuneBackend) handleStartup() error {
	startupMessage, err := p.backend.ReceiveStartupMessage()
	if err != nil {
		return fmt.Errorf("error receiving startup message: %w", err)
	}

	switch startupMessage.(type) {
	case *pgproto3.StartupMessage:
		buf := (&pgproto3.AuthenticationOk{}).Encode(nil)
		buf = (&pgproto3.NoticeResponse{Message: "\u001b[36m-------------------\x1B[0m"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "▄▄▄             ╷╷╷"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "████▄          ││││    \033[1m\u001b[36mMaterialize TUI\u001b[0m \033[0m"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: " ▀████▄     ╷│ ││││    ==============="}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "█▄ ▀████▄ ╷│││ ││││    \u001b[4mCommands\u001b[0m:"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "███▄ ▀████▄│││ ││││     -\033[1m\u001b[36m SHOW\u001b[0m HELP; \033[0m"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "▀████▄ ▀████▄│ ││││     -\033[1m\u001b[36m SHOW\u001b[0m DOCS; \033[0m"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "▄ ▀████▄ ▀████▄││││     -\033[1m\u001b[36m SHOW\u001b[0m GITHUB; \033[0m"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "██▄ ▀████▄ ▀████▄││     -\033[1m\u001b[36m SHOW\u001b[0m DEMOS; \033[0m"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "████▄ ▀████▄ ▀████▄"}).Encode(buf)
		buf = (&pgproto3.NoticeResponse{Message: "\u001b[36m-------------------\x1B[0m"}).Encode(buf)
		buf = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
		_, err = p.conn.Write(buf)
		if err != nil {
			return fmt.Errorf("error sending ready for query: %w", err)
		}
	case *pgproto3.SSLRequest:
		_, err = p.conn.Write([]byte("N"))
		if err != nil {
			return fmt.Errorf("error sending deny SSL request: %w", err)
		}
		return p.handleStartup()
	default:
		return fmt.Errorf("unknown startup message: %#v", startupMessage)
	}

	return nil
}

func (p *PgFortuneBackend) Close() error {
	return p.conn.Close()
}
