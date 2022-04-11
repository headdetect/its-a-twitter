
insert into users (username, password, createdAt) values (
  'admin', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into users (username, password, createdAt) values (
  'test', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into users (username, password, createdAt) values (
  'basic', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into users (username, password, createdAt) values (
  'lurker', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into users (username, password, createdAt) values (
  'chuck', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
);

insert into users (username, password, createdAt) values (
  'lily', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', strftime('%s', 'now')
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

-- Mark tweets as created in the past --

insert into tweets (userId, text, createdAt) values (
  1, 'First tweet ever', strftime('%s', 'now') - 1000
);

insert into tweets (userId, text, createdAt) values (
  3, 'Second tweet ever', strftime('%s', 'now') - 999
);

insert into tweets (userId, text, createdAt) values (
  2, '3rd tweet ever', strftime('%s', 'now') - 998
);

insert into tweets (userId, text, createdAt) values (
  1, 'This is a different tweet', strftime('%s', 'now') - 997
);

insert into tweets (userId, text, createdAt) values (
  3, 'Lorem ipsem dolor', strftime('%s', 'now') - 996
);

insert into tweets (userId, text, createdAt) values (
  3, 'beep boop boop boop beep bop beepy boop', strftime('%s', 'now') - 995
);

insert into tweets (userId, text, createdAt) values (
  4, 'One and only tweet', strftime('%s', 'now') - 994
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

insert into reactions (userId, tweetId, reaction) values (
  2, 1, '👏'
);

insert into reactions (userId, tweetId, reaction) values (
  3, 1, '🎉'
);

insert into reactions (userId, tweetId, reaction) values (
  4, 1, '🎉'
);