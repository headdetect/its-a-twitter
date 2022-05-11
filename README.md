# It's a Twitter

This is part of the "It's A" series. 
A short project series that aims to showcase my technical abilities by recreating
some famous applications/problems and breaking them down into their hard parts.

# Why

 This application was created to showcase my ability to build a fully-featured project from scratch. The goal is to understand the challenges that an app like Twitter would face, and demonstrate how these challenges could be solved. This "It's-a" project was chosen because Twitter has a clearly defined feature-set with interesting problems to tackle. The goal is not to re-create Twitter in its entirety, but to create a project with low risk and ease of use that is reminiscent of Twitter. 

# Scope

There are three main objectives for this project:
1. Design a system from start to finish &mdash; from concept to code. 
2. Explain the engineering decisions exercised.
3. Describe how to efficiently scale the application if it were the size of Twitter. 

This is just a recreation of the essential parts that make Twitter work. The scope has been cut to accommodate a shorter timeline of around 40-50 hours of work. This gives us
the ability to focus on the more specific details one might not consider when building an application meant for a public audience.

Given the goals of this project and timeline, there were a handful of engineering decisions that only make sense in this specific scenario. Throughout the application, there are comments and notes that convey the intent of a solution and how the approach taken for this purpose would differ if it were a solution for an enterprise system.

# Features

Below is a list of product facing features (things users directly interact with) and non-functional features (positive consequences of the system design)

## Product Features

- Send out short text snippets (tweet)
- Retweet
- Post media
- Emoji reactions (not a twitter thing, but I like it)
- Personalized Timelines
- Log in and manage account (display name, profile pic, and password)

## Non-functional features

- Short TTFB, less than 200ms for api calls (except logging in because [bcrypt takes a long time](https://auth0.com/blog/hashing-in-action-understanding-bcrypt/#Motivation-Behind--bcrypt-) on purpose)
- Able to run application in Heroku
- Front end able to be hosted using github pages

# Engineering Decisions

## Dependencies

I wanted to accomplish the ambitions of the project with as little dependencies as possible. In a real world scenario, dependencies (while they do give the ability to enable rapid development) can cause a down stream of issues if the dependencies break. As an application becomes more dependent on other providers to keep it functional, the potential for the application to become unavailable becomes unavoidable. Most of the largest companies with applications that tout millions of Daily Active Users have created their own building blocks and services so they can ensure the most critical of the blocks are as _cohesive and available_ and possible. As a test of ability, I've decided to reduce the amount of dependencies as much as possible and treat this project like it has a future.

Given this, there will be dependencies added in the interest of brevity and in the interest of showcasing my ability to utilize those dependencies as efficiently as possible.

## Media Store

In a large application like twitter, the assets uploaded by users (videos, pics, galleries, etc) would be stored independent of the application and the applications's filesystem. There are many services that accomplish this goal of separation of concerns, things like Amazon's S3 is a perfect example.

In this application, we'll keep uploaded media in the same filesystem and store it in of the `./media` (or whatever is set in `MEDIA_PATH`) folder and serve it from there as requested. There is no server side caching stood up in this application in the interest of simplicity. All assets are served with headers to instruct the web browsers to: "store this asset for as long as possible".

## Database

Much like the media store in the section above, the database is a very simple system. It's a SQLite database stored in the same filesystem that the executable runs in. While it may not be a production grade database, it still gave me the ability to feature my ability to accomplish tasks using raw SQL. (Some of these can be found in the model structures [here](https://github.com/headdetect/its-a-twitter/tree/master/api/model))

# Future Scaling

As stated in several instances throughout this write-up, if this were a true production grade application with the number of users twitter sees on a daily basis, these are a number of ways it would need to scale. While there might be hundreds of changes that would need to be made in order to reach the same level of availability, these highlight some of the larger additions that _could_ be introduced in order to scale at the capacity needed to sustain this theoretical user base size.

## Database

An application with the popularity of Twitter would have huge chunks of data flowing all over the organization and would require a database to handle that. There are an infinite amount of ways you could divide and conquer all the data that might be streaming into the application. But by possessing only a single database server would be an insurmountable bottle neck for an application like of that caliber.

Beyond the caching, increasing physical resources (ram, cpu, disk space, etc), and creating read-only replicas there's not a lot of ways to scale a database.

As such, something like a database shard would help accomplish the task of distributing the load that Twitter-like traffic would generate while also maintaining data integrity. Database sharding allows the data to be shared across different servers and allows the database to scale horizontally instead of vertically.

At a high level view, the database structure is maintained across the database swarm, but the data (the rows in the database) are split amongst several different instances. There is no right answer when it comes to picking how data is partitioned and in most cases, it's application specific. The goal is to reach an equal balance of data between all the replicas so that the load is evenly distributed.

## Load Balancing

Load balancing is the process of distributing load across multiple instances of the application. It's a crucial and frequently required feature for web applications that receive any respectable amount of traffic. The biggest capability that is unlocked when propping up a load balancer is the ability to scale up and down the number of instances as the flow of user traffic changes. Allowing your application to handle any amount of traffic it might receive.

As designed from the beginning, this demo app is not fit to be place behind a load balancer and scaled up with multiple instances. As such, there are a few changes that would need to be made in order to enable that ability. The database is set-up to be a single file sqlite database stored in the file system of each of the instances. This prevents data from being shared among different application instances and would partition the data to be per-instance, which means things like session, tweets, profiles, etc would not persist between different instances. Obviously this is a major blocker in terms of scalability for a popular application, but in our case for a demo app, it is acceptable. 

## Caching

The first line of defense for an overstretched application is usually improving/implementing some level of caching. The types of caching, the TTL, and scale of caching is situation dependent, but for an application that is similar to this one, there are two main types of caching that could be explored in order to reduce load. Bundle caching and media/asset caching.

### Bundle Caching

Because the application is designed as an SPA of sorts, the front end code could be aggressively cached. If the JavaScript bundle were to be loaded with a hash tied to the specific build version, the bundle could be cached indefinably. This would reduce the TTFB given the only requirement is to load the compressed javascript application pack. After the pack is loaded and all the scripts have been JIT compiled, we can bring up some intermediary loading screen while assets, tweets, and other content is fetched and loaded from different sources (if they weren't already being pre-fetched) 

### CDN/Media Caching

Media would be a large portion of this application and having to load the assets from a single source would take forever. We would need to store the different medias in different edge nodes around the world so it could be loaded quicker. This would allow users in different parts of the world to load the assets as if it were being served from a few miles away (because it would be)

The absolute best solution would be to duplicate all media being uploaded to every edge node. But alas, money doesn't grow on trees and housing that amount of data in every edge node is just a logistical nightmare. So as we break it down, there are a few major issues with this solution. Can't store that much data and don't have the bandwidth to duplicate **every** asset to **every** edge node. 

So a potential solution for this would be to create a classification for each user to determine how far their audience reach goes and how many users might see a piece of media. Someone in a small village in Europe probably isn't going to be looking at that picture a local politician posted in Utah. We can keep the image local to the audience that it will *most likely* reach. 

Subsequently, if that politicians image sparks some rage and it goes viral, we would have more incentive to cache that image in more and more edge nodes. We would need to employ a strategy of storing images in some edge nodes only after someone requests it.

# How to build the app

[See the build docs](BUILD.md)