airspace is a command line utility to provide an interactive shell for use with an ADS-B reciever

![airspace demo](https://raw.githubusercontent.com/vicgarcia/airspace/main/demo.gif)

airspace is intended for install on a Raspberry PI running FlightAware's PiAware configuration

Review and follow the full instructions from [FlightAware](https://flightaware.com/adsb/piaware/build)
- Assembling a PiAware device with a Raspberry Pi and ADS-B radio dongle
- Flash the PiAware OS image to your MicroSD card

After flashing the device follow instructions for [additional options](https://flightaware.com/adsb/piaware/build/optional)
- Enable SSH access to the PiAware device by creating an empty file `ssh` at the root of the MicroSD card
- After logging in via SSH, expand the filesystem to use the entire MicroSD card with `sudo raspi-config`

SSH to the device using `ssh pi@<device ip>` with the password `flightaware`

At this point your device will be ready to use with FlightAware

Visit the device IP in the browser and select the link to 'Claim this feeder to associate it with your FlightAware account' to complete setup with FlightAware

Install some basic stuff
```
sudo apt install vim git tmux
```

Add authorized keys for remote access to the pi user
```
cd ~/.ssh
touch authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

Disable login w/ password by configuring sshd with `sudo vim /etc/ssh/sshd_config`
```
...
PasswordAuthentication no
...
UsePAM no
...
PermitRootLogin no
...
```

Then configure tmux with `vim ~/.tmux.config`
```
new-session

set -g base-index 1
setw -g pane-base-index 1

set -g prefix C-a
unbind C-b

unbind-key Tab

bind \\ split-window -h
bind - split-window -v

bind -n C-p next
bind -n C-o prev
bind -n C-l select-pane -t :.+

set-option -ga terminal-overrides ',*:enacs@:smacs@:rmacs@:acsc@'

set-option -g status-style bg=black,fg=white
set-option -g window-status-current-style bg=white,fg=black
set-option -g window-status-current-style bg=white,fg=black
set-option -g pane-border-style bg=black,fg=white
set-option -g pane-active-border-style bg=black,fg=white
set-option -g window-status-format " #I "
set-option -g window-status-current-format " #I "
set -g status-left "| "
set -g status-right "#H |"
```

After configuring tmux, disconnect from the device ssh session and reconnect using tmux
```
ssh pi@<device ip> -t "tmux attach-session -d -t pi || tmux new-session -s pi"
```

Setup [FlightRadar24](https://flightradar24.com) feeder based on [instruction here](https://forum.flightradar24.com/forum/radar-forums/flightradar24-feeding-data-to-flightradar24/11792-beginner-feed-both-fr24-und-fa-with-raspberry-pi-3-model-b-flightaware-pro-stick)
```
sudo bash -c "$(wget -O - https://repo-feed.flightradar24.com/install_fr24_rpi.sh)"
```

Use the email address and sharing key associated with your FlighRadar24 account

During the install select these options
```
Enter your receiver type (1-7)$:4
...
Enter your connection type (1-2)$:2
...
Step 4.3A - Please enter your receiver's IP address/hostname
$:127.0.0.1

Step 4.3B - Please enter your receiver's data port number
$:30005

Step 5.1 - Would you like to enable RAW data feed on port 30334 (yes/no)$:no

Step 5.2 - Would you like to enable Basestation data feed on port 30003 (yes/no)$:no

Step 6 - Please select desired logfile mode:
 0 -  Disabled
 1 -  48 hour, 24h rotation
 2 -  72 hour, 24h rotation
Select logfile mode (0-2)$:2

Saving settings to /etc/fr24feed.ini...OK
Settings saved, please run "sudo systemctl restart fr24feed" to use new configuration.
Installation and configuration completed!
```

Restart fr24 for new configuration to take effect
```
sudo systemctl restart fr24feed
```

Check the status of everything
```
sudo systemctl status fr24feed
sudo systemctl status dump1090-fa
sudo systemctl status piaware
```

Install python 3.8, via [this guide](https://community.home-assistant.io/t/home-assistant-core-python-3-8-backport-for-debian-buster/234859)

Install python3 -distutils and -venv, which will be used by the python 3.8 install
```
sudo apt install python3-distutils python3-venv
```

Add the GPG key for the new repository
```
wget https://pascalroeleven.nl/deb-pascalroeleven.gpg
sudo apt-key add deb-pascalroeleven.gpg
```

Add the new repository with `sudo vim /etc/apt/sources.list`
```
sudo vim /etc/apt/sources.list
...
deb http://deb.pascalroeleven.nl/python3.8 buster-backports main
```

Update apt and install python 3.8
```
sudo apt update
sudo apt install python3.8 python3.8-venv python3.8-dev
```

Install the Go compiler
```
sudo apt install golang
```

Clone the airspace repo in `/opt`
```
cd /opt
sudo mkdir airspace
sudo chown -R pi:pi airspace
cd airspace
git clone https://github.com/vicgarcia/airspace.git .
```

Create the airspace config with `vim /opt/airspace/cmd/airspace/.config`
```
vim /opt/airspace/cmd/airspace/.config
...
AIRCRAFT_JSON_PATH=/run/dump1090-fa/aircraft.json
AVIATION_STACK_API_KEY=
```

Run `setup-db.sh` script to create the faa registration database
```
./setup-db.sh
```

Compile the airspace application
```
cd /opt/airspace/cmd/airspace
go build
```

Create launch script in `/usr/local/bin`
```
cd /usr/local/bin
sudo touch airspace
sudo chmod 755 airspace
sudo chmod +x airspace
vim airspace
...
#!/bin/bash
pushd /opt/airspace/cmd/airspace > /dev/null
/opt/airspace/cmd/airspace/airspace
popd > /dev/null
...
```

Call `airspace` over ssh
```
ssh pi@<piaware ip address> -t airspace
```

Add in `.zshrc`
```
alias airspace="ssh pi@<piaware ip address> -t airspace"
```
