# Pipeman: a simulator of networked broadcast domains, v0.5

Pipeman is a nondeterministic simulator of broadcast domains organized into networks, using simple TCP. Multiple broadcast domains can be specified in a JSON file, and the server listens for connections from ordinary TCP clients. Each `send()` event from a client gets broadcast to all the other nodes in the client's broadcast domains, optionally  with a certain chance for data loss. It's useful for experimenting with wireless mesh protocols.

## Example

This JSON config file describes a network with three broadcast domains and four nodes:

    {
        "type": "pipeman",
        "port": 4096,
        "buffer_size": 64,
        "network": [
            {
                "name": "red",
                "nodes": ["eenie", "meanie"],
                "loss": 0
            },
            {
                "name": "green",
                "nodes": ["meanie", "mynie"],
                "loss": 0.5
            },
            {
                "name": "blue",
                "nodes": ["meanie", "moe"],
                "loss": 0.1
            }
        ]
    }

Here's what's happening:

* The nodes are named "eenie", "meanie", "mynie" and "moe"
* "eenie" and "meanie" are in a broadcast domain named "red", "meanie" and "mynie" are in their own broadcast domain named "green", and "meanie" and "moe" are in their own, named "blue". For example, whatever the "eenie" node sends, is received by "meanie". Note that nodes can belong to multiple domains. If the nodes (re)broadcast messages, loops are naturally possible ("broadcast storms") and there are no guards against this case.
* Domains "green" and "blue" have a certain random chance of data loss, 0.5 and 0.1 (in the range of 0 = never and 1 = always) 

## Buffer size and data loss

Data is always read and written by the server in chunks of `buffer_size`, and random data loss chance is always calculated for such a buffer before data is received by each individual node. Consider the following extreme cases:

* `buffer_size = 1` means there is a chance of data loss for every single byte, in each domain where `loss` is non-zero. It's also very inefficient.
* Setting `loss` to 0 in each domain, and having a large `buffer_size` (for example 16384) results in very efficient operation, but without data loss simulation.

## The protocol 

The protocol used is (currently) very simple: As soon as a node connects to the server it must send its name (case sensitive) followed by a newline. Immediately after that, it can send (and) receive whatever data it needs to.

## Command line usage

    usage: pipeman.exe [-c config.json]
    -v
            Verbose output
    -c string
            Short for config file name (default "config.json")
    -config string
            Config file name (default "config.json")

The usual command line for starting Pipeman is something like `pipeman -c config.json`.

See the `pipeman_demo.py` example in the `demo` directory for an example which uses threads to simulate multiple nodes connecting to a single Pipeman server.

## License

Copyright (c) 2016, Ivan Voras <ivoras@gmail.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.