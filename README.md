# srl-ndk-git

This agent extends SRL to act as a client to [github](https://github.com). It allows to create branches, commit files and create Pull-requests and use Githib as a change management system for the configuration files

## Installation

```bash
# install the rpm with yum without downloading the RPM
yum -y https://github.com/srl-wim/container-lab/releases/download/v0.4.0/container-lab_0.4.0_linux_amd64.rpm

# or when rpm is downloaded to the host
sudo rpm -i container-lab_0.4.0-next_linux_amd64.rpm
```

## Open items

* action: should be an atomic command without commit stay
* action yang: enum branch, commit, pull-requst
* yang: space in the name
* yang: make namespace variable -> not it is fixed to mgmt
* yang: store token in hash form
* local file -> now it is fixed
* Why is ygot not generating the struct so we can marshal the data
