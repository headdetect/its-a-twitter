
insert into users (username, password, email, createdAt, profilePicPath) values (
  'ashley', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'ashley@example.com', strftime('%s', 'now'), 'u-ashley.jpg'
);

insert into users (username, password, email, createdAt, profilePicPath) values (
  'bray', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'bray@example.com', strftime('%s', 'now'), 'u-bray.jpg'
);

insert into users (username, password, email, createdAt, profilePicPath) values (
  'chuck', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'chuck@example.com', strftime('%s', 'now'), 'u-chuck.jpg'
);

insert into users (username, password, email, createdAt, profilePicPath) values (
  'diane', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'diane@example.com', strftime('%s', 'now'), 'u-diane.jpg'
);

insert into users (username, password, email, createdAt, profilePicPath) values (
  'elon', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'elon@example.com', strftime('%s', 'now'), 'u-elon.jpg'
);

insert into users (username, password, email, createdAt, profilePicPath) values (
  'lily', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'lily@example.com', strftime('%s', 'now'), 'u-lily.jpg'
);


insert into follows (userId, followedUserId) values (
  1, 3
);

insert into follows (userId, followedUserId) values (
  2, 1
);

insert into follows (userId, followedUserId) values (
  2, 3
);

insert into follows (userId, followedUserId) values (
  3, 1
);

insert into follows (userId, followedUserId) values (
  4, 1
);

-- Mark tweets as created in the past so we can verify orderings --

-- 1
insert into tweets (userId, text, createdAt) values (
  1, 'First tweet ever', strftime('%s', 'now') - 20000000
);

-- 2
insert into tweets (userId, text, createdAt, mediaPath) values (
  3, 'Second tweet ever', strftime('%s', 'now') - 10000000, 't-00.gif'
);

-- 3
insert into tweets (userId, text, createdAt) values (
  2, '3rd tweet ever', strftime('%s', 'now') - 5000000
);

-- 4
insert into tweets (userId, text, createdAt) values (
  1, 'This is a different tweet', strftime('%s', 'now') - 2500000
);

-- 5
insert into tweets (userId, text, createdAt, mediaPath) values (
  3, 'Lorem ipsem dolor', strftime('%s', 'now') - 1250000, 't-01.jpg'
);

-- 6
insert into tweets (userId, text, createdAt) values (
  3, 'beep boop boop boop beep bop beepy boop', strftime('%s', 'now') - 625000
);

-- 7
insert into tweets (userId, text, createdAt) values (
  5, 'One and only tweet', strftime('%s', 'now') - 312500
);

-- 8
insert into tweets (userId, text, createdAt) values (
  4, 'One and only tweet', strftime('%s', 'now') - 312500
);

-- 9
insert into tweets (userId, text, createdAt, mediaPath) values (
  6, 'One and only tweet', strftime('%s', 'now') - 312500, 't-02.jpg'
);

-- 10
insert into tweets (userId, text, createdAt) values (
  5, 'One and only tweet', strftime('%s', 'now') - 312500
);

insert into retweets (userId, tweetId) values (
  3, 1
);

insert into retweets (userId, tweetId) values (
  4, 1
);

insert into retweets (userId, tweetId) values (
  3, 7
);

insert into retweets (userId, tweetId) values (
  3, 3
);

insert into retweets (userId, tweetId) values (
  1, 6
);

insert into retweets (userId, tweetId) values (
  6, 2
);

insert into retweets (userId, tweetId) values (
  6, 4
);

insert into retweets (userId, tweetId) values (
  6, 1
);

insert into reactions (userId, tweetId, reaction) values (
  2, 1, 'shocked'
);

insert into reactions (userId, tweetId, reaction) values (
  3, 1, 'party'
);

insert into reactions (userId, tweetId, reaction) values (
  4, 1, 'party'
);