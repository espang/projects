# -*- coding: utf-8 -*-
"""
Created on Fri Feb 26 09:11:06 2016

@author: eikes
"""
import json
import logging
import multiprocessing
import re
import time

from collections import namedtuple

import numpy as np
import zmq
from zmq.log.handlers import PUBHandler
import gurobipy as gp

from result import MyOptResult
from log_wrappers import logged, formatters


VariableConstraints = namedtuple(
    'VariableConstraint',
    ['constr', 'var', 'timestep', 'idn']
)


class Worker(multiprocessing.Process):

    def __init__(self, worker_nbr, add_pull, add_push, add_sub, add_log, 
                 components, lp_file, trc_file, results):
        super(Worker, self).__init__()
        self.worker_nbr = worker_nbr
        self.add_pull = add_pull
        self.add_push = add_push
        self.add_sub = add_sub
        self.add_log = add_log
        self.components = components
        self.result_dict = { r.pk: r for r in results }
        self.lp_file = lp_file
        self.trc_file = trc_file

    @logged(performance=True)
    def _read_model(self, worker=-1):
        self.log.debug('Worker {nbr:02d}: read model from {fname}'.format(
            fname=self.lp_file,
            nbr=worker,
            )
        )
        self.model = gp.read(self.lp_file)
        #self.model.setParam(gp.GRB.Param.OutputFlag, 0)
        #self.model.setParam('Presolve', 0)
        self.model.setParam(gp.GRB.Param.Threads, 1)

    @logged(performance=True)
    def _read_tracefile(self, worker=-1):
        with open(self.trc_file, 'r') as f:
            content = f.read()
        match = re.search('Anzahl Zeitschritte: (\d*)', content)
        if not match:
            raise Exception('wrong file format in tracefile')
        self.timesteps = int(match.group(1))
        
        pattern = '{name}\s*IDN\s\s\s:\s(\d*)'
        for c in self.components:
            match = re.search(pattern.format(name=c.name), content)
            if not match:
                raise Exception('Component {name} not found'.format(name=c.name))
            idn = int(match.group(1))
            c.set_idn(idn)

    def _handle_var(self, name, component, timestep):
        coeff = component.array[timestep]
        var = self.model.getVarByName(name)
        if var is None:
            self.log.warn('Variable {var} is not in the model'.format(
                var=name,
                )
            )
            return
        col = self.model.getCol(var)
        for i in range(col.size()):
            constr = col.getConstr(i)
            row = self.model.getRow(constr)
            for k in range(row.size()):
                if row.getVar(k).sameAs(var):
                    if row.getCoeff(k) == -coeff:
                        return(
                            VariableConstraints(
                                        constr=constr,
                                        var=var,
                                        timestep=timestep,
                                        idn=component.idn,
                            )
                        )
    
    @logged(performance=True)
    def _get_variable_constraints(self, worker=-1):
        self.variable_constraints = []
        for c in self.components:
            for var_name in c.var_names:
                for timestep in range(self.timesteps):
                    n = '{component}.{h}'.format(component=var_name, h=hex(timestep)[2:])
                    self.variable_constraints.append(
                        self._handle_var(n, c, timestep),
                    )

    @logged(performance=True)
    def _get_variables_for_results(self, worker=-1):
        self.results_vars = {}
        for pk, r in self.result_dict.items():            
            array = []
            comp_name = r.comp_name
            #self.log.info('Result {0} has name {1}'.format(pk, comp_name))
            for timestep in range(self.timesteps):
                n = '{comp_name}.{h}'.format(comp_name=comp_name, h=hex(timestep)[2:])
                array.append(self.model.getVarByName(n))
            self.results_vars[pk] = array

    def _get_data(self, simu_nbr):        
        return {            
            c.idn: c.get_data_for_simulation(simu_nbr) for c in self.components
        }

    @logged(performance=True)
    def _change_constraints(self, data, worker=-1):
        for c in self.variable_constraints:
            timestep = c.timestep
            idn = c.idn
            self.model.chgCoeff(c.constr, c.var, -data.get(idn)[timestep])

    @logged(performance=True)
    def _optimize(self, simu_nbr, worker=-1):
        t_start = time.time()    
        try:
            self.model.optimize()
            value = -self.model.objval
        except gp.GurobiError:
            self.log.warn('Error during optimization', exc_info=1)
            value = 0.0
        t_end = time.time()   
        # Do the work
        time.sleep(1)
        self.simu_result = MyOptResult(
            nbr=simu_nbr,
            value=value,
            time=t_end-t_start,
            worker=worker,
        )

    @logged(performance=True)
    def _calculate_results(self, worker=-1):
        for pk, _vars in self.results_vars.iteritems():
            tmp_list = []
            for v in _vars:
                tmp_list.append(v.x)
            result = self.result_dict.get(pk)
            value = result.calculate_value_from_array(tmp_list)
            name = result.label
            self.simu_result.add_result(name=name, value=value)

    def run(self):
        context = zmq.Context()       

        pub = context.socket(zmq.PUB)
        pub.connect(self.add_log)
        self.log = logging.getLogger()
        self.log.setLevel(logging.DEBUG)
        handler = PUBHandler(pub)
        handler.formatters = formatters
        self.log.addHandler(handler)
            
        # Socket to receive messages on
        receiver = context.socket(zmq.PULL)
        receiver.connect(self.add_pull)
        
        # Socket to send messages to
        sender = context.socket(zmq.PUSH)
        sender.connect(self.add_push)
        
        # Socket for control input
        controller = context.socket(zmq.SUB)
        controller.connect(self.add_sub)
        controller.setsockopt(zmq.SUBSCRIBE, b"")
                
        poller = zmq.Poller()
        poller.register(receiver, zmq.POLLIN)
        poller.register(controller, zmq.POLLIN)
        
        self._read_model(worker=self.worker_nbr)
        self._read_tracefile(worker=self.worker_nbr)
        self._get_variable_constraints(worker=self.worker_nbr)
        self._get_variables_for_results(worker=self.worker_nbr)
        
        while True:
            socks = dict(poller.poll())
        
            if socks.get(receiver) == zmq.POLLIN:
                message = receiver.recv_string()
                
                simu_nbr = int(message)
                
                data = self._get_data(simu_nbr)
                self._change_constraints(data, worker=self.worker_nbr)
                self._optimize(simu_nbr, worker=self.worker_nbr)
                self._calculate_results(worker=self.worker_nbr)
                
                # Send results to sink
                sender.send(
                    json.dumps(
                        self.simu_result,
                        default=MyOptResult.serialize_instance,
                    )
                )
        
            # Any waiting controller command acts as 'KILL'
            if socks.get(controller) == zmq.POLLIN:
                self.log.debug('Worker {nbr:02d}: Finished stop worker'.format(
                    nbr=self.worker_nbr,            
                ))                
                break

if __name__ == '__main__':
    print 'not meant to be executed!'