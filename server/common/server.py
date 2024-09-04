import socket
import logging
import signal
from common.protocol import Protocol
from common.utils import Bet, store_bets, load_bets, has_won
import threading
from threading import Lock

BET_MESSAGE = 1
FINISH_MESSAGE = 4
class Server:
    def __init__(self, port, listen_backlog, ammount_clients):
        """
        # Initialize server socket
        """
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._clients = []
        self.clients_threads = []
        self._is_active = True  
        self.ammount_clients_done = 0
        self.total_clients = ammount_clients
        self.lock = Lock()  


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        signal.signal(signal.SIGTERM, self.__handle_sigterm)

        while self._is_active:
            try:
                client_sock = self.__accept_new_connection()
                client_thread = threading.Thread(target=self.__handle_client_connection, args=(client_sock,))
                client_thread.start()
            except OSError as e:
                if not self._is_active:
                    break
                logging.error(f"action: server | result: fail | error: {e}")
                break

            
    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            protocol = Protocol(client_sock)
            code = protocol.receive_code()
            if code == BET_MESSAGE:
                size, bets  = protocol.receive_bets()
                self.__handle_bets(bets, size, protocol)
            elif code == FINISH_MESSAGE:
                self.ammount_clients_done += 1
                if self.ammount_clients_done == self.total_clients:
                    self.__handle_draw() 
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}", e) 
        finally:
            if (code != FINISH_MESSAGE):
                client_sock.close()
                for client in self._clients:
                    if client == client_sock:
                        self._clients.remove(client)


    def __handle_draw(self):
        try: 
            logging.info("action: sorteo | result: success")
            with self.lock:
                bets = load_bets()
            winners = [bet for bet in bets if has_won(bet)]
            for client in self._clients:
                protocol = Protocol(client)
                protocol.send_winners(winners)
        finally: 
            for client in self._clients:
                client.close()
            self._clients.clear()



    def __handle_bets(self, bets, size, protocol):
        with self.lock:
            store_bets(bets)
        if len(bets) != size:
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {size}")
            protocol.send_error()
            return
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {size}")
        protocol.send_success()
    
    def __handle_sigterm(self, signum, frame):
        logging.info(f"action: shutdown | result: received signal {signum}")
        self._is_active = False
        for client in self._clients:
            client.close()
        self._clients.clear()
        if self._server_socket:
            self._server_socket.close()
        logging.info("action: server_close | result: success")


    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        self._clients.append(c)

        return c
