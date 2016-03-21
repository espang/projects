# -*- coding: utf-8 -*-
"""
Created on Tue Mar  8 19:25:44 2016

@author: eikes
"""
import socket

import echo_client


class Server:
    
    def __init__(self, address, handlers, backlog=3):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, True)
        self.sock.bind(address)
        self.sock.listen(backlog)
        self.handlers = handlers
        
    def fileno(self):
        return self.sock.fileno()

    def recv(self):
        client, addr = self.sock.accept()
        self.handlers.append(echo_client.Client(client, self.handlers))

    def wants_to_send(self):
        return False
