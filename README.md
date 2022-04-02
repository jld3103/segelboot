# Segelboot

Segelboot is a wrapper around `efibootmgr` and is designed to manage your EFISTUB UEFI boot entries.  

## Installation

`go install github.com/jld3103/segelboot@latest`

## Usage

See `example.conf` for an example configuration.  
You need to know what `efibootmgr` does and how it works.  

When you have created your config file at `/etc/segelboot.conf`, simply run `segelboot` with root privileges,
and it will create all the necessary entries.  
On subsequent runs segelboot will recreate all entries that it recognizes. You should not change the section names
in the config file, otherwise segelboot won't recognize the existing entry
(although it will create a new entry that will overwrite the old entry if it is on the same partition).  

To delete all entries segelboot ever created run `segelboot --delete` with root privileges.

## History

For a while I had a simple bash script for recreating the entries, but it can get quite annoying when
you change where the disk is located in your system and the script was pretty naive in general and could have easily
messed with my system (which luckily never happened). At some point I just decide to make it a proper tool.  

I know there already is a tool that is designed to automatically create the entries
(I forgot the name and can't find it anymore), but it didn't work for me and just did some weird stuff.

## Name

Segelboot is the German word for sailing boat, and I had the idea from the U-Boot project.
