package common


import (
	"encoding/binary"
	"net"
)

type Bet struct {
	Agency	   string
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

// SerializeBet serializes a Bet structure into a byte slice and sends it
func (p *Protocol) sendBets(bets []Bet) (string, error) {
	betsAmount := htonl(len(bets))
	_, err := p.conn.Write(betsAmount)
	if err != nil {
		return "Failed to send bets amount", err
	}

	for _, bet := range bets {
		// Serialize the bet type
		betBytes := htonl(int(BET))
		_, err := p.conn.Write(betBytes)
		if err != nil {
			return "Failed to send bet type", err
		}

		// Serialize the agency	
		agencyBytes := []byte(bet.Agency)
		agencySizeBytes := htonl(len(agencyBytes))
		_, err = p.conn.Write(agencySizeBytes)
		if err != nil {
			return "Failed to send agency size", err
		}
		_, err = p.conn.Write(agencyBytes)
		if err != nil {
			return "Failed to send agency", err
		}

		// Serialize the name
		nameBytes := []byte(bet.Name)
		nameSizeBytes := htonl(len(nameBytes))
		_, err = p.conn.Write(nameSizeBytes)
		if err != nil {
			return "Failed to send name size", err
		}
		_, err = p.conn.Write(nameBytes)
		if err != nil {
			return "Failed to send name", err
		}

		// Serialize the surname
		surnameBytes := []byte(bet.Surname)
		surnameSizeBytes := htonl(len(surnameBytes))
		_, err = p.conn.Write(surnameSizeBytes)
		if err != nil {
			return "Failed to send surname size", err
		}
		_, err = p.conn.Write(surnameBytes)
		if err != nil {
			return "Failed to send surname", err
		}

		// Serialize the document
		documentBytes := htonl(int(bet.Document))
		_, err = p.conn.Write(documentBytes)
		if err != nil {
			return "Failed to send document", err
		}

		// Serialize the birthdate
		birthdateBytes := []byte(bet.Birthdate)
		_, err = p.conn.Write(birthdateBytes)
		if err != nil {
			return "Failed to send birthdate", err
		}

		// Serialize the number
		numberBytes := htonl(int(bet.Number))
		_, err = p.conn.Write(numberBytes)
		if err != nil {
			return "Failed to send number", err
		}
	}
	return "success", nil
}

func (p *Protocol) receiveMessage() (int, error) {
	messageBytes := make([]byte, 4)
	_, err := p.conn.Read(messageBytes)
	if err != nil {
		return ERROR, err
	}
	messageType := int(ntohl(messageBytes))
	return messageType, nil
}
