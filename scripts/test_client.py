#!/usr/bin/env python3

import socket

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(('127.0.0.1', 7712))

print('setting db to test')
s.sendall(b'db test\n')
print(s.recv(1024))
print('getting foo')
s.sendall(b'get foo\n')
print(s.recv(1024))
print('setting foo key to bar')
s.sendall(b'set foo bar\n')
print(s.recv(1024))
print('getting foo')
s.sendall(b'get foo\n')
print(s.recv(1024))
s.close()
