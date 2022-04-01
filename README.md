# It’s a Twitter

This is part of the “It’s A” series. A short project series that recreates the apps I use.

Here’s a list of other projects that are a part of these series:

- TODO

# Why

As I’m growing in my career I wanted to push myself to understand the hard parts that make up the apps I use on a daily bases. The goal is to understand the challenges that building an app like Twitter would face, and demonstrates how I would solve them. The outcomes I’m hoping to gain are two fold: Gain a better understanding of the system design of a social networking platform and learn Go.

# Scope

Obviously we don’t want to recreate the full Twitter experience as a result, this is just a recreation of the essential parts that make Twitter work. The scope has been cut to accommodate my, self assigned, timeline of a single week. This includes the infamous lack of editing a tweet. Mostly because I don’t want to implement it.

That gives me time to do a deep dive into the interesting challenges and solutions while also not completely consuming my life. This is a project series of course. I have to pace myself.

# Features

This project will feature some of twitters main functions. Below is a list of product facing features (the things that users will be able to use) and non-functional features (features that are not user facing).

## Product Features

- Send out short text snippets (tweet)
- Retweet
- Post media
- Emoji reactions (not a twitter thing, but I like it)
- Personalized Timelines
- Log in and manage account (display name, profile pic, and password)

## Non-functional features

- TTFB < 200ms
- Generated timelines w/30 second window.

# Engineering Decisions

## Dependencies

I wanted to accomplish the ambitions of the project with as little dependencies as possible. In a real world scenario, dependencies (while they do give the ability to enable rapid development) can cause a stream of issues down the line. As an application becomes more dependent on other providers to keep it functional, the potential for the application to become unavailable becomes unavoidable. Most of these large companies with millions of Daily Active Users will need to create their own building blocks to ensure the most cohesive and available system relating to their solutions. As a test of ability, I’ve decided to reduce the amount of dependencies as much as possible and treat this project like it has a future.

However, there will be decisions made in the interest of brevity. In most situations should be called out in the code where they’re being implemented. But for choices made for the overall application, I’ll list them below.

- For the frontend, I used React because it’s stable, lightweight, and easy to use. It will be the largest front end dependency but should remain as the only major dependency.

## Data

Because I don’t want to throw money at this project we will be storing media locally as opposed to some object storage provider like Amazon’s S3. They’ll be stored & cached in the server’s filesystem. 

## Validation
There is barely any. Why? Because it doesn't really matter here.

## Auth Scheme

# Design

TODO - Insert API Spec

TODO - Insert Schema

# Future Scaling

If this were a true production grade application with the number of users twitter sees on a daily basis, these are the types of ways it would need to scale.

## Database

There would be thousands and thousands of tweets going out every second of every day. 

TODO

## Storage

TODO

## Load Balancing

TODO

## Caching

TODO

### Bundle Caching

Because the application is designed as an SPA of sorts, the front end code could be aggressively cached. If the JavaScript bundle were to be loaded with a hash tied to the specific build version, the bundle could be cached indefinably.

### CDN/Media Caching

Media would be a large portion of this application and having to load the assets from a single source would take forever. We would need to store the different medias in different edge nodes around the world so it could be loaded quicker. This would allow users in different parts of the world to load the assets as if it were being served from a few miles away (because it would be)

The absolute best solution would be to duplicate all media being uploaded to every edge node. But alas, money doesn’t grow on trees and housing that amount of data in every edge node is just a logistical nightmare. So as we break it down, there are a few major issues with this solution. Can’t store that much data and don’t have the bandwidth to duplicate **every** asset to **every** edge node. 

So a potential solution for this would be to create a classification for each user to determine how far their audience reach goes and how many users might see a piece of media. Someone in a small village in Europe probably isn’t going to be looking at that picture a local politician posted in Utah. We can keep the image local to the audience that it will *most likely* reach. 

Subsequently, if that politicians image sparks some rage and it goes viral, we would have more incentive to cache that image in more and more edge nodes. We would need to employ a strategy of storing images in some edge nodes only after someone requests it.