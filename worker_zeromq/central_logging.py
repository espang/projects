# -*- coding: utf-8 -*-
"""
Created on Fri Feb 26 09:11:06 2016

@author: eikes
"""
import logging
import multiprocessing

import zmq

LOG_LEVELS = {
    'DEBUG': logging.DEBUG,
    'INFO': logging.INFO,
    'WARN': logging.WARN,
    'ERROR': logging.ERROR,
    'CRITICAL': logging.CRITICAL,
}

LOG_FILE_NAME = 'optimization.log' 
LOG_PATTERN = '%(asctime)s - %(message)s'


class Logger(multiprocessing.Process):

    def __init__(self, add_sub, add_control):
        super(Logger, self).__init__()
        self.add_sub = add_sub
        self.add_control = add_control

    def run(self):
        context = zmq.Context()

        socket_fd = context.socket(zmq.SUB)
        socket_fd.bind(self.add_sub)
        socket_fd.setsockopt(zmq.SUBSCRIBE, "")

        # Socket for control input
        controller = context.socket(zmq.SUB)
        controller.connect(self.add_control)
        controller.setsockopt(zmq.SUBSCRIBE, b"")
        
        poller = zmq.Poller()
        poller.register(socket_fd, zmq.POLLIN)
        poller.register(controller, zmq.POLLIN)
    
        filehandler = logging.FileHandler(LOG_FILE_NAME, mode='w')
        log = logging.getLogger()
        log.setLevel(logging.DEBUG)
        filehandler.setLevel(logging.DEBUG)
        formatter = logging.Formatter(LOG_PATTERN)
        filehandler.setFormatter(formatter)
        log.addHandler(filehandler)
        
        while True:
            socks = dict(poller.poll())
            if socks.get(socket_fd) == zmq.POLLIN:
                topic, message = socket_fd.recv_multipart()
                pos = topic.find('.')
                level = topic
                if pos > 0: 
                    level = topic[:pos]
                if message.endswith('\n'):
                    message = message[:-1]
                log_msg = getattr(logging, level.lower())
                if pos > 0:
                    message = topic[pos+1:] + " | " + message
                log_msg(message)
            if socks.get(controller) == zmq.POLLIN:
                log.debug('Kill command: stop logger')                
                break

if __name__ == '__main__':
    print 'not meant to be executed!'