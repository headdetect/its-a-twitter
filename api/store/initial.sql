create table if not exists users (
  id integer not null primary key, 
  username text,
  displayName text,
  password text,
  createdAt integer
);

insert into users(id, username, displayName, password, createdAt) values(
  1, "admin", "admin", "$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW", 1647840227
)

insert into users(id, username, displayName, password, createdAt) values(
  2, "test", "test", "$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW", 1647840228
)