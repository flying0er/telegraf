# Iptables Plugin

The iptables plugin gathers packets and bytes counters for rules within a set of table and chain from the Linux's iptables firewall.

Rules are identified through associated comment. Rules without comment are ignored.

The iptables command requires CAP_NET_ADMIN and CAP_NET_RAW capabilities. You have several options to grant telegraf to run iptables:

* Run telegraf as root. This is strongly discouraged.
* Configure systemd to run telegraf with CAP_NET_ADMIN and CAP_NET_RAW. This is the simplest and recommended option.
* Configure sudo to grant telegraf to run iptables. This is the most restrictive option, but require sudo setup.

### Using systemd capabilities

You may run `systemctl edit telegraf.service` and add the following:

```
[Service]
CapabilityBoundingSet=CAP_NET_RAW CAP_NET_ADMIN
AmbientCapabilities=CAP_NET_RAW CAP_NET_ADMIN
```

Since telegraf will fork a process to run iptables, `AmbientCapabilities` is required to transmit the capabilities bounding set to the forked process.

### Using sudo

You may edit your sudo configuration with the following:

```sudo
telegraf ALL=(root) NOPASSWD: /usr/bin/iptables -nvL *
```

### Configuration:

```toml
  # use sudo to run iptables
  use_sudo = false
  # defines the table to monitor:
  table = "filter"
  # defines the chains to monitor:
  chains = [ "INPUT" ]
```

### Measurements & Fields:


- iptables
    - pkts (integer, count)
    - bytes (integer, bytes)

### Tags:

- All measurements have the following tags:
    - table
    - chain
    - ruleid

The `ruleid` is the comment associated to the rule.

### Example Output:

```
$ iptables -nvL INPUT
Chain INPUT (policy DROP 0 packets, 0 bytes)
pkts bytes target     prot opt in     out     source               destination
100   1024   ACCEPT     tcp  --  *      *       192.168.0.0/24       0.0.0.0/0            tcp dpt:22 /* ssh */
 42   2048   ACCEPT     tcp  --  *      *       192.168.0.0/24       0.0.0.0/0            tcp dpt:80 /* httpd */
```

```
$ ./telegraf -config telegraf.conf -input-filter iptables -test
iptables,table=filter,chain=INPUT,ruleid=ssh pkts=100i,bytes=1024i 1453831884664956455
iptables,table=filter,chain=INPUT,ruleid=httpd pkts=42i,bytes=2048i 1453831884664956455
```
