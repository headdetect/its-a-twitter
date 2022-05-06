create table if not exists users (
  id integer not null primary key, 
  username string not null unique,
  password string not null,
  email string not null,
  profilePicPath string unique,
  createdAt integer not null
);

create table if not exists follows (
  userId integer not null,
  followedUserId integer not null,

  foreign key (userId) references users (id) on delete cascade
  foreign key (followedUserId) references users (id) on delete cascade

  primary key (userId, followedUserId)
);

create table if not exists tweets (
  id integer not null primary key, 
  userId integer not null,
  text text not null,
  mediaPath string unique,
  createdAt integer not null,

  foreign key (userId) references users (id) on delete cascade
);

create table if not exists retweets (
  userId integer not null,
  tweetId integer not null,

  foreign key (userId) references users (id) on delete cascade
  foreign key (tweetId) references tweets (id) on delete cascade

  primary key (userId, tweetId)
);

create table if not exists reactions (
  userId integer not null,
  tweetId integer not null,
  reaction string not null,

  foreign key (userId) references users (id) on delete cascade
  foreign key (tweetId) references tweets (id) on delete cascade
);