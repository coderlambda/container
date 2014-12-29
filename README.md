container
=========

to study how to build a container step by step

## container basics
To build a container you should know following things:

* namespace
  To be continued
* cgroup
  To be continued
I aufs
  To be continued

## network 

The simplest way to manage network interface is using ip command.
But you may want to use netlink APIs insted.

Here is some infomation about netlink API:

Basics:
* man 7 netlink
* man 3 netlink
* man 7 rtnetlink
* man 3 rtnetlink
* http://www.linuxjournal.com/article/7356?page=0,0
* http://iijean.blogspot.com/2010/03/howto-get-list-of-network-interfaces-in.html
* http://mwnnlin.blogspot.com/2010/12/get-list-of-network-interfaces-using.html

If you fallowing the articles up, sometimes you will get an empty list, it's because the change of the netlink API.

Here is the key change:
* http://permalink.gmane.org/gmane.linux.network/220770
* http://lxr.free-electrons.com/source/include/linux/if_link.h

There is a example code in the src named "getLinkList.c".It tested on 3.13.0-32-generic







