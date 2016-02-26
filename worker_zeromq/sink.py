# -*- coding: utf-8 -*-
"""
Created on Fri Feb 26 09:11:06 2016

@author: eikes
"""
import json
import logging
import multiprocessing
import time

import zmq
from zmq.log.handlers import PUBHandler

import numpy as np
import seaborn as sns
import matplotlib.pyplot as plt

from log_wrappers import formatters
from result import MyOptResult
import resource as rs


class Sink(multiprocessing.Process):
    
    def __init__(self, add_push, add_rec, add_pub, add_log, tasks):
        super(Sink, self).__init__()
        self.tasks = tasks
        self.add_log = add_log
        self.add_push = add_push
        self.add_rec = add_rec
        self.add_pub = add_pub

    def run(self):
        context = zmq.Context()
        pub = context.socket(zmq.PUB)
        pub.connect(self.add_log)
        self.log = logging.getLogger()
        self.log.setLevel(logging.DEBUG)
        handler = PUBHandler(pub)
        handler.formatters = formatters
        self.log.addHandler(handler)
        self.log.debug('start_sink')
        
        self.receiver = context.socket(zmq.PULL)
        self.receiver.bind(self.add_rec)
    
        self.controller = context.socket(zmq.PUB)
        self.controller.bind(self.add_pub)
    
        #socket to sinks
        self.main = context.socket(zmq.PUSH)
        self.main.connect(self.add_push)
        
        # Message from main: start
        self.receiver.recv()
    
        #Measure time!
        t_start = time.time()
        results = []
        for task_nbr in range(self.tasks):
            raw_json_data = self.receiver.recv()
            res = json.loads(
                raw_json_data,
                object_hook=MyOptResult.unserialize_object
            )
            results.append(res)
        t_end = time.time()
        t_duration = t_end - t_start
        self.log.debug(
            'Collected {count} results'.format(
            count=len(results),        
            )
        )
        self.log.debug(
            'Total elapsed time: {duration} s'.format(duration=t_duration)
        )
    
        self.controller.send(b'KILL')
        self.handle_results(results)
        time.sleep(1)
        self.main.send(b'0')
        
    def handle_results(self, results):
        with open('results.csv', 'w') as f:
            f.write('%s\n' % MyOptResult.header())
            for result in sorted(results, key=lambda x: x.nbr):
                f.write('%s\n' % result)        
        sns.set_context('talk', font_scale=1.8)
        plt.figure(figsize=(12,8))
        dist = [ result.value for result in results ]
        mean = np.mean(dist)
        try:
            sns.distplot(dist, norm_hist=False)#, fit=stats.gamma)#, rug=True, hist=False)
            ymin, ymax = plt.ylim()
            plt.vlines(rs.S_VALUE, ymin, ymax, color='#ffd700', linestyle='--', label='IW')
            plt.vlines(mean, ymin, ymax, color='#32cd32', linestyle='--', label='Mean')
            plt.legend()
            plt.savefig('results.png')
        except:
            self.log.warn('errors plotting the results...', exc_info=1)

if __name__ == '__main__':
    print 'not meant to be executed!'