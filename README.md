# cryptomarketdata

Data is inserted into TimescaleDB - a PostgreSQL extension.

Tables needed:

CREATE TABLE btc_dayohlc (
  time        TIMESTAMPTZ       NOT NULL,
  open DOUBLE PRECISION  NOT NULL,
  high DOUBLE PRECISION  NOT NULL,
  low    DOUBLE PRECISION  NOT NULL,
  close DOUBLE PRECISION  NOT NULL
);

CREATE TABLE eth_dayohlc (
  time        TIMESTAMPTZ       NOT NULL,
  open DOUBLE PRECISION  NOT NULL,
  high DOUBLE PRECISION  NOT NULL,
  low    DOUBLE PRECISION  NOT NULL,
  close DOUBLE PRECISION  NOT NULL
);
