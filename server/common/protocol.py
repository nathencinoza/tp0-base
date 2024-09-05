from common.utils import Bet
import logging  
BET_MESSAGE = 1
SUCCESS_CODE = 2
ERROR_CODE = 3
FINISH_CODE = 4

SIZE = 4
DATE_SIZE = 10


class Protocol:

    def __init__(self, socket):
        self.socket = socket
    
    def is_closed(self):
        try:
            self.socket.getpeername()
            return False
        except:
            return True
        
        
    def receive_exact(self, size):
        data = b''
        while len(data) < size:
            packet = self.socket.recv(size - len(data))
            if not packet:
                raise Exception("Connection closed unexpectedly")
            data += packet
        return data
    
    def receive_code(self):
        return int.from_bytes(self.receive_exact(SIZE), byteorder='big')

    def receive_code(self):
        code = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
        if code == FINISH_CODE:
            return "FINISH"
        else:
            return "BET"

    def receive_bets(self):
        bets = []
        bets_size = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
        for _ in range(bets_size):

            agency = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
            
            name_size = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
            name = self.receive_exact(name_size).decode()
            
            surname_size = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
            surname = self.receive_exact(surname_size).decode()
            
            document = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
            birthdate = self.receive_exact(DATE_SIZE).decode()
            number = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
            
            bet = Bet(agency, name, surname, document, birthdate, number)
            bets.append(bet)
        return bets_size, bets
    

    def send_success(self):
        self.socket.sendall(SUCCESS_CODE.to_bytes(SIZE, byteorder='big'))

    def send_error(self):
        self.socket.sendall(ERROR_CODE.to_bytes(SIZE, byteorder='big'))

    def send_winners(self, winners): 
        self.socket.sendall(len(winners).to_bytes(SIZE, byteorder='big'))
        for winner in winners:
            document_size = len(winner.document)
            self.socket.sendall(document_size.to_bytes(SIZE, byteorder='big'))
            self.socket.sendall(winner.document.encode())
