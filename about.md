---
layout: page
title: About 
permalink: /about/
---

Here, you'll find a collection of works about everything I found time to do, outside of class work.

I'll just go ahead and brain dump everything I can think of that I've touched in the last 20 years that I can recall.

## Linux

I use Ubuntu as a dev environment under [Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install-win10) [v2!](https://docs.microsoft.com/en-us/windows/wsl/compare-versions)

Given that I've used Ubuntu since 12.04 LTS, I've also used Debian 9 and 10 in Virtual Machines.

I'm familiar with CentOS 6 through 8 (RIP...) as we use this as the primary OS for deployments in my current professional role.  This directly translates to Amazon Linux 2, which is RHEL based like CentOS.

## Windows

Yeesh, my first computer was Windows 98.  I've been on the Windows train from 98, through XP, Vista, 7, 8, 8.1, and finally 10 (long live!) 

## Containers

You can see from my repositories in my [GitHub account](https://github.com/alex4108) that I'm a huge fan of using docker to ship images to production.  I love how docker is lightweight, easy to use, and keeps things moving in the development cycle.

I'm a fan of [docker-compose] for it's simplicity.

I've yet to really get under the hood with k8s, but I have deployed a few helm charts.

I haven't used RHEL's podman, and haven't heard much about it, other than it exists.

## Ansible

I've used ansible once in a [personal project](https://github.com/alex4108/wp-ansible-kubernetes/tree/main/mariadb-galera), to deploy a [MySQL Galera cluster](https://galeracluster.com/products/)

## Databases

I've used MongoDB, MySQL, and SQLite in personal projects.  In the workplace, we use [Amazon RDS](https://aws.amazon.com/rds/) to host (mostly) AWS Aurora MySQL clusters.  

Microsoft SQL Server in the workplace seems to be nice for the legacy windows apps.  I like how much documentation exists, and how [docs.microsoft.com](https://docs.microsoft.com) stays up to date, even if you're running a 10 year old SQL Server edition.

## Virtual Machines

I use Hyper-V on my Windows 10 machine to run VMs locally.  Right now, there's just a Windows Server 2019 VM I used to test [AstroInstaller](https://github.com/alex4108/AstroInstaller).  AstroInstaller is some powershell to automatically install the [Astroneer]() [Dedicated Server]().  The idea was to use this project as a launch configuration for a Virtual Machine template in Azure, however I haven't quite gotten there yet...

In the professional world, we use [Amazon EC2](https://aws.amazon.com/ec2/) and [Azure VMs](https://azure.microsoft.com/en-us/services/virtual-machines/) to run our production workloads.  There are a few cases where we prefer [DigitalOcean](https://www.digitalocean.com/), for absolutely-not-critical-at-all workloads that could disappear in an instant with no repercussions.  Aside from VMs in the cloud, I also operate a few [Windows Server Failover Clusters](https://docs.microsoft.com/en-us/windows-server/failover-clustering/failover-clustering-overview) running Hyper-V VMs, for local/on-prem services.  The first week I had this in production, I was giddy, I could finally rest well knowing our praised on-prem services would reboot in the event of a host failure, which was all to frequent for that time.

I also have a spare piece of metal, a PowerEdge R710 with 2 Xeon X5650's and 64GB of RAM.  It runs [xcp-ng](https://xcp-ng.org/), a community fork of [XenServer](https://xenserver.org/).  Watch out for a blog post highlighting a project I did with [Nate Montgomery](https://github.com/DeathBringer12345), one of my classmates at UTD.


## Python

My first encounter was Python 2.6 in the 9th grade.  I quickly enjoyed the syntax and extensibility of the language.  I've continued to use Python in projects far and wide.  

My favorite Python projects to-date:
* [Approova, a Discord Bot](https://github.com/alex4108/Approova)
* [Stock Price Alerter, because I watch Slack more than my phone during the day](https://github.com/alex4108/stock_price_alert_system)
* [SoundCloud Likes Downloader (defunct)](https://github.com/alex4108/scLikesDownloader) - I wrote this for a friend who DJ'ed back in high school.  His external hard drive with all his music died, with about 1000 songs, and no backups!  Luckily, most of his library consisted of his Likes in SoundCloud, thus this one was born.

## NodeJS / JavaScript

I have my apathety towards node and the concurrency under the hood.  I frequently produce race conditions in Node, because I still haven't fully grasped [Promises](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) yet.  I like NodeJS, but not `node_modules`, it feels like there's a lot of unnecessary bloat in there.

Regardless, I'm a fan of Node and will probably continue to use Express in my future projects.  Maybe one day I'll learn promises.

My favorite NodeJS projects to-date:
* [Scoreboard](https://github.com/alex4108/scoreboard).  A simple webapp to keep track of the score during board games.

### React

Aside, I'm not much of a frontend developer or UX designer, I've never really been creative brain to that level.  Regardless, that hasn't stopped me from _trying_ to learn React.

You can see my only React example here:
* [Scoreboard](https://github.com/alex4108/scoreboard)

## PHP

My first server-side language!  No, really!  I was a huge fan of PHP5 back in the day, but the reality is there are much more powerful and secure languages that have come into the production realm since I taught myself.  While I try to stay out of PHP these days, it always finds it's way back to me...

Some PHP Projects:
* [RingCentral Live Boards](https://github.com/alex4108/ringcentral-live-boards).  This one produces pseudo-live dashboards by repeatedly polling the RingCentral API for call data in our organization, and allows our teams to see their Queue's status in a much more unified fashion.
* [GSheets Sync](https://github.com/alex4108/GSheetsSync).  When I worked in Support, we posted the schedule on google sheets every few weeks.  This one ripped the google sheet and made Google Calendar appointments for all my shifts, making it easy to share my work schedule with the loved ones.
* [StreamMachine](https://github.com/alex4108/StreamMachine).  Built during the COVID-19 pandemic, it was meant as a way for my inner circle to be able to stream content in 4K, without paying for premium services at Discord or TwitchTV.  

## Bash

I _really_ like bash because of it's ease of use.  For the simpler things, it can get the job done quickly, and without much fuss.

Some bash scripts I've written, and like:
* [docker-puller](https://github.com/alex4108/docker-puller), because otherwise my containers will go away on DockerHub.
* [spot-reaper](https://github.com/alex4108/spot-reaper), so we can fire notifications when [AWS Spot Instances](https://github.com/alex4108/spot-reaper) get reclaimed or otherwise get a termination signal.

## Java

While I don't have any of my own Java projects published, we use Java heavily in the workplace for REST APIs using [Spring Framework](https://spring.io/)

## C/C++

WHY?  Sure, I've written some, and could probably hack around in an existing project if I needed.  Definitely not my cuppa, but I did learn it at school.

In my opinion, there are far more languages which make memory management safer, like Java or Rust.

## Rust

Eventually, I want to write some Rust.  One day.

## Continuous Integration & Deployment, Pipelines, and Source Control!

In the workplace, we just have a [Jenkins](https://www.jenkins.io/) instance wired up to [BitBucket Hooks](https://confluence.atlassian.com/bitbucketserver/using-repository-hooks-776639836.html) to run tests and deployments as needed.

Personally, [Travis-CI](https://travis-ci.com) makes my life really easy with managing _so many_ projects.  You can check out all my deployments on my [Travis-CI Page](https://travis-ci.com/alex4108).  

I'm really glad Travis-CI chose to be part of the [GitHub Education Pack](https://education.github.com/pack), otherwise I probably would've gone on to learn [GitHub Actions](https://github.com/features/actions).

For SCM, I've only ever used Git.

# Cloud!

This one gets it's own section, because why not?

## Amazon Web Services

In my current role with a Point of Sale SaaS provider, we heavily use AWS to operate development, test, and production workloads.  I'll go ahead and list every AWS service I've used and worked with here.

* RDS

Because managing databases without it is just too much.  RDS gives us the availability and fault tolerances we need to operate at scale.

* Systems Manager (SSM)

_Mostly_ use this for parameter store, to keep secrets out of ECS Task Definitions.

* EC2

Per-customer isolated virtual machines, for PCI compliance.

* ECS + ECR

To run our web service APIs

* S3

For static website hosting, and just general object storage.

* Lambda + API Gateway

Because serverless is the future!

* Kineses Firehose

For big-data pipelines so we can aggregate POS metrics and data into our Redshift warehouse

* Redshift

See above!

* EFS

When we need an NFS mount

* DynamoDB

This one gets used very lightly because of costs, it holds some configuration data for some of our web services to load into cache upon init.

* ElastiCache

Because managed services for Redis and Memcached is the way!  Build cattle, not pets!

* VPC

Really?  Without a VPC you get nowhere.  Regardless, NAT Gateways are not Internet Gateways, and VPN tunnels come in handy when you need a tunnel onsite.  Similarly, VPC peering makes for an easy way to have _one_ jumphost to rule them all!

* Route 53

Because DNS needs to be availabile, everywhere, and update, fast!

* CloudFormation

Honestly, not a huge fan.  Prefer Terraform.  CloudFormation's "Delete and Create" update functionality really killed me a few times getting started with it.

* CloudWatch

Metrics, Metrics, Metrics!  Metrics give us valuable insights about how our services are operating, be it AWS provided metrics or our own custom metrics.

Also, CloudWatch Logs.

* Ground Station

If I ever needed a satellite link, I would probably use this.  Seems quicker than provisioning a satellite array...

* Elastic Load Balancing

Because Multi-AZ architectures are the way.  You don't want a single AZ outage to offline your services, do you?

* Certificate Manager

When I can't use acme.sh, like in the context of [Application Load Balancers](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) I use ACM.

## Azure!

It's not number two in my heart, but it's number two in my skillset.  As I grow, I'll use Azure more for my online services.

* Azure VMs

Probably the most useful service I've found yet.  With built-in's like [Tardigrade](https://azure.microsoft.com/en-us/blog/improving-azure-virtual-machines-resiliency-with-project-tardigrade/) and the possibility of adding resources without a requiring a VM reboot, I definitely prefer Azure VM's over Amazon EC2.

Plus, Azure had native VM [Backup](https://docs.microsoft.com/en-us/azure/backup/backup-azure-vms-introduction) before [AWS Backup](https://aws.amazon.com/backup/) existed.  +1 for Azure!

* Azure Storage Accounts

Blob storage is _much_ cheaper than S3 if used appropriately.  Plus with [Azure Files](https://azure.microsoft.com/en-us/services/storage/files/), I can provision SMB mounts with a few button clicks or a Terraform deployment.  Dang, that was easy!  Time for coffee...

Alas, I haven't worked much with Queues in Azure Storage Accounts.  One day I might pick that up...

* Azure Load Balancing

I really enjoy how LB's in Azure offer the ability to pool outbound connections from the VMs, as well as incoming connections, unlike Amazon's Elastic Load Balancing.  

* Azure Automate

Okay this one is really cool, I used the python3 beta runtime once to automate cleaning of scheduled backups.  I must say, getting off the ground with python3 because of the weird way Automate handles dependencies was tricky, but once I figured it out, it was all smooth sailing.  I like Azure Automate!

* Azure Traffic Manager

The service is comparable to Amazon's Route 53 [route policies](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy.html), and became super useful during one of our org's migrations into Azure that I headed up.  The use case was simple: we wanted to incrementally shift production load to Azure, so we could monitor for failure and correct as needed, before opening the flood gates.  

* More to come!

I want to get more hands on experience using Azure App Services, Azure Functions, and CosmosDB.  These seem like the tools I need to evolve my Serverless on Azure skills.

# Free Time

Wow, if you got here, you must be committed to learning more about me!

## Rock Climbing

When I'm not at the desk, typing away on some blog post or project, you can find me at [Summit Climbing, Yoga and Fitness](https://www.summitgyms.com/) when the weather isn't in my favor, or at [Horshoe Canyon Ranch](https://horseshoecanyonduderanch.com/), climbing rocks.  Sometimes, I look off the cliff and think "dang, I feel like a mountain goat right now, how did I even get up here!?"  I've found rock climbing to be my favorite hobby I've discovered yet, as it really requires all my focus.  I can't be half paying attention to climbing, distracted by something that happened previously in regards to work, school, or one of the other numerous stresses that exists as a part life.  Rock climbing lets me use my sharp skills for something more than just code.

## Video Games

And like anyone with a computer, I play some video games.  Most notably, right now, some games I've enjoyed are

* [Rocket League](https://www.rocketleague.com/)
* [Rogue Company](https://www.roguecompany.com/), because third person just has a nice perspective
* [Sid Meier's Railroads](https://store.steampowered.com/app/7600/Sid_Meiers_Railroads/)
* [Deep Rock Galactic](https://www.deeprockgalactic.com/) "Rock and Stone!"
* [Counter-Strike: Global Offensive](https://store.steampowered.com/app/730/CounterStrike_Global_Offensive/) (I started on Source)

## Photography

I don't have any samples to share, but I enjoy taking photos of nature when time permits!