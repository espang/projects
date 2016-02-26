# -*- coding: utf-8 -*-
"""
Created on Fri Feb 26 09:11:06 2016

@author: eikes
"""
import numpy as np


class VariableResult(object):
    
    calc_types = {
        1: '_sum',
        2: '_starts',
    }
    default_calc = calc_types.get(1)

    def __init__(self, pk, label, comp_name, calc_type):
        self.pk = pk        
        self.label = label
        self.comp_name = comp_name
        self.calc_func = getattr(self, VariableResult.calc_types.get(
            calc_type, VariableResult.default_calc
        ))

    def _sum(self, array):
        return sum(array)

    def _starts(self, array):
        diff_arr = np.diff(np.array(array))
        starts = np.sum(diff_arr == 1)
        #shutdowns = np.sum(diff_arr == -1)
        return starts

    def calculate_value_from_array(self, array):
        return self.calc_func(array)


class MyOptResult(object):

    result_names = set()
    
    def __init__(self, nbr, value, time, worker):
        self.nbr = nbr        
        self.value = value
        self.time = time
        self.worker = worker
        self.result_dict = {}

    def add_result(self, name, value):
        self.result_dict[name] = value
        MyOptResult.result_names.add(name)

    @classmethod
    def header(cls):
        headers = ['Simulation', 'Worker', 'Value', 'Time in s']
        for label in sorted(cls.result_names):
            headers.append(label)
        return ';'.join(headers)

    def __repr__(self):
        representations = [ self.nbr, self.worker, self.value, '%.3f'%self.time ]
        for label in sorted(cls.result_names):
            representations.append(self.result_dict.get(label, 'nan'))
        return ';'.join(map(str, representations))

    @staticmethod
    def serialize_instance(obj):
        d = { '__classname__' : obj.__class__.__name__ }
        d.update(vars(obj))
        return d

    @staticmethod
    def unserialize_object(d):
        clsname = d.pop('__classname__', None)
        if clsname:
            cls = eval(clsname)
            obj = cls.__new__(cls)
            for key, value in d.items():
                setattr(obj, key, value)
            return obj
        else:
            return d
