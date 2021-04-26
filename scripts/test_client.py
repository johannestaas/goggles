#!/usr/bin/env python3

import socket
import time

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(('127.0.0.1', 7712))


def send(msg):
    print(f'sending {msg!r}')
    s.sendall(msg.encode() + b'\n')


def get(expected):
    result = s.recv(1024).decode()
    print(f'received: {result!r}\n')
    if result != expected:
        print(f'expected: {expected!r}\n')
        raise AssertionError


send('db test')
get('\n')

send('get foo')
get('\n')

send('set 1 foo bar')
get('\n')

send('get foo')
get('bar\n')

send('get foo')
get('bar\n')

print('getting foo after 1 seconds')
time.sleep(1)
send('get foo')
get('\n')

send('set 0 foo bar')
get('\n')
print('getting foo after 1 second')
time.sleep(1)
send('get foo')
get('bar\n')

print('dropping kvstore')
send('drop')
get('\n')
print('getting back test kvstore')
send('db test')
get('\n')
print('getting foo, which should be empty now')
send('get foo')
get('\n')

print('setting foo key to bar for XXX seconds')
send('set XXX foo bar')
get('error: bad command string\n')

# XXX add timeout to exit
print('setting foo baz to nothing')
send('set 0 baz')
get('error: bad command\n')

print('PASS!')

s.close()
