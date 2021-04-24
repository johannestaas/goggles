#!/usr/bin/env python3

import socket
import time

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(('127.0.0.1', 7712))


def send(msg):
    print(f'sending {msg!r}')
    s.sendall(msg.encode() + b'\n')


def get():
    print(f'received: {s.recv(1024)!r}\n')


send('db test')
get()

send('get foo')
get()

send('set 1 foo bar')
get()

send('get foo')
get()

send('get foo')
get()

print('getting foo after 1 seconds')
time.sleep(1)
send('get foo')
get()

send('set 0 foo bar')
get()
print('getting foo after 1 second')
time.sleep(1)
send('get foo')
get()

print('setting foo key to bar for XXX seconds')
send('set XXX foo bar')
get()

print('setting foo baz to nothing')
send('set 0 baz')
get()

s.close()
