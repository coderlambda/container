#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <linux/netlink.h>
#include <linux/rtnetlink.h>
#include <string.h>

struct request {
  struct nlmsghdr hdr;         //netlink msg header - man 7 netlink
  struct ifinfomsg ifinfomsg;  //netlink route ifinfomsg - man 7 rtnetlink
  struct rtattr ext_req __attribute__((aligned(NLMSG_ALIGNTO)));
  __u32 ext_filter_mask;
};

int BUFFER_SIZE_LARGE_ENOUGH = 8192;

void rtnl_print_link(struct nlmsghdr *h) {
  struct ifinfomsg *iface;
  struct rtattr *attribute;
  int len;

  iface = NLMSG_DATA(h);
  len = NLMSG_LENGTH(sizeof(*iface));

  for (attribute = IFLA_RTA(iface); RTA_OK(attribute, len); attribute = RTA_NEXT(attribute, len)) {
    switch(attribute->rta_type) {
    case IFLA_IFNAME:
      printf("Interface %d : %s\n", iface->ifi_index, (char *)RTA_DATA(attribute));
      break;
    default:
      break;
    }
  }
}

int get_netlink_socket() {
  int sock;
  if((sock = socket(AF_NETLINK, SOCK_RAW, NETLINK_ROUTE)) < 0) {
    return -1;
  }

  struct sockaddr_nl nladdr;
  memset(&nladdr, 0, sizeof(nladdr));
  nladdr.nl_family = AF_NETLINK;

  if(bind(sock, (struct sockaddr*)&nladdr, sizeof(nladdr)) < 0) {
    return -1;
  }

  return sock;
}

int main() {
  int sock = get_netlink_socket();
  if(sock < 0) {
    return -1;
  }


  struct request req;
  memset(&req, 0, sizeof(req));

  // build request
  req.hdr.nlmsg_len = sizeof(req);
  req.hdr.nlmsg_type = RTM_GETLINK;//get link info
  req.hdr.nlmsg_flags = NLM_F_REQUEST | NLM_F_DUMP;
  // req.hdr.nlmsg_seq = 1234; for custom usage
  req.ifinfomsg.ifi_family = AF_UNSPEC;

  req.ext_req.rta_type = IFLA_EXT_MASK;
  req.ext_req.rta_len = RTA_LENGTH(sizeof(__u32));
  req.ext_filter_mask = RTEXT_FILTER_VF;
  
  // send request;
  send(sock, (void *)&req, sizeof(req), 0);

  // if you don't receive all the message, check the buffer size.
  char reply[BUFFER_SIZE_LARGE_ENOUGH];
  int len;
  len = recv(sock, reply, sizeof(reply), 0); 
  if (len < 0) {
    return -1;
  }

  struct nlmsghdr *msg_ptr;
  msg_ptr = (struct nlmsghdr*)reply;

  while(NLMSG_OK(msg_ptr, len)) {
    if (msg_ptr->nlmsg_type == 16) { // here we just care about this type of msgs
      rtnl_print_link(msg_ptr); 
    } else {
      break;
    }
    msg_ptr = NLMSG_NEXT(msg_ptr, len);
  }

  return 0;
}
