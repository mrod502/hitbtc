CREATE TABLE IF NOT EXISTS hitbtc_ticker (
    TickerID VARCHAR(20) NOT NULL PRIMARY KEY,
    Symbol VARCHAR(20),
    Ask DECIMAL,
    Bid DECIMAL,
    "Last" DECIMAL,
    "Open" DECIMAL,
    Low DECIMAL,
    High DECIMAL,
    Volume DECIMAL,
    VolumeQuote DECIMAL,
    "Timestamp" TIMESTAMP
);