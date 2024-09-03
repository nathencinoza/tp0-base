package common

import (
	"net"
	"sync" // ver si tengo que usarlo
	"time"
	"strconv"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
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
	
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {	
		document, err := strconv.Atoi(os.Getenv("DOCUMENTO"))
		if err != nil {
			log.Errorf("action: convert_document | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		number, err := strconv.Atoi(os.Getenv("NUMERO"))
		if err != nil {
			log.Errorf("action: convert_number | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		name := os.Getenv("NOMBRE")
		surname := os.Getenv("APELLIDO")
		birthdate := os.Getenv("NACIMIENTO")

		bet := Bet{
			Name:      name,
			Surname:   surname,
			Document:  document,
			Birthdate: birthdate,
			Number:    number,
		}

		c.createClientSocket()

		protocol := NewProtocol(c.conn)

		_, err = protocol.sendBet(bet)

		if err != nil {
			log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
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
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
				document,
				number,
			)
		} else {
			log.Errorf("action: receive_message | result: fail | client_id: %v | response: %v",
				c.config.ID,
				response,
			)
		}
		c.conn.Close()
		time.Sleep(c.config.LoopPeriod)

	}
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