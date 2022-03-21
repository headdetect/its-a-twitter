create table if not exists users (
  id integer not null primary key, 
  username text,
  displayName text,
  profilePicHash text,
  password text,
  createdAt integer
);

insert into users(id, username, displayName, profilePicHash, password, createdAt) values(
  1, "admin", "admin", "admin", "$2a$14$Xyy1t.JsY.LHkN23Rsw5meOcmrVAueiRFEE8m5.YJy.vUi0.qowcW", 1647840227
)