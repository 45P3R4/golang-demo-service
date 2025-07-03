CREATE TABLE payments(
   transaction   VARCHAR(20) NOT NULL PRIMARY KEY
  ,request_id    VARCHAR(30)
  ,currency      VARCHAR(3)
  ,provider      VARCHAR(10)
  ,amount        INTEGER  NOT NULL
  ,payment_dt    INTEGER
  ,bank          VARCHAR(10)
  ,delivery_cost INTEGER	  NOT NULL
  ,goods_total   INTEGER  NOT NULL
  ,custom_fee    INTEGER  NOT NULL
);

CREATE TABLE items(
  chrt_id      INTEGER
  ,track_number VARCHAR(20) NOT NULL
  ,price        INTEGER  NOT NULL
  ,rid          VARCHAR(21) PRIMARY KEY
  ,name         VARCHAR(30) NOT NULL
  ,sale         INTEGER  NOT NULL
  ,size         INTEGER  NOT NULL
  ,total_price  INTEGER  NOT NULL
  ,nm_id        INTEGER  NOT NULL
  ,brand        VARCHAR(20)
  ,status       INTEGER  NOT NULL
);

CREATE TABLE deliveries(
   delivery_uid VARCHAR(20) NOT NULL PRIMARY KEY
  ,name        VARCHAR(20) NOT NULL
  ,phone       VARCHAR(11)  NOT NULL
  ,zip         INTEGER  NOT NULL
  ,city        VARCHAR(20) NOT NULL
  ,address     VARCHAR(20) NOT NULL
  ,region      VARCHAR(20) NOT NULL
  ,email       VARCHAR(20) NOT NULL
);

CREATE TABLE orders(
   order_uid          VARCHAR(20) NOT NULL PRIMARY KEY
  ,track_number       VARCHAR(20) NOT NULL
  ,entry              VARCHAR(5) NOT NULL
  ,delivery_uid           VARCHAR(20) NOT NULL
  ,payment_transaction    VARCHAR(20) NOT NULL
  ,items_rid              VARCHAR(21)[] NOT NULL
  ,locale             VARCHAR(3) NOT NULL
  ,internal_signature VARCHAR(30)
  ,customer_id        VARCHAR(20) NOT NULL
  ,delivery_service   VARCHAR(10) NOT NULL
  ,shardkey           INTEGER  NOT NULL
  ,sm_id              INTEGER  NOT NULL
  ,date_created       TIME  NOT NULL
  ,oof_shard          INTEGER  NOT NULL
);