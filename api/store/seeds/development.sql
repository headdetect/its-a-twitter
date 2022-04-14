-- admin:password
insert into users (username, password, email, createdAt) values (
  'admin', '$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW', 'admin@example.com', strftime('%s', 'now')
);