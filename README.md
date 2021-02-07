# ssh_honeypot

You can use this package to reduce the impact of brute force login attemts on ssh, this package imitate an ssh server, and on each login attempt pause the execution for a short time, and log the origin of the attempt on a separate log file.


### Installation

1- Download and extract the binary: `wget https://github.com/Shadi/ssh_honeypot/releases/download/v0.1.0/ssh_honeypot_0.1.0_linux_amd64.tar.gz -q -O- | tar xz`

2- Move the executable to a location in your $PATH: `mv ssh_honeypot /usr/local/bin`

3- Run it: `ssh_honeypot`

### Flags:

* `-p`: Port to use, default is 22 
* `-l`: Path of file to use for writing login attempts, default is `./attempts.log`
* `-c`: Max concurrent attempts, default is 20
* `-w`: Pause duration  n seconds before failing a login attempt, default is 15

### Usage

1- Change the default port that is being used by your ssh server, usually this is specified in the file `/etc/ssh/sshd_config`, search  for the entry port, and change the value to a different port, it might be commented out if you are still using the default, after that restart the ssh service: `systemctl restart ssh`

2- If you are using a firewall make sure you allow access to the new port, for instance if you specify the port 2242 in step(1) and you are using ufw firewall you can use the command: `ufw allow 2242/tcp`, from now on when you login you need to specify the new port using ssh `-p` flag: `ssh root@address -p NewPort`

3- run the package `ssh_honeypot`, this will use the default package configuration, it will pause execution for 15 seconds between attempts and write attempts to the log file `./attempts.log`

### Automatically Start at Boot

If you are using systemd then you can use the unit in [systemd/sshp.service](systemd/sshp.service) to automatically start ssh_honeybot on boot:

1- Copy the file to `cp systemd/sshp.service /etc/systemd/system`
2- Reload systemd: `systemctl daemon-reload`
3- Enable and start the service: `systemctl enable --now sshp.service`