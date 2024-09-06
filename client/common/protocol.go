package common


import (
	"encoding/binary"
	"net"
	"bytes"
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
	FINISH int = 4
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

func (p *Protocol)  sendAgency(agency int) (error) {
	agencyBytes := htonl(agency)
	err := writeFully(p.conn, agencyBytes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Protocol) sendBets(bets []Bet) (string, error) {
	var buffer bytes.Buffer

	betBytes := htonl(int(BET))
	err := writeFully(p.conn, betBytes)
	if err != nil {
		return "Failed to send bet message", err
	}

	betsAmount := htonl(len(bets))
	err = writeFully(p.conn, betsAmount)
	if err != nil {
		return "Failed to send bet amount", err
	}
	for _, bet := range bets {

		nameBytes := []byte(bet.Name)
		binary.Write(&buffer, binary.BigEndian, int32(len(nameBytes))) 
		buffer.Write(nameBytes)                                       

		surnameBytes := []byte(bet.Surname)
		binary.Write(&buffer, binary.BigEndian, int32(len(surnameBytes))) 
		buffer.Write(surnameBytes)                                        

		binary.Write(&buffer, binary.BigEndian, int32(bet.Document))

		birthdateBytes := []byte(bet.Birthdate)
		buffer.Write(birthdateBytes)

		binary.Write(&buffer, binary.BigEndian, int32(bet.Number))
	}

	err = writeFully(p.conn, buffer.Bytes())
	if err != nil {
		return "Failed to send batch", err
	}

	return "Batch sent", nil
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



func (p *Protocol) receiveWinners() (int, error) {
    amountOfWinnersBytes := make([]byte, 4)
    _, err := readFully(p.conn, amountOfWinnersBytes)
	if err != nil {
		return 0, err
	}
    amountOfWinners := int(ntohl(amountOfWinnersBytes))

	for i := 0; i < amountOfWinners; i++ {
		numberBytes := make([]byte, 4)
		_, err := p.conn.Read(numberBytes)
		if err != nil {
			return 0, err
		}

		prizeBytes := make([]byte, 4)
		_, err = p.conn.Read(prizeBytes)
		if err != nil {
			return 0, err
		}
	}

	return amountOfWinners, nil
}
func (p *Protocol) sendFinish() error {
	finishBytes := htonl(int(FINISH))
	err := writeFully(p.conn, finishBytes)
	if err != nil {
		return err
	}
	return nil
}