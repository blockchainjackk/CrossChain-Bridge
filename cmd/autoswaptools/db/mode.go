package db

import (
	"database/sql"
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db/ddl"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db/types"
	"github.com/anyswap/CrossChain-Bridge/log"
)

// closeRows closes the input sql.Rows, logging any error.
func closeRows(rows *sql.Rows) {
	if e := rows.Close(); e != nil {
		log.Fatalf("Close of Query failed: %v", e)
	}
}

// SqlExecutor is implemented by both sql.DB and sql.Tx.
type SqlExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// SqlQueryer is implemented by both sql.DB and sql.Tx.
type SqlQueryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// SqlExecQueryer is implemented by both sql.DB and sql.Tx.
type SqlExecQueryer interface {
	SqlExecutor
	SqlQueryer
}

// sqlExec executes the SQL statement string with any optional arguments, and
// returns the number of rows affected.
func sqlExec(db SqlExecutor, stmt, execErrPrefix string, args ...interface{}) (int64, error) {
	res, err := db.Exec(stmt, args...)
	if err != nil {
		return 0, fmt.Errorf("%v: %w", execErrPrefix, err)
	}
	if res == nil {
		return 0, nil
	}

	var N int64
	N, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error in RowsAffected: %w", err)
	}
	return N, err
}

// sqlExecStmt executes the prepared SQL statement with any optional arguments,
// and returns the number of rows affected.
func sqlExecStmt(stmt *sql.Stmt, execErrPrefix string, args ...interface{}) (int64, error) {
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("%v: %w", execErrPrefix, err)
	}
	if res == nil {
		return 0, nil
	}

	var N int64
	N, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error in RowsAffected: %w", err)
	}
	return N, err
}

// TableExists checks if the specified table exists.
func (db *CrossChainDB) TableExists(tableName string) (bool, error) {
	rows, err := db.db.Query(`select relname from pg_class where relname = $1`,
		tableName)
	if err != nil {
		return false, err
	}

	defer func() {
		if e := rows.Close(); e != nil {
			log.Fatalf("Close of Query failed: %v", e)
		}
	}()
	return rows.Next(), nil
}

// CreateTable creates a table with the given name using the provided SQL
// statement, if it does not already exist.
func (db *CrossChainDB) CreateTable(tableName, stmt string) error {
	exists, err := db.TableExists(tableName)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf(`Creating the "%s" table.`, tableName)
		_, err = db.db.Exec(stmt)
		if err != nil {
			return err
		}
	} else {
		//log.Tracef(`Table "%s" exists.`, tableName)
		log.Printf(`Table "%s" exists.`, tableName)
	}

	return err
}

func (db *CrossChainDB) RetrieveAddressCount() (int64, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectAddressCountRow)
	if err != nil {
		return 0, err
	}

	defer closeRows(rows)
	rows.Next()
	var count int64
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return count, nil
}
func (db *CrossChainDB) InsertAddress(key string, address string, balance int64) error {
	stmt, err := db.db.Prepare(ddl.InsertAddressRow)
	if err != nil {
		log.Fatalf("insert address err: %v", err)
		return err
	}
	_, err = stmt.Exec(key, address, balance)

	return err
}

func (db *CrossChainDB) RetrieveAddressesToSwapIn(value int64) ([]*types.AddressInfo, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectAddressToSwapInRow, value)
	if err != nil {
		return nil, err
	}

	defer closeRows(rows)
	addresses := make([]*types.AddressInfo, 0)
	for rows.Next() {
		var key string
		var address string
		var balance int64
		err = rows.Scan(&key, &address, &balance)
		if err != nil {
			return nil, err
		}
		a := &types.AddressInfo{
			Key:     key,
			Address: address,
			Balance: balance,
		}
		addresses = append(addresses, a)
	}
	err = rows.Err()

	return addresses, err
}
func (db *CrossChainDB) RetrieveTxsToSwapOut() ([]string, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectTxToSwapOutRow)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)
	txs := make([]string, 0)
	for rows.Next() {
		var tx string
		err = rows.Scan(&tx)
		if err != nil {
			return nil, err
		}

		txs = append(txs, tx)
	}
	err = rows.Err()
	return txs, err
}
func (db *CrossChainDB) RetrieveAddressFromSwapIn2() ([]string, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectToAddressFromSwapIn2)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)
	addrs := make([]string, 0)
	for rows.Next() {
		var addr string
		err = rows.Scan(&addr)
		if err != nil {
			return nil, err
		}

		addrs = append(addrs, addr)
	}
	err = rows.Err()
	return addrs, err
}

func (db *CrossChainDB) RetrieveKeyByAddress(address string) (string, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectKeyByAddress, address)
	if err != nil {
		return "", err
	}
	defer closeRows(rows)
	var key string
	rows.Next()

	err = rows.Scan(&key)
	if err != nil {
		return "", err
	}
	err = rows.Err()

	return key, err
}
func (db *CrossChainDB) RetrieveBindByAddress(address string) (string, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectBindByToAddress, address)
	if err != nil {
		return "", err
	}
	defer closeRows(rows)
	var bind string
	rows.Next()

	err = rows.Scan(&bind)
	if err != nil {
		return "", err
	}
	err = rows.Err()

	return bind, err
}
func (db *CrossChainDB) UpdateAddrBalance(balance int64, addr string) error {
	stmt, err := db.db.Prepare(ddl.UpdateAddrBalanceRow)
	if err != nil {
		log.Fatalf("update address balance err: %v", err)
		return err
	}
	_, err = stmt.Exec(balance, addr)

	return err

}
func (db *CrossChainDB) InsertTxInSwapIn(txId string, status int64) error {
	stmt, err := db.db.Prepare(ddl.InsertTxInSwapIn)
	if err != nil {
		log.Fatalf("insert address err: %v", err)
		return err
	}
	_, err = stmt.Exec(txId, status)

	return err
}
func (db *CrossChainDB) UpdateSwapIn(
	txId string, fromAddress string, toAddress string, singInfo string, status int64) error {
	stmt, err := db.db.Prepare(ddl.UpdateSwapInRow)
	if err != nil {
		log.Fatalf("insert address err: %v", err)
		return err
	}
	_, err = stmt.Exec(txId, fromAddress, toAddress, singInfo, status)

	return err

}
func (db *CrossChainDB) RetrieveTx2SwapIn() ([]string, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectTxToSwapInRow)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)
	txIds := make([]string, 0)
	for rows.Next() {
		var txId string

		err = rows.Scan(&txId)
		if err != nil {
			return nil, err
		}

		txIds = append(txIds, txId)
	}
	err = rows.Err()

	return txIds, err
}

func (db *CrossChainDB) RetrieveToAddressFromSwapIn() ([]string, []string, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectToAddressFromSwapIn)
	if err != nil {
		return nil, nil, err
	}
	defer closeRows(rows)
	fromAddrs := make([]string, 0)
	toAddrs := make([]string, 0)
	for rows.Next() {
		var from string
		var to string
		err = rows.Scan(&from, &to)
		if err != nil {
			return nil, nil, err
		}

		fromAddrs = append(fromAddrs, from)
		toAddrs = append(toAddrs, to)
	}
	err = rows.Err()

	return fromAddrs, toAddrs, err
}
func (db *CrossChainDB) RetrieveAddressToSwapOut(balance int64) ([]*types.AddressInfo, error) {
	var rows *sql.Rows
	rows, err := db.db.Query(ddl.SelectAddressToSwapOutRow, balance, 10)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)
	addrInfo := make([]*types.AddressInfo, 0)

	for rows.Next() {
		var key string
		var address string
		var balance int64
		err = rows.Scan(&key, &address, &balance)
		if err != nil {
			return nil, err
		}
		info := &types.AddressInfo{
			Key:     key,
			Address: address,
			Balance: balance,
		}

		addrInfo = append(addrInfo, info)
	}
	err = rows.Err()

	return addrInfo, err
}

func (db *CrossChainDB) InsertTxInSwapOut(txId string, from string, bind string, status int64) error {
	stmt, err := db.db.Prepare(ddl.InsertTxOutSwapRow)
	if err != nil {
		log.Fatalf("insert address err: %v", err)
		return err
	}
	_, err = stmt.Exec(txId, from, bind, status)

	return err
}
func (db *CrossChainDB) UpdateSwapOutStatus(txId string, status int64) error {
	stmt, err := db.db.Prepare(ddl.UpdateSwapOutStatusRow)
	if err != nil {
		log.Fatalf("insert address err: %v", err)
		return err
	}
	_, err = stmt.Exec(txId, status)

	return err

}
