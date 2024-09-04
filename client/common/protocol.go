package common


import (
	"encoding/binary"
	"net"
)

type Bet struct {
	Name       string
	Surname    string
	Document   int
	Birthdate  string
	Number     int
}

const (
    BET int = 1
    OK  int = 2
	ERROR int = 3
)
type Protocol struct {
	conn net.Conn
}

// Create a Protocol object
func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{
		conn: conn,
	}
}

func htonl(value int) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(value))
	return bytes
}

func ntohl(value []byte) uint32 {
	return binary.BigEndian.Uint32(value)
}

func readFully(conn net.Conn, buf []byte) (int, error) {
    totalRead := 0
    for totalRead < len(buf) {
        n, err := conn.Read(buf[totalRead:])
        if err != nil {
            return totalRead, err
        }
        totalRead += n
    }
    return totalRead, nil
}

func writeFully(conn net.Conn, buf []byte) error {
    totalWritten := 0
    for totalWritten < len(buf) {
        n, err := conn.Write(buf[totalWritten:])
        if err != nil {
            return err
        }
        totalWritten += n
    }
    return nil
}


// SerializeBet serializes a Bet structure into a byte slice and sends it
func (p *Protocol) sendBet(bet Bet) (string, error) {

	// Serialize the bet type
	betBytes := htonl(int(BET))
	err := writeFully(p.conn, betBytes)
	if err != nil {
		return "Failed to send bet type", err
	}

	// Serialize the name
	nameBytes := []byte(bet.Name)
	nameSizeBytes := htonl(len(nameBytes))
	err = writeFully(p.conn, nameSizeBytes)
	if err != nil {
		return "Failed to send name size", err
	}
	err = writeFully(p.conn, nameBytes)
	if err != nil {
		return "Failed to send name", err
	}

	// Serialize the surname
	surnameBytes := []byte(bet.Surname)
	surnameSizeBytes := htonl(len(surnameBytes))
	err = writeFully(p.conn, surnameSizeBytes)
	if err != nil {
		return "Failed to send surname size", err
	}
	err = writeFully(p.conn, surnameBytes)
	if err != nil {
		return "Failed to send surname", err
	}


	// Serialize the document
	documentBytes := htonl(int(bet.Document))
	err = writeFully(p.conn, documentBytes)
	if err != nil {
		return "Failed to send document", err
	}

	// Serialize the birthdate
	birthdateBytes := []byte(bet.Birthdate)
	err = writeFully(p.conn, birthdateBytes)
	if err != nil {
		return "Failed to send birthdate", err
	}

	// Serialize the number
	numberBytes := htonl(int(bet.Number))
	err = writeFully(p.conn, numberBytes)
	if err != nil {
		return "Failed to send number", err
	}

	return "Bet sent", nil
}

func (p *Protocol) receiveMessage() (int, error) {
	messageBytes := make([]byte, 4)
	_, err := readFully(p.conn, messageBytes)
	if err != nil {
		return 0, err
	}
	messageType := int(ntohl(messageBytes))
	return messageType, nil
}