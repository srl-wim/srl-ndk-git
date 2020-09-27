# srl-ndk-git

This agent extends SR-Linux to act as a client to [github](https://github.com). It allows to create branches, commit files and create Pull-requests and use Githib as a change management system for the configuration files of SRL devices.
This SW package extends both the CLI with python and the Agent is written in [golang](https://golang.org).

## Github

### Create a github account

Given we interact with github you should be in posession of a github account.

[setup github account](https://github.com)

### Create a github token

You should create a github token, which is used to authenticate the client. The following procedure describe the steps:

[create github token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)

### Setup a github repo

The github client interact with github within a repository. As such you should setup a gihub repository with the following procedure.

[setup github repo](https://docs.github.com/en/github/getting-started-with-github/create-a-repo)

The repo name we will be used as a basis to store the files you commit to the repo.

## Installing the srl-ndk-git agent

### Yum install from Internet

Below is a procedure to install the agent using yum, but there are other methods that can be used. For the procedure below the SRL instance should have access to the Internet.

```bash
login to the SRL instance ssh admin@<ip address>
from the command prompt execute bash
sudo yum install https://github.com/srl-wim/srl-ndk-git/releases/download/v0.1.0/srl-ndk-git_0.1.0_linux_amd64.rpm -y
```

Example:

```
[henderiw@srlinux-2 ~]$ ssh admin@172.19.19.11
admin@172.19.19.11's password:
Last login: Sun Sep 20 04:45:21 2020 from 172.19.19.1
Using configuration file(s): []
Welcome to the srlinux CLI.
Type 'help' (and press <ENTER>) if you need any help using this.
--{ running }--[  ]--
A:wan2# bash
bash-4.2$ sudo yum install https://github.com/srl-wim/srl-ndk-git/releases/download/v0.1.0/srl-ndk-git_0.1.0_linux_amd64.rpm -y
```

### Yum install with local download

Another alternative is downloading the image locally and installing it locally using the following command

```
sudo yum install srl-ndk-git_0.1.0_linux_amd64.rpm
```

### Python extensions

Given we use python extensions for the cli you might have to logout and login from the cli

```
quit
```

## Configuration

Now that your agent is installed in the system we have to activate it in the system. We do this in the following way.

### Loading the agent

login into the system

```
ssh admin@<ip address> 
```

First we need to load the agent:

```
/ tools system app-management application app_mgr reload
```

When you show the applications running on the system the agent should be visible

```
A:wan2# show system application
  +-------------------+--------+--------------------+-------------------------+--------------------------+
  |       Name        |  PID   |       State        |         Version         |       Last Change        |
  +===================+========+====================+=========================+==========================+
  | aaa_mgr           | 1538   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.385Z |
  | acl_mgr           | 1547   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.385Z |
  | app_mgr           | 1480   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.596Z |
  | arp_nd_mgr        | 1556   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.385Z |
  | bfd_mgr           |        | waiting-for-config |                         |                          |
  | bgp_mgr           |        | waiting-for-config |                         |                          |
  | chassis_mgr       | 1565   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.386Z |
  | dev_mgr           | 1504   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:07.212Z |
  | dhcp_client_mgr   | 1578   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.386Z |
  | dnsmasq-mgmt      | 306956 | running            |                         | 2020-09-19T06:06:02.122Z |
  | fib_mgr           | 1589   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.387Z |
  | gnmi_server       | 1826   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:09.226Z |
  | idb_server        | 1529   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:07.458Z |
  | isis_mgr          |        | waiting-for-config |                         |                          |
  | json_rpc          | 1831   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:09.263Z |
  | l2_mac_learn_mgr  | 1598   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.387Z |
  | l2_mac_mgr        | 1612   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.388Z |
  | l2_static_mac_mgr |        | waiting-for-config |                         |                          |
  | lag_mgr           | 1621   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.388Z |
  | linux_mgr         | 1630   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.388Z |
  | lldp_mgr          |        | waiting-for-config |                         |                          |
  | log_mgr           | 1639   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.389Z |
  | mcid_mgr          | 1649   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.389Z |
  | mgmt_server       | 1663   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.389Z |
  | mpls_mgr          |        | waiting-for-config |                         |                          |
  | ndk-git           | 333407 | running            | v20.6.1-286-g118bc27b34 | 2020-09-19T06:41:54.103Z |
  | net_inst_mgr      | 1673   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.390Z |
  | oam_mgr           | 1688   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.390Z |
  | ospf_mgr          |        | waiting-for-config |                         |                          |
  | plcy_mgr          |        | waiting-for-config |                         |                          |
  | qos_mgr           | 1830   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:09.239Z |
  | sdk_mgr           | 1703   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.391Z |
  | sshd-mgmt         | 2083   | running            |                         | 2020-09-13T17:36:11.924Z |
  | static_route_mgr  |        | waiting-for-config |                         |                          |
  | supportd          | 1489   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:07.061Z |
  | xdp_cpm           | 1719   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.391Z |
  | xdp_lc_1          | 1729   | running            | v20.6.1-286-g118bc27b34 | 2020-09-13T17:36:08.392Z |
  +-------------------+--------+--------------------+-------------------------+--------------------------+
```

You see the ndk-git agent appear in the application list. It has a PID of 333407 in this example. If you defined in the YML file that the ndk-git agent should have waited for configuration, the PID would not have been allocated since there was no configuration in the system and hence the agent process would not have started.

### Configuring DNS

DNS needs to be configured to ensure the github client finds the github api server. We use google dns in this example but other DNS servers can be used.

```
enter candidate
set / system dns network-instance mgmt
set / system dns server-list [ 8.8.8.8 8.8.4.4]
commit stay
```

I sometimes have to also update the local dns entries like this

```
enter candidate
set system dns host-entry api.github.com ipv4-address 140.82.121.6
commit stay
```

### Configuring the agent

Next step is configuring the git client. In order for the client to operate it should be configured with the following parameters:

* owner or organization: depending on how you created the repo you should use either owner or organization. This field should be configured with the name of the github owner or organization
* repo: name of github repository
* token: id of the token
* branch: name of the branch that will be created in the repo
* file: name of the file that should be used for the file that will be committed in the repo
* author: name of the Author that will be used during commit/pull-request
* author-email: email of the Author that will be used during commit/pull-request

Example:

```
enter candidate
set / git-client organization srl-wim
set / git-client repo srl-lab-config-store
set / git-client filename wan1-config.json
set / git-client author "Wim Henderickx"
set / git-client author-email wim.henderickx@gmail.com
set / git-client branch test-branchJE
set / git-client token ****
commit stay
```

### Github actions

In order to perform a github action and workflow you can execute the following sequence of commands:

Create a branch

```
/ tools git branch
```

Commit; the file we use to commit to github is the startup file.

```
/ tools git commit "first release"
```

Once you are ok with the updates you can create a pull request

```
/ tools git pull-request "first release"
```

### Aliasing


## State

The github client stores the operational state of the transactions.

```
A:wan1# info from state / git-client
    organization srl-wim
    owner ""
    repo srl-lab-config-store
    file wan1-config.json
    token *************************
    author WHenderickx
    author-email wim.henderickx@gmail.com
    branch test-branch2
    action pull-request
    oper-state up
    statistics {
        success 2
        failure 0
    }
--{ candidate shared default }--[ git-client ]-
```

## Logging

Information that the agent is providing is also send to /var/log/srlinux/stdout/<agentname>.log and can be sent to syslog, etc.


## Open items

* CommitMessage, PRmessage, etc -> how to handle this better versus hard coding
* Validation before sending to the RPC server
* yang: store token in hash form
* Enabling this through a proxy
* Aliasing

## Ongoing
* to be tested: network-instance

## To be discussed
* We use the startup file as source
* Why is ygot not generating the struct so we can marshal the data
