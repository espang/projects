# -*- coding: utf-8 -*-
"""
Created on Tue Mar  8 22:11:11 2016

@author: eikes
"""

class Client:
    
    def __init__(self, sock, handlers):
        self.sock = sock
        self.handlers = handlers
        self.out = bytearray()
        
    def fileno(self):
        return self.sock.fileno()

    def wants_to_send(self):
        return b'\n' in self.out
        
    def send(self):
        idx = self.out.find(b'\n')        
        n = self.sock.send(self.out[:idx+1])
        self.out = self.out[n:]

    def recv(self):
        data = None
        try:
            data = self.sock.recv(8192)
        except ConnectionResetError:
            # close to remove client from handlers
            self.close()
            return          
        if not data:
            self.close()
        else:
            self.out.extend(data)

    def close(self):
        self.sock.close()
        self.handlers.remove(self)

    