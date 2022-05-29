
# Avly Trader
[![MPLv2 License](https://img.shields.io/badge/license-MPLv2-blue.svg?style=flat-square)](https://www.mozilla.org/MPL/2.0/)
[![Version](https://img.shields.io/docker/v/9tmark/avly-trader?style=flat-square&sort=semver)](https://hub.docker.com/r/9tmark/avly-trader)
[![Docker Pulls](https://img.shields.io/docker/pulls/9tmark/avly-trader.svg?style=flat-square)](https://hub.docker.com/r/9tmark/avly-trader)
[![Docker Image Size](https://img.shields.io/docker/image-size/9tmark/avly-trader.svg?style=flat-square&sort=date)](https://hub.docker.com/r/9tmark/avly-trader)
## Disclaimer
Please check [LICENSE](LICENSE) before contributing or using *Avly Trader*. Please keep in mind: Use at own risk. No copyright or trademark infringement intended.

## About the project
Over the years, the **MetaTrader** (registered trademark) by **MetaQuotes Ltd** became a commonly used software for automated trading setups and trading in general. With *Avly Trader* you will be able to put your trading setup into the cloud, using your own provider, infrastructure and settings. There are no limitations whatsoever–you can do anything you would be able to do on your local trading setup. On top, you enhance the scalability, can monitor it as you like and be in control all the time.

## Special features
In case you are ready and suited for your advanced cloud trading setup, consider the following benefits at your side:
- Flexible component for your cloud infrastructure: Build and deploy an image which fulfills your requirements–or just download a [released version](https://github.com/9tmark/avly-trader/releases) and [run it](#manual-usage) wherever
- No restrictions in terms of general usage or integration of external tools (whatever runs, is fine)
- Be responsible for your own system security... well, for some people this *is* a benefit
- Feel at home: Use your trading setup the familiar way. Login to the UI via [VNC](https://en.wikipedia.org/wiki/Virtual_Network_Computing) and do your thing, as always
- There's no default way implemented for sharing files (scripts, config, indicators etc.) between instances, but there are several options, depending on your infrastructure

## Getting started
**Prerequisites:**
- a VM or cluster infrastructure or whatever for deployment
- a **docker**-ready system in case you'd like to build the docker images locally
- **go** version 1.18 in case you'd like to build the CLI from source 
- a **third-party** folder:
	- content:
		- a file called **mt5setup.exe** which is the official installer from MetaQuotes Ltd
		- a file called **wine-mono-7.1.1-x86.msi** you will get at [WineHQ](https://wiki.winehq.org/Mono)
		- a file called **wine_gecko-2.47-x86_64.msi** you will get at [WineHQ](https://wiki.winehq.org/Gecko)
		- a file called **winetricks** you will get [here](https://github.com/Winetricks/winetricks/blob/master/src/winetricks)
	- this folder is required to be present as a volume for the docker container to run and set itself up (see [compose file](resources/02-run/compose/docker-compose.yml))
- not mandatory but recommended: an empty **logs** folder for the container to store persistent logs (see [compose file](resources/02-run/compose/docker-compose.yml))

### Quick installation
On your target system, create a `.avly` folder inside the home directory. This folder should contain a separate log folder for every instance and a single `third-party` folder. Make sure all necessary files are inside it (see [prerequisites](#getting-started)).
Open a terminal:
```sh
$ wget https://raw.githubusercontent.com/9tmark/avly-trader/main/resources/02-run/compose/docker-compose.yml
```
Use your preferred editor to edit the compose file. Let's assume it's `vim` for now:
```sh
$ vim docker-compose.yml
```
Scroll to the `volumes` section:
```yml
volumes:
      - /etc/timezone:/etc/timezone:ro
      # - <path to logs on host>:/var/log/avly-trader
      # This line is required (see README):
      # - <path to third-party on host>:/opt/third-party
```
Make sure it look something like this (see [prerequisites](#getting-started)):
```yml
volumes:
      - /etc/timezone:/etc/timezone:ro
      - /home/<username>/.avly/mt5001:/var/log/avly-trader
      - /home/<username>/.avly/third-party:/opt/third-party
```
Great! The hard part is done. Now it's time to start trading:
```sh
$ docker-compose pull app
```
```sh
$ docker-compose up -d
```
You can execute the command  without detach flag (`-d`) or observe the log files. There you will see if and when the container is ready. **Please wait!** The boot time can vary. As soon as the container is ready, you can start trading.

Finally you can connect this instance, using the VNC client of your choice. On your local computer it would be `localhost:55900`. **ATTENTION! Make sure you DO NOT expose any of these ports or directories to the public internet.** A simple solution could be using a VM at your preferred cloud provider. Usually, by default, they are only reachable via SSH (Port 22), secured with the public key method. Most VNC clients will allow you to establish a VNC connection via an [**SSH tunnel**](https://askubuntu.com/questions/1090177/use-remmina-1-2-0-with-ssh-tunneling).

Could work like this:
1. **SSH** connection on **public IP**, port **22**, authentication via **public key**
2. Use the **SSH** connection as a tunnel: **VNC** connection on **localhost**, port **55900**

It's an important topic. Please do some research and pick a secure approach which fits for you.

### Manual usage
If you're looking for a more customizable way to go, see the `help` output of the `avly` command:
```
  Usage of avly:
  -c
  -clean-up
        dispose remains of target process
  -d
  -drain
        shut down VNC server
  -e
  -enter
        run startup routine as container process
  -f
  -fledge
        (safely) pull up VNC server
  -l
  -launch
        (safely) launch target executable
  -m
  -mute
        mute output unless error occurs
  -p
  -prepare
        verify perquisites for a workstation to work properly
  -s
  -stop
        stop target process
```
What basically happens inside the container, is the execution `avly -e`. This command is **NOT recommended** to be executed on a personal computer.
