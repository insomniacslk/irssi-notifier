#!/usr/bin/env python3

import argparse
import subprocess

import requests


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('-s', '--server', required=True,
        help='The remote server to connect to')
    parser.add_argument('-p', '--port', type=int, default=8080,
        help='The server port to connect to (default: %(default)s)')
    parser.add_argument('-e', '--endpoint', default='/tail',
        help='The endpoint to connect to (default: %(default)s)')
    parser.add_argument('-c', '--certificate',
        help='The server certificate used to validate the TLS connection')
    return parser.parse_args()


def main():
    args = parse_args()
    url = 'https://{server}:{port}/{endpoint}'.format(
        server=args.server,
        port=args.port,
        endpoint=args.endpoint
    )

    r = requests.get(url, stream=True, verify=args.certificate)
    for chunk in r.iter_lines():
        print(chunk)
        subprocess.run(['notify-send', 'IRC', chunk.decode('utf-8')])


if __name__ == '__main__':
    main()
