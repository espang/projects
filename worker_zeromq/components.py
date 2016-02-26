# -*- coding: utf-8 -*-
"""
Created on Fri Feb 26 09:11:06 2016

@author: eikes
"""
import numpy as np


def read_data(fname):
    array = []
    with open(fname, 'r') as data_file:
        for line in data_file:
            array.append(float(line))
    return np.array(array)


class Component():
    name_pattern = 'K{idn}{shortname}'
    
    def __init__(self, name, cmp_type, fname, connections, newdata_filename, factor):
        self.name = name
        self.type = cmp_type
        self.connections = connections
        self.array = read_data(fname)
        self.filename = newdata_filename
        self.factor = factor
        
    def set_idn(self, idn):
        self.idn = idn
        self.var_names = [
            Component.name_pattern.format(idn=idn, shortname=shortname)
            for shortname in self.connections
        ]

    def get_data_for_simulation(self, simu_nbr):
        return self.factor * read_data(self.filename.format(simu_nbr=simu_nbr))
