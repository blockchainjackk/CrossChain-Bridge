package ddl

const (
	BscAddressesTableName = "bscaddress"

	CreateAddressTable = `CREATE TABLE IF NOT EXISTS bscaddress (
		id SERIAL8 PRIMARY KEY,
       	key TEXT ,
		address TEXT ,
		balance INT8 

	);`

	SelectAddressCountRow = `SELECT count(address) from bscaddress;`

	SelectAddressToSwapInRow = `SELECT key,address,balance from bscaddress  ORDER BY balance  LIMIT $1;`

	SelectKeyByAddress = `SELECT key from bscaddress where address = $1;`

	InsertAddressRow = `INSERT INTO bscaddress (key,address,balance) VALUES ($1, $2, $3);`

	SelectAddressToSwapOutRow = `SELECT  key,address,balance  from  bscaddress  where balance >= $1 LIMIT $2;`
	UpdateAddrBalanceRow      = `UPDATE bscaddress SET balance = $1  WHERE address = $2;`

	SwapInTaleName    = "swapin"
	CreateSwapInTable = `CREATE TABLE IF NOT EXISTS swapin (
		id SERIAL8 PRIMARY KEY,
		txid TEXT,
       	fromaddress TEXT ,
		toaddress TEXT ,
		signinfo TEXT,
		status  INT8
	);`
	InsertTxInSwapIn = `INSERT INTO swapin (txid,status) VALUES ($1,$2);`
	//InsertSwapInRow     = `INSERT INTO swapin (txid,fromaddressstatus,toaddress,signinfo,status) VALUES ($1,$2,$3,$4,$5)`
	SelectTxToSwapInRow   = `SELECT txid  from  swapin  where status = 0;`
	SelectBindByToAddress = `SELECT fromaddress  from  swapin  where toaddress = $1 LIMIT 1;`

	UpdateSwapInRow           = `UPDATE swapin SET fromaddress = $2,toaddress = $3,signinfo = $4,status = $5  WHERE txid = $1;`
	SelectToAddressFromSwapIn = `SELECT fromaddress,toaddress  from  swapin  where status = 1;`

	SelectToAddressFromSwapIn2 = `SELECT DISTINCT toaddress  from  swapin  where status = 1;`

	SwapOutTaleName = "swapout"

	CreateSwapOutTable = `CREATE TABLE IF NOT EXISTS swapout (
		id SERIAL8 PRIMARY KEY,
		txid TEXT,
       	fromaddress TEXT ,
		bind TEXT ,
		status  INT8
	);`
	InsertTxOutSwapRow     = `INSERT INTO swapout (txid,fromaddress,bind,status) VALUES ($1,$2,$3,$4);`
	SelectTxToSwapOutRow   = `SELECT txid  from  swapout  where status = 0;`
	UpdateSwapOutStatusRow = `UPDATE swapout SET status = $2  WHERE txid = $1;`
)
