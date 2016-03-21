# -*- coding: utf-8 -*-
"""
Created on Tue Mar  8 22:15:53 2016

Tests a running server-process

uses py.test for testing

@author: eikes
"""
from contextlib import contextmanager
import socket

import cfg

@contextmanager
def connection():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(cfg.address)
    try:
        yield s
    finally:
        s.close()

def test_simple_connection():
    with connection() as s:    
        s.send(b'hello\n')
        res = s.recv(8192)
        assert res == b'hello\n'
        
        s.send(b'line1 \n line 2 part 1')
        res = s.recv(8192)
        assert res == b'line1 \n'
        
        s.send(b'line 2 part 2 \n')
        res = s.recv(8192)
        assert res == b' line 2 part 1line 2 part 2 \n'

def test_multiple_clients_one():
    with connection() as s1:
        msg1 = b'hello. Wait for it'
        s1.send(msg1)
        
        with connection() as s2:
            s2.send(b'hello\n')
            res = s2.recv(8192)
            assert res == b'hello\n'
        
        s1.send(b'\n')
        res = s1.recv(8192)
        assert res == b'hello. Wait for it\n'

def test_multiple_clients_two():
    with connection() as s1:
        s1.send(b'hello. Wait for it\n')
        
        with connection() as s2:
            s2.send(b'hello\n')
            res = s2.recv(8192)
            assert res == b'hello\n'
        
        res = s1.recv(8192)
        assert res == b'hello. Wait for it\n'
    
def test_long_message():
    with connection() as s:
        # send a msg bigger than 8192 bytes
        msg = b'abcdefghi ' * 1000 + b'\n'
        n = s.send(msg)
        assert n == 10001
        res1 = s.recv(8192)
        res2 = s.recv(8192)
        assert res1 == msg[:8192]
        assert res2 == msg[8192:]