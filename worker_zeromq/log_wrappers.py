# -*- coding: utf-8 -*-
"""
Created on Fri Feb 26 09:11:06 2016

@author: eikes
"""
import logging
import time
from functools import wraps


formatters = {
    logging.DEBUG: logging.Formatter("%(filename)s:%(lineno)d | %(message)s"),
    logging.INFO: logging.Formatter("%(message)s"),
    logging.WARN: logging.Formatter("%(filename)s:%(lineno)d | %(message)s"),
    logging.ERROR: logging.Formatter("%(filename)s:%(lineno)d | %(message)s"),
    logging.CRITICAL: logging.Formatter("%(filename)s:%(lineno)d | %(message)s"),
}

def logged(performance=False):
    def decorate(func):
        log = logging.getLogger()
        
        @wraps(func)
        def wrapper(*args, **kwargs):
            worker_number = kwargs.get('worker', -2)
            modul = func.__module__
            name = func.__name__
            start = time.time()
            try:
                result = func(*args, **kwargs)
            except:
                log.error(
                    '[{modul}] - Error in func {name}'.format(
                        name=name,
                        modul=modul,
                    ),
                    exc_info=1
                )
                raise
            if performance is True:
                msg = '[{modul}] {nbr:02d} - Function {name} finished after {time:.3f} s'.format(
                            nbr=worker_number,                            
                            modul=modul,
                            name=name,
                            time=time.time()-start,
                )
                log.info(msg)
            return result
        return wrapper
    return decorate