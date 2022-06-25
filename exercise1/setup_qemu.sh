#!/bin/bash


WDIR=$(pwd)
KERNEL_LIST=$(ls /boot/vmlinuz*)
KERNEL=$(echo $KERNEL_LIST | awk '{print $1}')
INITRD_LIST=$(ls /boot/initrd*)
INITRD=$(echo $INITRD_LIST | awk '{print $1}')

get_kernel ()
{
	wget -c https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-5.18.5.tar.xz
	tar -xf linux*
	rm linux-5.18.5.tar.xz
	cd linux*
	make defconfig
	RES=$?
	if [[ $RES != 0 ]]
	then
		echo "Something failed, exit."
		exit 1
	fi
	make
	cp arch/x86/boot/bzImage $WDIR/
}


echo "Check disk space, more than 3GB of free space are required."
FREE_SPACE=$(df -k . | awk 'FNR == 2 {print $4}')
if [[ $FREE_SPACE -lt 3145728 ]]
then
	echo "Not enough space"
	exit 1
fi	
echo "Update environment."
apt update >> /dev/shm/log.txt 2>&1
apt -y install make >> /dev/shm/log.txt 2>&1
apt -y install bzip2 >> /dev/shm/log.txt 2>&1
apt -y install libncurses-dev flex bison openssl libssl-dev dkms libelf-dev libudev-dev libpci-dev libiberty-dev autoconf >> /dev/shm/log.txt 2>&1
mkdir src
cd src

echo "Download busybox"
wget -c http://www.busybox.net/downloads/busybox-1.35.0.tar.bz2 >> /dev/shm/log.txt 2>&1
RES=$?
if [[ $RES != 0 ]] 
then
	echo "Unable to download busybox, check your internet connection"
	exit 1
fi
tar -xf busybox-1.35.0.tar.bz2 >> /dev/shm/log.txt 2>&1
cd busybox-1.35.0
echo "Make defconfig and install busybox"
{
make defconfig 
make install
} &>/dev/shm/log.txt

cd $WDIR
dd if=/dev/zero of=rootfs.ext4 bs=1M count=3072 >> /dev/shm/log.txt 2>&1

echo "Create ex4 filesystem"
mkfs.ext4 -F rootfs.ext4 >> /dev/shm/log.txt 2>&1
mkdir tmpfs
mount -t ext4 -o loop rootfs.ext4 tmpfs/ >> /dev/shm/log.txt 2>&1
echo "Populate filesystem"
cp -av src/busybox-1.35.0/_install/* tmpfs/ >> /dev/shm/log.txt 2>&1
cd tmpfs
mkdir -pv {bin,dev,sbin,etc/init.d,proc,sys/kernel/debug,usr/{bin,sbin},lib,lib64,mnt/root,root} >> /dev/shm/log.txt 2>&1
cp -av /dev/{null,console,tty,sda} dev/ >> /dev/shm/log.txt 2>&1
cp -rP /lib/* lib/ >> /dev/shm/log.txt 2>&1
cp -rP /lib64/* lib64/ >> /dev/shm/log.txt 2>&1
ln -s sbin/init init >> /dev/shm/log.txt 2>&1
echo -e "#!/bin/sh\n
	echo -e \"\nHELLO WORLD\"
	" >> etc/init.d/rcS
chmod +x etc/init.d/rcS >> /dev/shm/log.txt 2>&1
cd $WDIR
mknod tmpfs/dev/tty1 c 4 1 >> /dev/shm/log.txt 2>&1
mknod tmpfs/dev/tty2 c 4 2 >> /dev/shm/log.txt 2>&1
cd $WDIR

umount tmpfs

sudo qemu-system-x86_64 -m 1024 -nographic -kernel $KERNEL -initrd $INITRD -append "root=/dev/sda rw console=ttyS0,115200" -hda rootfs.ext4
RES=$?

if [[ $RES != 0 ]]
then
	get_kernel
	cd $WDIR
	sudo qemu-system-x86_64 -m 1024 -nographic -kernel bzImage -append "root=/dev/sda rw console=ttyS0,115200" -hda rootfs.ext4
	RES=$?
fi

if [[ $RES != 0 ]]
then
	echo "Something going wrong."
	exit 1
else
	exit 0
fi

