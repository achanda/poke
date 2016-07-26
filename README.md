This is a port scanner that I wrote for fun. The primary design goal here is simplicity. Currently, it has the following features:
* Scan a hostname or IP
* Use goroutines to scan
* Supports IPv6
* Tries to guess service names (using a library that I wrote)

## Caveats ##
* UDP support is not fully implemented yet
* TCP SYN and UDP scans use raw sockets, so these modes need root permission
* CLI error reporting has not been tested well

## Usage ##
Since this uses goroutines, the host system might run out of file descriptors. Please increase ulimit.
```bash
ulimit -n 100000
```
### Scanning a single IP ###
```bash
# ./poke -host 159.203.228.218 -ports 1:65535 -scanner s
Scanning 159.203.228.239...
5672/tcp open amqp
80/tcp open http
3306/tcp open mysql
3260/tcp open iscsi-target
22/tcp open ssh
4369/tcp open epmd
5000/tcp: open
53/tcp open domain
35357/tcp: open
37729/tcp: open
```
### Scanning a host ###
```bash
# ./poke -host google.com -ports 1:100 -scanner s
Scanning 216.58.192.14...
80/tcp open http
```
### Scanning a CIDR ###
```bash
# ./poke -host 159.203.228.239/30 -ports 1:100 -scanner s
Scanning 159.203.228.236...
80/tcp open http
22/tcp open ssh
21/tcp open fsp
Scanning 159.203.228.237...
22/tcp open ssh
80/tcp open http
Scanning 159.203.228.238...
80/tcp open http
Scanning 159.203.228.239...
22/tcp open ssh
80/tcp open http
53/tcp open domain
```
### IPv6 scanning ###
```bash
# ./poke -host 2604:a880:1:20::9f9:9001 -ports 1:100 -scanner s
Scanning 2604:a880:1:20::9f9:9001...
53/tcp open domain
22/tcp open ssh
80/tcp open http
```
## Developing ##
This uses glide to manage packages. The `update` makefile target installs and gets all dependencies if they are not present. Tests can be run by:
```
make test
```
Building
```
make build
```
