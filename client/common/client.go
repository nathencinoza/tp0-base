package common

import (
	"net"
	"sync" // ver si tengo que usarlo
	"time"
	"strconv"
	"os"
	"encoding/csv"

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
	stopCh chan struct{}  
	wg     sync.WaitGroup
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		stopCh: make(chan struct{}),
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	c.createClientSocket()
	protocol := NewProtocol(c.conn)
	log.Infof("protocolo creado", protocol)

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {	
		file_path := os.Getenv("FILE_PATH")
		log.Infof("action: open_file | result: success | client_id: %v | file_path: %v",
			c.config.ID,
			file_path,
		)
		file, err := os.Open(file_path)
		if err != nil {
			log.Errorf("action: open_file | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			log.Errorf("action: read_csv | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		batchSize := 0
		batchBets := make([]Bet, 0, c.config.MaxBatch)

		for _, record := range records {
			if len(record) != 5 {
				log.Errorf("action: parse_record | result: skip | client_id: %v | record: %v",
					c.config.ID,
					record,
				)
				return
			}

			document, err := strconv.Atoi(record[2])
			if err != nil {
				log.Errorf("action: parse_document | result: skip | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}

			number, err := strconv.Atoi(record[4])
			if err != nil {
				log.Errorf("action: parse_number | result: skip | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}

			bet := Bet{
				Name:      record[0],
				Surname:   record[1],
				Document:  document,
				Birthdate: record[3],
				Number:    number,
			}
			batchBets = append(batchBets, bet)
			batchSize++

			if batchSize == c.config.MaxBatch {
				batchSize = 0
				_, err = protocol.sendBets(batchBets)
				
				if err != nil {
					log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v",
						c.config.ID,
						err,
					)
					return
				}
				response, err := protocol.receiveMessage()
				if err != nil {
					log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
						c.config.ID,
						err,
					)
					return
				}
		
				if response == OK {
					log.Infof("action: apuestas_enviadas | result: success | dni: %v | cantidad: %v",
						document,
						c.config.MaxBatch,
					)
				} else {
					log.Errorf("action: apuestas_enviadas | result: fail | client_id: %v | cantidad: %v",
						c.config.ID,
						c.config.MaxBatch,
					)
				}
			}
		}
		time.Sleep(c.config.LoopPeriod)
		
	}
	c.conn.Close()
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)

}



// Stop Gracefully stops the client by closing the stop channel and waiting for
// the loop to finish its current iteration.
func (c *Client) Stop() {
	close(c.stopCh)
	c.wg.Wait()

	if c.conn != nil {
		c.conn.Close()
		log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
	}
	log.Infof("action: stop_client | result: success | client_id: %v", c.config.ID)
}