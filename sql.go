package hitbtc

const (
	//CreateTableTicker - query to create table in database for HitBTC ticker
	CreateTableTicker = "CREATE TABLE IF NOT EXISTS hitbtc_ticker (TickerID VARCHAR(20) NOT NULL PRIMARY KEY,Symbol VARCHAR(20),Ask DECIMAL,Bid DECIMAL,`Last` DECIMAL,`Open` DECIMAL,Low DECIMAL,High DECIMAL,Volume DECIMAL,VolumeQuote DECIMAL,`Timestamp` TIMESTAMP,);"
)
