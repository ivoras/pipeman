# Pipeman broadcast domains network simulator

Pipeman is a simulator of broadcast domains in networks, using TCP. Multiple broadcast domains can be specified in a JSON file, and the server listens for connections from ordinary TCP clients. Each `send()` event from a client gets broadcast to all the other nodes in the client's broadcast groups, possibly with a certain chance for data loss. It's useful for experimenting with wireless mesh protocols.

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

* The nodes are "eenie", "meanie", "mynie" and "moe"
* "eenie" and "meanie" are in a broadcast domain named "red", "meanie" and "mynie" are in their own broadcast domain named "green", and "meanie" and "moe" are in their own, named "blue". For example, whatever the "eenie" node sends, get received by "meanie". Note that nodes can belong to multiple domains. If the nodes (re)broadcast messages, loops are naturally possible.
* Domains "green" and "blue" have a certain chance of data loss, 0.5 and 0.1 (in the range of 0 = never and 1 = always) 

## Buffer size and data loss

Data is always read and written by the server in chunks of `buffer_size`, and random data loss chance is always calculated for such a buffer before data is received by each individual node. Consider the following extreme cases:

* `buffer_size = 1` means there is a chance of data loss for every single byte, in each domain where `loss` is non-zero. It's also very inefficient.
* Setting `loss` to 0 in each domain, and having a large `buffer_size` (for example 16384) results in very efficient operation, but without data loss simulation.

## The protocol 

The protocol used is (currently) very simple: As soon as a node connects to the server it must send its name followed by a newline. Immediately after that, it can send (and) receive whatever data it needs to.