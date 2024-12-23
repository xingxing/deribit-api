package models

import (
	"context"
	"errors"
	"io"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// ObjectStream is a jsonrpc2.ObjectStream that uses a WebSocket to
// send and receive JSON-RPC 2.0 objects.
type ObjectStream struct {
	conn *websocket.Conn
}

// NewObjectStream creates a new jsonrpc2.ObjectStream for sending and
// receiving JSON-RPC 2.0 objects over a WebSocket.
func NewObjectStream(conn *websocket.Conn) ObjectStream {
	return ObjectStream{conn: conn}
}

// WriteObject implements jsonrpc2.ObjectStream.
func (t ObjectStream) WriteObject(obj interface{}) error {
	return wsjson.Write(context.Background(), t.conn, obj)
}

// ReadObject implements jsonrpc2.ObjectStream.
func (t ObjectStream) ReadObject(v interface{}) error {
	err := wsjson.Read(context.Background(), t.conn, v)
	var e *websocket.CloseError
	if errors.As(err, &e) {
		if e.Code == websocket.StatusNormalClosure && e.Error() == io.ErrUnexpectedEOF.Error() {
			// unwrapping this error.
			err = io.ErrUnexpectedEOF
		}
	}
	return err
}

// Close implements jsonrpc2.ObjectStream.
func (t ObjectStream) Close() error {
	return t.conn.Close(websocket.StatusNormalClosure, "")
}
