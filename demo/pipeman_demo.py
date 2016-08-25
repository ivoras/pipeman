import sys,os,socket,select,time
from threading import Thread

HOST = 'localhost'
PORT = 4096

def worker(name):
    bname = name.encode('utf-8')
    cn = socket.create_connection((HOST, PORT))
    cn.sendall(b"%s\n" % bname) # Send the node name as the first thing
    x = 0
    tm = time.time()
    while True:
        r, w, e = select.select([cn.fileno()], [], [], 0.1)
        if len(r) > 0:
            while True:
                buf = cn.recv(4096)
                print("%s received: %s" % (name, repr(buf)))
                if len(buf) != 4096:
                    break
        if time.time() - tm >= 1: # every second
            cn.sendall(b"%s%d\n" % (bname, x))
            x += 1
            if x == 10: # send 10 messages
                break
            tm = time.time()
    print("exit", name)


def main():
    threads = (Thread(target=worker, args=('eenie', )), Thread(target=worker, args=('meanie',)), Thread(target=worker, args=('mynie',)), Thread(target=worker, args=('moe',)))
    for th in threads:
        th.start()
    for th in threads:
        th.join()


if __name__ == '__main__':
    main()