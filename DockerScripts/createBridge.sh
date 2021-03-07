#!/bin/bash
#

echo "Hello, $(whoami)!"
echo "Linking Mac VLAN from ${0}"

# https://collabnix.com/2-minutes-to-docker-macvlan-networking-a-beginners-guide/
# https://blog.oddbit.com/post/2018-03-12-using-docker-macvlan-networks/
# https://dockstarter.com/advanced/macvlan/
# https://gist.github.com/xirixiz

# Create a mac0 bridge network attached to the physical ens18, and add the ip range scope
sudo ip link add mac0 link ens18 type macvlan mode bridge
# Specify part of the eth0 scope you'd like to reserve for mac0
sudo ip addr add 10.0.50.222/32 dev mac0
# Bring up the mac0 adapter
sudo ip link set mac0 up
# Route local traffic to mac0
sudo ip route add 10.0.50.0/24 dev mac0

# https://linuxconfig.org/how-to-run-script-on-startup-on-ubuntu-20-04-focal-fossa-server-desktop
# https://unix.stackexchange.com/questions/47695/how-to-write-startup-script-for-systemd/47715#47715