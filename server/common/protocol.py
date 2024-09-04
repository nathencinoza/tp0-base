from common.utils import Bet

BET_MESSAGE = 1
SUCCESS_CODE = 2
ERROR_CODE = 3

SIZE = 4
DATE_SIZE = 10


class Protocol:

    def __init__(self, socket):
        self.socket = socket

    def receive_exact(self, size):
        data = b''
        while len(data) < size:
            packet = self.socket.recv(size - len(data))
            if not packet:
                raise Exception("Connection closed unexpectedly")
            data += packet
        return data

    def receive_bets(self):
        bets = []
        bets_size = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
        for _ in range(bets_size):
            code = int.from_bytes(self.receive_exact(SIZE), byteorder='big')
            if code != BET_MESSAGE:
                raise Exception("Invalid code")

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


        