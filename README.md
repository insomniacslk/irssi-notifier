# irssi-notifier
Desktop notifications with irssi and fnotify, over unauthenticated HTTPS. Client in Python, server in Go

# When to use it

If you run irssi in tmux/screen on some remote server, connect to it via ssh/mosh, and want to have client-side notifications.
WARNING: there is no authentication. Anybody can connect to the server and see what `fnotify` says. If you want authentication please send a PR.

# How to use it

* Install [fnotify for irssi](https://scripts.irssi.org/scripts/fnotify.pl) on the server (the machine running irssi)
* check out the server code on the server host, build it (`go build`) and run it (`./server -h` to see the available options)
* run the client code on the client (`./irssi-notifier.py`). It requires Python 3 and requests, and the server certificate.

By default the client uses `notify-send`. If you need another terminal notifier feel free to open an issue or send a PR.
