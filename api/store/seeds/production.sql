-- Users
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (1, 'ashley', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'ashley@example.com', 'u-ashley.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (2, 'bray', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'bray@example.com', 'u-bray.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (3, 'chuck', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'chuck@example.com', 'u-chuck.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (4, 'diane', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'diane@example.com', 'u-diane.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (5, 'elon', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'elon@example.com', 'u-elon.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (6, 'lily', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'lily@example.com', 'u-lily.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (7, 'charlie', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'charlie@example.com', 'u-charlie.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (8, 'mac', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'mac@example.com', 'u-mac.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (9, 'dennis', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'dennis@example.com', 'u-dennis.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (10, 'frank', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'frank@example.com', 'u-frank.jpg', 1651906207);
INSERT INTO users (id, username, password, email, profilePicPath, createdAt) VALUES (11, 'dee', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'dee@example.com', 'u-dee.jpg', 1651906207);

-- Tweets
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (1, 1, 'First tweet ever', null, 1631906207);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (2, 3, 'Me when it''s morphin'' time', 't-00.gif', 1641906207);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (3, 2, '3rd tweet ever', null, 1646906207);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (4, 1, 'my stomach when I run a query on the prod DB: ⬆️⬇️⬆️⬇️🔄', null, 1649406207);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (5, 3, 'Hittin gym or whatever', 't-01.jpg', 1650656207);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (6, 3, '> can you see my screen?
me: yes, for the last two years, yes, when you press that "share screen" button, that shares your screen and we can see it', null, 1651281207);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (7, 5, 'Twitter algorithm should be open source', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (8, 4, 'Investigate what happened on 3/11 https://www.youtube.com/watch?v=dQw4w9WgXcQ', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (9, 6, 'One and only tweet', 't-02.jpg', 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (10, 5, 'I dunno … seems kinda fungible', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (11, 10, 'I''m not gonna be buried in a grave. When I''m dead, just throw me in the trash.', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (12, 8, 'I''m not fat. I''m cultivating mass.', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (13, 7, 'His Neck Is High. Makes Me Trust Him.', 't-03.jpg', 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (14, 8, 'Well, maybe it boils down to this, smart guy: Computers are for losers.', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (15, 9, 'I''m To Remember Every Man I''ve Seen Fall Into A Plate Of Spaghetti? Get Real.', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (16, 8, 'And you know what happens with Tokyo drifting? It leads to bickering. Which, of course, leads to karate', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (17, 10, 'Well, I don''t know how many years on this Earth I got left. I''m gonna get real weird with it.', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (18, 11, 'I''m not asking you to do much. Just turn a blind eye while I rob this place stupid.', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (19, 8, 'I don''t appreciate being paraphrased. Now, I choose my words very deliberately.', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (20, 7, 'Oh, shit. Look at that door, dude. See that door right there? The one marked "pirate?" You think a pirate lives in there?', null, 1651593707);
INSERT INTO tweets (id, userId, text, mediaPath, createdAt) VALUES (21, 7, 'Can I stop you, though? You keep using this word "jabroni." And…it''s awesome!', null, 1651593707);

-- Retweets
INSERT INTO retweets (userId, tweetId) VALUES (3, 1);
INSERT INTO retweets (userId, tweetId) VALUES (4, 1);
INSERT INTO retweets (userId, tweetId) VALUES (3, 7);
INSERT INTO retweets (userId, tweetId) VALUES (3, 3);
INSERT INTO retweets (userId, tweetId) VALUES (1, 6);
INSERT INTO retweets (userId, tweetId) VALUES (6, 2);
INSERT INTO retweets (userId, tweetId) VALUES (6, 4);
INSERT INTO retweets (userId, tweetId) VALUES (6, 1);
INSERT INTO retweets (userId, tweetId) VALUES (10, 13);
INSERT INTO retweets (userId, tweetId) VALUES (7, 14);
INSERT INTO retweets (userId, tweetId) VALUES (8, 14);
INSERT INTO retweets (userId, tweetId) VALUES (11, 14);
INSERT INTO retweets (userId, tweetId) VALUES (7, 17);
INSERT INTO retweets (userId, tweetId) VALUES (9, 21);

-- Reactions
INSERT INTO reactions (userId, tweetId, reaction) VALUES (2, 1, 'shocked');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (3, 1, 'party');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (4, 1, 'party');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (10, 1, 'party');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (10, 12, 'laugh');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (10, 15, 'raisedHands');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (10, 21, 'heart');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (10, 7, 'heart');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (10, 6, 'laugh');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (3, 12, 'laugh');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (4, 12, 'laugh');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (5, 12, 'laugh');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (6, 12, 'laugh');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (7, 18, 'sad');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (8, 18, 'sad');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (1, 20, 'laugh');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (2, 20, 'shocked');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (9, 20, 'shocked');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (7, 19, 'shocked');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (7, 10, 'laugh');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (7, 5, 'shocked');
INSERT INTO reactions (userId, tweetId, reaction) VALUES (7, 4, 'laugh');

-- Follows
INSERT INTO follows (userId, followedUserId) VALUES (1, 3);
INSERT INTO follows (userId, followedUserId) VALUES (2, 1);
INSERT INTO follows (userId, followedUserId) VALUES (2, 3);
INSERT INTO follows (userId, followedUserId) VALUES (3, 1);
INSERT INTO follows (userId, followedUserId) VALUES (4, 1);
INSERT INTO follows (userId, followedUserId) VALUES (7, 8);
INSERT INTO follows (userId, followedUserId) VALUES (7, 9);
INSERT INTO follows (userId, followedUserId) VALUES (7, 10);
INSERT INTO follows (userId, followedUserId) VALUES (7, 11);
INSERT INTO follows (userId, followedUserId) VALUES (8, 7);
INSERT INTO follows (userId, followedUserId) VALUES (8, 9);
INSERT INTO follows (userId, followedUserId) VALUES (8, 10);
INSERT INTO follows (userId, followedUserId) VALUES (8, 11);
INSERT INTO follows (userId, followedUserId) VALUES (9, 7);
INSERT INTO follows (userId, followedUserId) VALUES (9, 8);
INSERT INTO follows (userId, followedUserId) VALUES (9, 10);
INSERT INTO follows (userId, followedUserId) VALUES (9, 11);
INSERT INTO follows (userId, followedUserId) VALUES (10, 7);
INSERT INTO follows (userId, followedUserId) VALUES (10, 8);
INSERT INTO follows (userId, followedUserId) VALUES (10, 9);
INSERT INTO follows (userId, followedUserId) VALUES (10, 11);
INSERT INTO follows (userId, followedUserId) VALUES (11, 7);
INSERT INTO follows (userId, followedUserId) VALUES (11, 8);
INSERT INTO follows (userId, followedUserId) VALUES (11, 9);
INSERT INTO follows (userId, followedUserId) VALUES (11, 10);