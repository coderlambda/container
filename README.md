container
=========

to study how to build a container step by step

## container basics
To build a container you should know following things:

* namespace

    man 7 namespace

A namespace wraps a global system resource in an abstraction that makes it appear to the processes within the namespace that they have their own isolated instance of the global resource.  Changes to the global resource are visible to other processes that are members of the namespace, but are invisible to other processes.  One use of namespaces is to implement containers.
  
The releated systemcalls are not releated in go, so if you want to build a container, you need wrap the syscall youself.

e.g setns, see src/go/system/setns_linux.go (copy from docker)

References:
* https://blog.jtlebi.fr/2013/12/22/introduction-to-linux-namespaces-part-1-uts/
* http://stackoverflow.com/questions/16977988/details-of-syscall-rawsyscall-syscall-syscall-in-go 
* http://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/commit/?id=7b21fddd087678a70ad64afc0f632e0f1071b092
  
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







