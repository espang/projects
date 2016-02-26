# -*- coding: utf-8 -*-
"""
Created on Fri Feb 26 09:11:06 2016

@author: eikes
"""
import ConfigParser

from components import Component
from result import VariableResult

_config = ConfigParser.ConfigParser()
_config.read('scenario.cfg')

_section = 'MySection'
_results = 'results'


def _create_comp(index):
    global _config, _section
    connections = map(str.strip, _config.get(
        _section,
        'comp.{0}.connections'.format(index),
    ).split(','))
    return Component(
        _config.get(_section, 'comp.{0}.name'.format(index)),
        _config.get(_section, 'comp.{0}.type'.format(index)),
        _config.get(_section, 'comp.{0}.reference_values'.format(index)),
        connections,
        _config.get(_section, 'comp.{0}.replace_values'.format(index)),
        _config.getfloat(_section, 'comp.{0}.factor'.format(index)),
        )


def _create_results():
    global _config, _results
    quantity = _config.getint(_results, 'quantity')
    results = []
    for i in range(1, quantity+1):
        label = _config.get(_results, 'result.{0}.name'.format(i))
        comp = _config.get(_results, 'result.{0}.comp'.format(i))
        calc_type = _config.getint(_results, 'result.{0}.type'.format(i))
        results.append(
            VariableResult(pk=i, label=label, comp_name=comp, calc_type=calc_type)        
        )
    return results

LP_FILE_PATH = _config.get(_section, 'lp')
TRC_FILE_PATH = _config.get(_section, 'trc')

QUANTITY = _config.getint(_section, 'quantity')

COMPONENTS = [ _create_comp(i) for i in range(1, QUANTITY+1) ]
RESULTS = _create_results()

SIMULATIONS = _config.getint(_section, 'simulations')
WORKER = _config.getint(_section, 'worker')

S_VALUE = float(1.5855e+07)




