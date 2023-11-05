# BorderOS

## Hardware supported

- EdgeRouter Lite (ERL-3)

## Design

Based on [Gokrazy](https://gokrazy.org/), ideas from [Router7](https://router7.org/).

## Installation

### EdgeRouter Lite

- Unplug the internal flash drive, or obtain a new one
- `gok overwrite --full /dev/sdN`
  - Using `gok` from the fork of gokrazy-tools
- Connect to the serial console of the EdgeRouter and power it on
- Before the OS loads press any key to drop into u-boot
- Update the `bootcmd` variable
  - You may wish to save the default first, e.g. `setenv oldbootcmd $bootcmd`

  ```
  setenv bootcmd fatload usb 0 $loadaddr vmlinuz; bootoctlinux $loadaddr coremask=0x3 root=/dev/sda2 rootwait init=/user/brd mtdparts=phys_mapped_flash:512k(boot0),512k(boot1),64k@1024k(eeprom)
  ```
- Store this: `saveenv`
- Reboot the router

## Testing

The `brd` init program can be tested on a regular Linux system without needing the full gokrazy image.

Partition a system with with a disk layout as follows:

- 500MB EFI
- ext4 `/`
- ext4, unmounted, identically-sized to `/` partition
- ext4, unmounted, to be used as gokrazy `/perm`

This is the layout used by gokrazy.

Build `cmd/brd` and install it as `/brd`.

Create `/config` and add a `config.cue` with the system configuration.

Setup GRUB2 to boot using `brd`. On Debian this can be done by editing `/etc/grub.d/40_custom` and setting up a boot entry as follows:

```
#!/bin/sh
exec tail -n +3 $0
# This file provides an easy way to add custom menu entries.  Simply type the
# menu entries you want to add after this comment.  Be careful not to change
# the 'exec tail' line above.
menuentry 'BorderOS' {
	load_video
	insmod gzio
	insmod part_gpt
	insmod ext2
	echo	'Loading Linux 6.1.0-13-arm64 ...'
	linux	/boot/vmlinuz-6.1.0-13-arm64 root=/dev/nvme0n1p2 init=/brd
	echo	'Loading initial ramdisk ...'
	initrd	/boot/initrd.img-6.1.0-13-arm64
}
```

Replacing the kernel lines and root device to match your own system.

Then run `sudo update-grub`.
