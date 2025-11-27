package minecraft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// RCON packet types.
const (
	rconTypeAuth     int32 = 3
	rconTypeCommand  int32 = 2
	rconTypeResponse int32 = 0
)

// SendRconCommand connects to the server via RCON, authenticates, and sends a command.
func SendRconCommand(host, port, password, command string) (string, error) {
	if command == "" {
		return "", fmt.Errorf("command is required")
	}

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect to RCON at %s: %w", address, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))

	if err := sendPacket(conn, 1, rconTypeAuth, password); err != nil {
		return "", fmt.Errorf("failed to send auth packet: %w", err)
	}

	authResp, err := readPacket(conn)
	if err != nil {
		return "", fmt.Errorf("failed to read auth response: %w", err)
	}
	if authResp.ID == -1 {
		return "", fmt.Errorf("rcon authentication failed (check password)")
	}

	if err := sendPacket(conn, 2, rconTypeCommand, command); err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	resp, err := readPacket(conn)
	if err != nil {
		return "", fmt.Errorf("failed to read command response: %w", err)
	}

	return resp.Body, nil
}

type rconPacket struct {
	ID   int32
	Type int32
	Body string
}

func sendPacket(w io.Writer, id, pktType int32, body string) error {
	payload := bytes.Buffer{}
	if err := binary.Write(&payload, binary.LittleEndian, id); err != nil {
		return err
	}
	if err := binary.Write(&payload, binary.LittleEndian, pktType); err != nil {
		return err
	}
	payload.WriteString(body)
	payload.WriteByte(0x00)
	payload.WriteByte(0x00)

	length := int32(payload.Len())

	header := bytes.Buffer{}
	if err := binary.Write(&header, binary.LittleEndian, length); err != nil {
		return err
	}

	if _, err := w.Write(append(header.Bytes(), payload.Bytes()...)); err != nil {
		return err
	}
	return nil
}

func readPacket(r io.Reader) (*rconPacket, error) {
	var length int32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)
	var id, pktType int32
	if err := binary.Read(buf, binary.LittleEndian, &id); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &pktType); err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	if len(bodyBytes) >= 2 {
		bodyBytes = bodyBytes[:len(bodyBytes)-2]
	}

	return &rconPacket{
		ID:   id,
		Type: pktType,
		Body: string(bodyBytes),
	}, nil
}
