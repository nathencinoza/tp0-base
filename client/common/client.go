package common

import (
	"net"
	"time"
	"strconv"
	"os"
	"encoding/csv"
    "path/filepath"
    "strings"
	"fmt"
	"io"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	MaxBatch	  int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}
func (c *Client) StartClientLoop() {
    c.createClientSocket()
    defer c.conn.Close()
    protocol := NewProtocol(c.conn)

    filePath := os.Getenv("FILE_PATH")
    file, err := os.Open(filePath)
    if err != nil {
		log.Errorf("action: open_file | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
        return
    }
    defer file.Close()

	if err := c.sendAgency(protocol); err != nil {
		log.Errorf("action: send_agency | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

    reader := csv.NewReader(file)

	

    for {
        records, err := c.readBatchCSV(reader, c.config.MaxBatch) 
        if err == io.EOF && len(records) == 0 {
            break
        }
        if err != nil && err != io.EOF {
			log.Errorf("action: read_batch_csv | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
            return
        }

        batchBets := make([]Bet, 0, c.config.MaxBatch)
        for _, record := range records {
            bet, err := c.parseRecord(record)
            if err != nil {
				log.Errorf("action: parse_record | result: skip | client_id: %v | record: %v",
					c.config.ID,
					record,
				)
                continue
            }
            batchBets = append(batchBets, bet)
        }

        if len(batchBets) > 0 {
            if err := c.sendBatch(protocol, batchBets); err != nil {
                return
            }
        }
        if err == io.EOF {
            break
        }
    }

    protocol.sendFinish()
    c.handleWinners(protocol)
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) readBatchCSV(reader *csv.Reader, batchSize int) ([][]string, error) {
	var records [][]string
	for i := 0; i < batchSize; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			return records, io.EOF
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (c *Client) parseRecord(record []string) (Bet, error) {
    if len(record) != 5 {
        return Bet{}, fmt.Errorf("invalid record")
    }

    document, err := strconv.Atoi(record[2])
    if err != nil {
        return Bet{}, fmt.Errorf("invalid document")
    }

    number, err := strconv.Atoi(record[4])
    if err != nil {
        return Bet{}, fmt.Errorf("invalid number")
    }

    return Bet{
        Name:      record[0],
        Surname:   record[1],
        Document:  document,
        Birthdate: record[3],
        Number:    number,
    }, nil
}
func (c *Client) getAgencyFromFileName(filePath string) (int, error) {
    fileName := filepath.Base(filePath)
    agencyStr := strings.TrimPrefix(fileName, "agency-")
    agencyStr = strings.TrimSuffix(agencyStr, ".csv")
    return strconv.Atoi(agencyStr)
}
func (c* Client) sendAgency(protocol *Protocol) error {
	agency, err := c.getAgencyFromFileName(os.Getenv("FILE_PATH"))
	if err != nil {
		return err
	}
	err = protocol.sendAgency(agency)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) sendBatch(protocol *Protocol, bets []Bet) error {
    _, err := protocol.sendBets(bets)
    if err != nil {
		log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}

    response, err := protocol.receiveMessage()
    if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	if response == OK {
		log.Infof("action: apuestas_enviadas | result: success | client_id: %v| cantidad: %v",
		c.config.ID,
		len(bets),
	)
	} else {
		log.Errorf("action: apuestas_enviadas | result: fail | client_id: %v | cantidad: %v",
		c.config.ID,
		len(bets),
	)
	}
    return nil
}

func (c *Client) handleWinners(protocol *Protocol) {
	winners, err := protocol.receiveWinners()
    if err != nil {
        log.Errorf("action: consulta_ganadores | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
        return
    }
	log.Infof("action: consulta_ganadores | result: success | client_id: %v | winners: %v",
		c.config.ID,
		winners,
	)
}

// Stop Gracefully stops the client by closing the stop channel and waiting for
// the loop to finish its current iteration.
func (c *Client) Stop() {
	if c.conn != nil {
		c.conn.Close()
	}
	log.Infof("action: stop_client | result: success | client_id: %v", c.config.ID)
}