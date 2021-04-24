#!/usr/bin/env python3

import socket
import time

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(('127.0.0.1', 7712))

print('setting db to test')
s.sendall(b'db test\n')
print(s.recv(1024))

print('getting foo')
s.sendall(b'get foo\n')
print(s.recv(1024))

print('setting foo key to bar for 3 seconds')
s.sendall(b'set foo bar 3\n')
print(s.recv(1024))

print('getting foo')
s.sendall(b'get foo\n')
print(s.recv(1024))

print('getting foo')
s.sendall(b'get foo\n')
print(s.recv(1024))

print('getting foo after 3 seconds')
time.sleep(3)
s.sendall(b'get foo\n')
print(s.recv(1024))

print('bad duration')
print('setting foo key to bar for z seconds')
s.sendall(b'set foo bar z\n')
print(s.recv(1024))

s.close()
