---
layout: post
title:  "How we used XCP-ng to run Rosetta@Home"
date:   2021-02-04 02:38:00 -0600
categories: linux vm xcp-ng Rosetta@Home
---

# The Project

We were tasked by our [Professional Responsibility](https://catalog.utdallas.edu/2020/undergraduate/courses/cs3162) course to "execute a project that benefits more than 5 people."  Minimum requirement was 8 hours of work.  All in all, [Nate](https://github.com/DeathBringer12345) and I spent about 12 hours together getting this project to fruition.

This seemed a little shocking at first, we were given complete oversight of our projects, and were to run our ideas by our professor for approval prior to starting.

## How we landed at Rosetta@Home

So the first idea was pretty simple, set up a "private cloud" in the apartment's closet where we could host VMs and have a dedicated space to run projects for us, friends, and the world.  Unfortunately, this really only directly benefited the two of us, until it dawned: [Rosetta@Home](https://boinc.bakerlab.org/) is one of the many [BOINC](https://boinc.berkeley.edu/) projects out there.  If we could get this "private cloud" off the ground, and install a VM purely dedicated to Rosetta@Home, we could meet the 8 hour requirement, and help the world.

The server still runs, and is part of the [Cure-4-Cancer](https://boinc.bakerlab.org/rosetta/team_members.php?teamid=1070&offset=0&sort_by=expavg_credit) team.  You can check the box's stats [here](https://boinc.bakerlab.org/rosetta/show_user.php?userid=2344453)

# Procurement

I already had a PowerEdge R710 I scored of [r/homelabsales](https://reddit.com/r/homelabsales) for free, provided it was going to a good cause.  I think the original owner would've been proud to see how it helps the world now, considering it started as a test bed for me to fiddle with some bare-metal.  The box has 2 Xeon X5650's, and 64GB of RAM.  We allocated 18 vCPU's to the Rosetta@Home VM, leaving 6 vCPUs for Nate & I to run our projects with.  Frontier provides us with "upto" 500mbit up and down, so we had plenty of bandwidth to boot.

# Execution

## Physical Topology
The topology of the physical environment is dirt simple.  The box sits on a shelf in a closet in the apartment, and has a 16 port unmanaged switch above it, with a run going to the Frontier (residential Fiber ISP) router.  Because I really don't like residential/consumer grade router/firewall devices, we chose to apply a DMZ rule on the Frontier router, thus forwarding all ports from the apartment's edge to... the pfSense Virtual Machine running on the metal!

So how did we get to pfSense in a VM?  Easy, really.

## Setting up the "private cloud" VM Host

1. Install xcp-ng on the bare metal
1. Set up a management interface on an IP that won't be pfSense's Edge address.
1. Set up a vSwitch.  The idea here is the Edge is only exposed to pfSense, and all VMs will use the vSwitch to route through pfSense.  Thus, we have packet filtering, HAProxy, LetsEncrypt, OpenVPN, and a few other services running in pfSense to make the environment easier to manage remotely.
1. Set up pfSense, and make sure to disable [tcp checksums](https://support.citrix.com/article/CTX212540) in Xen, otherwise suffer a gnarly performance hit to the throughput of 16Kbit's a second!  All we ended up needing was `xe pif-param-set uuid=$PIFUUID other-config:ethtool-tx="off";`
1. Set up an OpenVPN installation in pfSense, so we can tunnel in to the virtual LAN where our VMs run
1. Set up an SMB Mount on the local windows desktop, so we can get a few ISO images ready to provision VMs.

From there, we just ran.  It took about 6-7 hours to work through the initial bugs, like getting pfSense off the ground ([Why do they have an expired Root CA in OS' default certificate package?  And why is that Expired CA still signing the pkg.pfsense.org repo?  Good question.](https://serverfault.com/questions/1019348/certificate-expiration-pfsense-wont-update-ddns)), and getting that Frontier router to trying to assign the public IP to pfSense.  We wanted to keep the apartment LAN as is, and just have the private "cloud" accessible.

## Setting up Ubuntu 20, Docker, and Boinc.

Getting Ubuntu installed was cake, and the feature that lets us import public keys from GitHub directly to OpenSSH during installation saved me from ever needing a password.  Except, of course, to [disable password prompts in sudo](https://askubuntu.com/questions/147241/execute-sudo-without-password)!  Even better, [installing Docker](https://docs.docker.com/engine/install/ubuntu/) was super-quick thanks to the folks there giving us the shell snippets we needed to run.

The final steps of getting the BOINC client and Ubuntu were fairly trivial, and pretty well documented.  We chose the [BOINC in docker](https://github.com/BOINC/boinc-client-docker) approach, because why not?  Docker makes life easy, if the documentation for the project is good.