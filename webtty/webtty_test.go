package webtty

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"testing"
)

type pipePair struct {
	*io.PipeReader
	*io.PipeWriter
}

// testSlave implements Slave over pipes. Reads return what is written
// to the slave's output pipe (simulating PTY output); writes go to the
// slave's input pipe (simulating PTY input).
type testSlave struct {
	*io.PipeReader
	*io.PipeWriter
}

func (testSlave) WindowTitleVariables() map[string]interface{} { return nil }
func (testSlave) ResizeTerminal(columns int, rows int) error   { return nil }

// runWebTTY starts Run in a goroutine and returns a channel that
// receives its result once the given pipes are closed and the context
// is canceled.
func runWebTTY(ctx context.Context, dt *WebTTY) <-chan error {
	runErr := make(chan error, 1)
	go func() {
		runErr <- dt.Run(ctx)
	}()
	return runErr
}

func TestWriteFromPTY(t *testing.T) {
	connInPipeReader, connInPipeWriter := io.Pipe()   // in to conn
	connOutPipeReader, connOutPipeWriter := io.Pipe() // out from conn

	conn := pipePair{
		connOutPipeReader,
		connInPipeWriter,
	}

	slaveOutPipeReader, slaveOutPipeWriter := io.Pipe() // out from slave
	slaveInPipeReader, slaveInPipeWriter := io.Pipe()   // in to slave
	_ = slaveInPipeReader
	slave := testSlave{
		slaveOutPipeReader,
		slaveInPipeWriter,
	}

	dt, err := New(conn, slave)
	if err != nil {
		t.Fatalf("Unexpected error from New(): %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	runErr := runWebTTY(ctx, dt)

	buf := make([]byte, 1024)

	// Run sends the initialize message before anything else.
	n, err := connInPipeReader.Read(buf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if buf[0] != SetWindowTitle {
		t.Fatalf("Unexpected message type `%c`", buf[0])
	}

	message := []byte("foobar")
	n, err = slaveOutPipeWriter.Write(message)
	if err != nil {
		t.Fatalf("Unexpected error from Write(): %s", err)
	}
	if n != len(message) {
		t.Fatalf("Write() accepted `%d` for message `%s`", n, message)
	}

	n, err = connInPipeReader.Read(buf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if buf[0] != Output {
		t.Fatalf("Unexpected message type `%c`", buf[0])
	}
	decoded := make([]byte, 1024)
	n, err = base64.StdEncoding.Decode(decoded, buf[1:n])
	if err != nil {
		t.Fatalf("Unexpected error from Decode(): %s", err)
	}
	if !bytes.Equal(decoded[:n], message) {
		t.Fatalf("Unexpected message received: `%s`", decoded[:n])
	}

	cancel()
	connOutPipeWriter.Close()
	slaveOutPipeWriter.Close()
	if err := <-runErr; err == nil {
		t.Fatal("Expected an error from Run() after shutdown")
	}
}

func TestWriteFromConn(t *testing.T) {
	connInPipeReader, connInPipeWriter := io.Pipe()   // in to conn
	connOutPipeReader, connOutPipeWriter := io.Pipe() // out from conn

	conn := pipePair{
		connOutPipeReader,
		connInPipeWriter,
	}

	slaveOutPipeReader, slaveOutPipeWriter := io.Pipe() // out from slave
	slaveInPipeReader, slaveInPipeWriter := io.Pipe()   // in to slave
	slave := testSlave{
		slaveOutPipeReader,
		slaveInPipeWriter,
	}

	dt, err := New(conn, slave, WithPermitWrite())
	if err != nil {
		t.Fatalf("Unexpected error from New(): %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	runErr := runWebTTY(ctx, dt)

	var (
		message []byte
		n       int
	)
	readBuf := make([]byte, 1024)

	// Run sends the initialize message before anything else.
	n, err = connInPipeReader.Read(readBuf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if readBuf[0] != SetWindowTitle {
		t.Fatalf("Unexpected message type `%c`", readBuf[0])
	}

	// input
	message = []byte("1hello\n") // Input message
	n, err = connOutPipeWriter.Write(message)
	if err != nil {
		t.Fatalf("Unexpected error from Write(): %s", err)
	}
	if n != len(message) {
		t.Fatalf("Write() accepted `%d` for message `%s`", n, message)
	}

	n, err = slaveInPipeReader.Read(readBuf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if !bytes.Equal(readBuf[:n], message[1:]) {
		t.Fatalf("Unexpected message received: `%s`", readBuf[:n])
	}

	// ping
	message = []byte("2") // Ping message
	n, err = connOutPipeWriter.Write(message)
	if err != nil {
		t.Fatalf("Unexpected error from Write(): %s", err)
	}
	if n != len(message) {
		t.Fatalf("Write() accepted `%d` for message `%s`", n, message)
	}

	n, err = connInPipeReader.Read(readBuf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if !bytes.Equal(readBuf[:n], []byte{Pong}) {
		t.Fatalf("Unexpected message received: `%s`", readBuf[:n])
	}

	// TODO: resize

	cancel()
	connOutPipeWriter.Close()
	slaveOutPipeWriter.Close()
	if err := <-runErr; err == nil {
		t.Fatal("Expected an error from Run() after shutdown")
	}
}
