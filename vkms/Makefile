obj-m := vkms.o
vkms-objs := vkms_drv.o vkms_plane.o vkms_output.o vkms_crtc.o vkms_gem.o vkms_composer.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(shell pwd) modules
	EXTRA_CFLAGS="-g -DDEBUG"
clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(shell pwd) clean

