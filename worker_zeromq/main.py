# -*- coding: utf-8 -*-
"""
Created on Fri Feb 26 09:11:06 2016

@author: eikes
"""
import logging
import logging.handlers
import sys
import time

import zmq
from zmq.log.handlers import PUBHandler

from addresses import ADD_VENTI_WORKER, ADD_SINK_RECEIVE, \
    ADD_SINK_PUBLISH, ADD_SINK_LH_RECEIVE, ADD_VENTI_LH_WORKER, \
    ADD_SINK_LH_PUBLISH, ADD_VENTI_RECEIVE, ADD_VENTI_LH_RECEIVE, \
    ADD_LH_LOGGING, ADD_LOGGING, ADD_LOG_CONTROLLER, ADD_LOG_LH_CONTROLLER
from sink import Sink
from optimization_worker import Worker
from central_logging import Logger
from log_wrappers import formatters
import resource as rs


def _start_logger(context):
    global log
    
    controller = context.socket(zmq.PUB)
    controller.bind(ADD_LOG_CONTROLLER)

    l = Logger(ADD_LOGGING, ADD_LOG_LH_CONTROLLER)
    l.start()
    
    time.sleep(2)
    
    pub = context.socket(zmq.PUB)
    pub.connect(ADD_LH_LOGGING)
    log = logging.getLogger('main')
    log.setLevel(logging.DEBUG)
    handler = PUBHandler(pub)
    handler.formatters = formatters
    log.addHandler(handler)
    
    return controller

def _start_sink(number_of_simus, context):
    s = Sink(
        ADD_VENTI_LH_RECEIVE,
        ADD_SINK_RECEIVE,
        ADD_SINK_PUBLISH,
        ADD_LH_LOGGING,
        number_of_simus,
    )
    s.start()
    
    #socket to sinks
    sink = context.socket(zmq.PUSH)
    sink.connect(ADD_SINK_LH_RECEIVE)

    #socket from sink for ready    
    receiver = context.socket(zmq.PULL)
    receiver.bind(ADD_VENTI_RECEIVE)
    
    return sink, receiver

def _start_workers(nbr_of_workers, fname, trace, context, comps):
    global log    
    #socket to workers
    sender = context.socket(zmq.PUSH)
    sender.bind(ADD_VENTI_WORKER)
    for nbr in range(nbr_of_workers):
        log.debug('start worker {nbr}'.format(nbr=nbr))
        w = Worker(
                   worker_nbr=nbr,
                   add_pull=ADD_VENTI_LH_WORKER,
                   add_push=ADD_SINK_LH_RECEIVE,
                   add_sub=ADD_SINK_LH_PUBLISH,
                   add_log=ADD_LH_LOGGING,
                   components=comps,
                   lp_file=fname,
                   trc_file=trace,
                   results=rs.RESULTS,)
        w.start()
    return sender

def _start_work(sink):
    global log    
    time.sleep(2)
    log.debug('Start sending tasks to worker...')
    
    # Message to sink: start
    sink.send(b'0')


def main():
    global log
    
    context = zmq.Context()
    controller = _start_logger(context)

    sink, sink_receiver = _start_sink(
        number_of_simus=rs.SIMULATIONS,
        context=context,
    )

    sender = _start_workers(
        context=context,        
        nbr_of_workers=rs.WORKER,
        fname=rs.LP_FILE_PATH,
        trace=rs.TRC_FILE_PATH,
        comps=rs.COMPONENTS,
    )    
    
    _start_work(sink)

    for simu_number in range(rs.SIMULATIONS):
        sender.send(b'{simu_number}'.format(simu_number=simu_number))

    sink_receiver.recv()
    log.debug('main done! Kill Logger')
    controller.send(b'KILL')
    return 0
    

if __name__ == '__main__':
    global log
    log = None
    try:
        res = main()
    except:
        res = 1
        if log:
            log.error('Error in main.py: ', exc_info=1)
        else:
            raise
    sys.exit(res)
