create table if not exists users (
  id integer not null primary key, 
  username text not null,
  password text not null,
  createdAt integer not null
);

create table if not exists follows (
  id integer not null primary key, 
  userId integer not null,
  followedUserId integer not null,

  foreign key (userId) references users (id) on delete cascade
  foreign key (followedUserId) references users (id) on delete cascade
);

create table if not exists tweets (
  id integer not null primary key, 
  userId integer not null,
  text text not null,
  mediaPath string,
  createdAt integer not null,

  foreign key (userId) references users (id) on delete cascade
);

create table if not exists retweets (
  id integer not null primary key, 
  userId integer not null,
  tweetId integer not null,

  foreign key (userId) references users (id) on delete cascade
  foreign key (tweetId) references tweets (id) on delete cascade
);

create table if not exists reactions (
  id integer not null primary key, 
  userId integer not null,
  tweetId integer not null,
  reaction text not null,

  foreign key (userId) references users (id) on delete cascade
  foreign key (tweetId) references tweets (id) on delete cascade
);