# -*- coding: utf-8 -*-
"""
Created on Tue Mar  8 22:13:35 2016

@author: eikes
"""
import select

import cfg
import echo_server


def main():
    #handlers is a list of open client-connections and a single
	# server acception new connections
    handlers = []
    handlers.append(echo_server.Server(cfg.address, handlers))
    
    #eventloop:
    while True:
        #checks which clients should get an answer
        senders = [h for h in handlers if h.wants_to_send()]
        recvs, sends, _ = select.select(handlers, senders, [], 2)
        for h in recvs:
            h.recv()
        for h in senders:
            h.send()


if __name__ == '__main__':
    main()