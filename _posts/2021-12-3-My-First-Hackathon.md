---
layout: post
title:  "My First Hackathon"
date:   2021-12-03 18:00:00 -0600
categories: linux docker discord hackathon utdallas
description: We built a discord bot to facilitate on-campus ride sharing.
---

# My First Hackathon

This semester at UT Dallas I participated in my first hackathon, [HackUTD VII](https://hackutd.co/).  It was a great experience to meet new people in my sector, and find passion in a project I wouldn't have otherwise thought of.

We built a [Discord Bot](https://github.com/blakelawyer/CommuteBot) that helps Riders & Drivers on campus car pool.  We found this use case to be unique, in that students are a speical breed of customer: They're broke, don't have a car, and want to see the world around them. Especially when they're in a new place, like University! Otherwise, some students just don't have cars, especially those who are experiencing financial hardships or just uncomfortable driving in the craze that is Dallas/Fort Worth Traffic!

You can check out a couple youtube videos to see how it works, namely from the [rider's perspective](https://youtu.be/fs8H7vkDtdY) and the [driver's perspective](https://youtu.be/rW02s0dfZqs), or read about our HackUTD submission on [DevPost](https://devpost.com/software/cometcommute-discord-bot).

We chose to build this as a chatbot for a couple reasons, namely UX and popularity of Discord on our campus.  Studiets are showing that UX gets improved if a user is able to step through a workflow that looks like an instant message.  Additionally, Discord is used heavily by our student body, almost every single student org maintains a Discord server for their communities, so we believe this would scale well with our peers.

# The Stack

I'd like to talk a little about the technologies we used: as one judge mentioned "this is almost an enterprise grade stack." We used a few key cloud-based technologies to help accelerate our deployment, not worrying (too much) about the infrastructure.

* [MongoDB Atlas](https://www.mongodb.com/atlas/database)
* [Amazon ECS (Elastic Container Service)](https://aws.amazon.com/ecs/)
* [Terraform](https://www.terraform.io/)
* [Travis-CI](https://www.travis-ci.com/)
* [NodeJS](https://nodejs.org/en/)
* [DiscordJS](https://discord.js.org/)

## MongoDB Atlas

I wish I would've taken time during the competition to see if Atlas had a terraform provider, but nonetheless, Atlas saved the day. After we came up with our idea and stopped to eat lunch, I came back to my team struggling. They weren't sure how we were supposed to get a MongoDB server erected on our University's network in a way that we could all access it. I had the chance to save the day here, "Just use Atlas!" In about 15 minutes, we had a cluster deployed to AWS on Atlas' Free Tier, that met the needs of the rest of our time with the project. Given that this Chatbot has the potential to be used by almost everyone at UTD, we wanted to make sure we were choosing a scalable and maintainable database architecture.

## Amazon ECS

I use ECS day-to-day at work for [Granbury Solutions](https://www.granburysolutions.com). My team is in love with ECS, from ECS Anywhere to "it just works," we've never seen significant issues in ECS that made us question our choices. Naturally, to have a highly available, scalable Discord Bot, I made the choice that we should run our production bot on ECS. My team had little familiarity with Docker, so I quickly cooked us up a Dockerfile and made sure the container launched and functioned as expected locally. Again for scalability, but also for availability, we chose ECS as Fargate tasks are super-easy to scale, and can be architected for high availability.

## Terraform

I've been using Terraform more and more lately, and I can't get enough. Terraform is so extensible, it fits for almost every cloud provider/infrastructure use case I can come up with by myself. It has a huge ecosystem of tools, from [Checkov](https://www.checkov.io) and [TFLint](https://github.com/terraform-linters/tflint) to cooler tools like [Infracost](https://www.infracost.io/) and [Inframap](https://github.com/cycloidio/inframap). 

## Travis-CI

I use Travis-CI in my personal projects on GitHub, for it's simplicity (and it's German, like me!) Travis makes life easy for us, in that we're assured our code builds and passes tests before it gets to the production environment in ECS, every time. We can easily track the status of our build, and even expose the build logs to the public, to prove the project is truly open-source.

## NodeJS

While I consider my Node to be inferior to those who work using Node on a daily basis, I'd consider it "good enough to get the job done." I'm by no means a professional when it comes to writing Node, handling Promises, or publishing modules. However this project forced my hand, and got me more hands on with Node. Overall, I think this was a good experience for me to get comfortable in something I've been wanting to use more.

## DiscordJS

I found the DiscordJS API to be complicated, but sound. Once one can read through the forest that is their [API Docs](https://discord.js.org/#/docs/main/stable/general/welcome) and seek out enough examples, the logical structure starts to unwind. I would've liked to use a whiteboard during the Hackathon so we could draw the Forest out, and provide a visualization to the whole team.

I think if our team had more time, or came into the competition knowing more about the DiscordJS API, we could've built this as a webhook-based Bot instead of using the current workflow. A webhook-based Bot would allow us to build a series of Lambda functions to serve a REST API, that our Bot would use to receive events from Discord, instead of monitoring the Discord Gateway for new messages in-app.