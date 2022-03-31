create table if not exists users (
  id integer not null primary key, 
  username text not null,
  password text not null,
  createdAt integer not null
);

create table if not exists follows (
  id integer not null primary key, 
  lurkerUserId integer not null,
  followedUserId integer not null,

  foreign key (lurkerUserId) references users (id) on delete cascade
  foreign key (followedUserId) references users (id) on delete cascade
);

create table if not exists tweets (
  id integer not null primary key, 
  text text not null,
  userId integer not null,
  mediaPath string,
  createdAt integer not null,

  foreign key (userId) references users (id) on delete cascade
);

create table if not exists retweets (
  id integer not null primary key, 
  originalTweetId integer not null,
  userId integer not null,

  foreign key (userId) references users (id) on delete cascade
  foreign key (originalTweetId) references tweets (id) on delete cascade
);

create table if not exists reactions (
  id integer not null primary key, 
  originalTweetId integer not null,
  reaction text not null,
  userId integer not null,

  foreign key (userId) references users (id) on delete cascade
  foreign key (originalTweetId) references tweets (id) on delete cascade
);

-- Seed --

insert into users (username, password, createdAt) values (
  'admin', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into users (username, password, createdAt) values (
  'test', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into users (username, password, createdAt) values (
  'user', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into users (username, password, createdAt) values (
  'lurker', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into follows (lurkerUserId, followedUserId) values (
  1, 3
);

insert into follows (lurkerUserId, followedUserId) values (
  3, 1
);

insert into follows (lurkerUserId, followedUserId) values (
  4, 1
);


insert into tweets (userId, text, createdAt) values (
  1, 'First tweet ever', strftime('%s', 'now')
);

insert into tweets (userId, text, createdAt) values (
  3, 'Second tweet ever', strftime('%s', 'now')
);

insert into tweets (userId, text, createdAt) values (
  3, '3rd tweet ever', strftime('%s', 'now')
);

insert into retweets (originalTweetId, userId) values (
  1, 3
);

insert into retweets (originalTweetId, userId) values (
  1, 4
);

insert into reactions (originalTweetId, userId, reaction) values (
  1, 3, '👏'
);

insert into reactions (originalTweetId, userId, reaction) values (
  1, 4, '🎉'
);