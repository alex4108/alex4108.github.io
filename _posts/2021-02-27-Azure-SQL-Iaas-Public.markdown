---
layout: post
title:  "Public access to an Azure SQL IaaS Failover Cluster Instance"
date:   2021-03-12 17:03:00 -0600
description: A journey of publicly exposing a SQL FCI on Azure IaaS using Azure Load Balancer, Azure Firewall, and HAProxy
# categories: 
---

## The Legacy Environment

If you read my other post about [how I moved 4600 SQL DBs in 90 minutes](https://schittko.me/2021/02/04/SQL-Server-4600-DBs/), you'll see the [Azure SQL IaaS environment](https://schittko.me/2021/02/04/SQL-Server-4600-DBs/#the-iaas-environment) that post refers to.  

I'll go ahead and reiterate the environment, just for sanity:

* 2x Azure Virtual Machines (D16s_v4)
* 1x SQL Server Failover Cluster Instance in Windows Server Failover Cluster
* 1x Azure Load Balancer - Internal - Floating IP configured for the SQL Server Instance's IP address
* 1x Azure Firewall - To allow users of the greater internet to reach our SQL Server.

Before we started up the cluster, we realized that there's a flaw in [the how-to doc over on docs.microsoft.com](https://docs.microsoft.com/en-us/azure/azure-sql/virtual-machines/windows/failover-cluster-instance-vnn-azure-load-balancer-configure): The statement that a public load balancer can be used for the failover cluster instance seems false.  At least, I could never get it to work.  

There were a few attempts we tried, detailed below, before settling on Azure Firewall to perform DNAT to the internal load balancer.  The Azure Firewall wasn't a silver bullet though: paying $1.00 per hour for a dumb packet forwarder seemed like a waste of money.  Although the service did get us through this pickle, we wanted to see how we could reduce costs further.

The exchange with Azure Support over this flaw has gone on for 5 weeks now, with multiple outages as a result of us trying various workarounds, at least 5 sets of hands who have looked at the problem, and no ideas that can lead to us shutting down that expensive Azure Firewall.  Not to say Azure Support isn't awesome: in the simpler issues I've faced, they make quick work of solving my mess and making sure I understand how to avoid breaking something again.

### Take 1: Public Load Balancer with VMs in Backend Pool

So this seemed like the most straightforward attempt.  We created a Azure Public Load Balancer, created a backend pool, and added the VMs to the pool.

As we could see in Wireshark, this ended up not working, as the packets forwarded by the load balancer would be destined for the Host's IP address, not the SQL FCi's IP address.  We tried enabling Floating IP, but this didn't do much, it simply forwards the packets to the host, but also rewrites the TCP Destination address to be the public load balancer's IP.  

Naturally, the SQL Host is running the SQL Server instance on a cluster shared address, and when the host would receive packets, it would simply observe "hey, I don't have anything listening on that interface, so I'm dropping you".

So, how do we DNAT packets from the world, to end up at the correct destination address?

### Take 1.1: Public Load Balancer with Internal Load Balancer in Backend Pool

We used the same public load balancer as above, this time creating a new backend pool with a single IP Address target.  We configured the target to be the Internal Load Balancer, and this too proved unsuccessful: packets would just get dropped between the public load balancer and the internal.

Surely there has to be a way to forward these packets...

### Take 2: Azure Firewall

Well, Azure Firewall ended up doing the DNAT properly, at the cost of $1.00 per hour.  Seemed a bit expensive for a dumb packet rewriter service.  Sure, there are many other benefits to Azure Firewall, like IDS and IPS;  Given our use cases, IDS and IPS weren't on our to-do list.

Overall, Azure Firewall seems like a good solution when the use case fits.  The problem is here, we're using size 14 shoes from Azure, to dress our toddler SQL Server.  It's just overkill.

### Take 3: HAProxy

Ahh HAProxy.  [Used by (at least) 646 companies](https://stackshare.io/haproxy), most with a global footprint, HAProxy is the tried and true reverse proxy that's suitable for basically any use-case you can throw at it.

So, this one was straight forward: Given we have a 2-node cluster, I created 2 more VM's running Ubuntu 18.04 and placed them in the same VNet as the SQL Server instances.  From there, I wired the Public Load Balancer in Take 1 to use the HAProxy VM's as it's backend.

HAProxy's configuration is dirt-simple, just set up a listener to forward to the internal load balancer's IP address (10.90.1.200)

```
listen l1
    bind 0.0.0.0:1433
    mode tcp
    timeout client 180000
    timeout server 180000
    server ilb 10.90.1.200:1433
```

And like magic, we now have a "dumb packet forwarder" at $0.09 per hour (x2 for high availability, so really $0.18) that allows our SQL Server Failover Cluster Instance on Azure IaaS to be publicly exposed.  

If we take the savings from Azure Firewall versus HAProxy, we're looking at a cost reduction of $0.82/hour, just to forward packets.  I hope with these savings, my CFO can buy me a coffee (or two)!

### Take 4: Tell Windows Server Failover Cluster about the Public Load Balancer's Address.

So this one was interesting to me.  It goes against everything I know as a junior network engineer, but hey, it works.

In [this doc](https://docs.microsoft.com/en-us/azure/azure-sql/virtual-machines/windows/failover-cluster-instance-vnn-azure-load-balancer-configure#configure-cluster-probe) there's a powershell snippet to configure "the health probe" of the WSFC Role for SQL Server.  This not only creates a health probe available on 59999/tcp, but also modifies the IP address that the SQL Server instance listens on.

```
$ClusterNetworkName = "<Cluster Network Name>"
$IPResourceName = "<SQL Server FCI / AG Listener IP Address Resource Name>" 
$ILBIP = "<n.n.n.n>" 
[int]$ProbePort = <nnnnn>

Import-Module FailoverClusters

Get-ClusterResource $IPResourceName | Set-ClusterParameter -Multiple @{"Address"="$ILBIP";"ProbePort"=$ProbePort;"SubnetMask"="255.255.255.255";"Network"="$ClusterNetworkName";"EnableDhcp"=0}
```

Turns out, you can set `<n.n.n.n>` to the Public IP associated with the Azure Load Balancer's Frontend configuration.  Provided Floating IP is enabled on the load balancer rule, this enables publicly routable connectivity to the SQL Server instance.

I'd like to credit Ashan Nanayakkara from Microsoft Azure Networking support who had the bright idea, "set this to the public IP," that eventually lead to this discovery.