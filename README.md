# Linux exercise

## Exercise 1
### Bootable linux qemu script

This shell script creates and runs a Linux filesystem image using QEMU that will print “hello world” after startup. 
The script assumes that the system already has:
Internet access, 3GB free on the disk, installed qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils and the user who will launch the script must have superuser privileges.
If the kernel running on the machine has a particular configuration, umount the boot partition.

## Exercise 2
### Shred function in go 
Shred(filename) function overwrites the given file 3 times with random data and delete the file afterwards.
