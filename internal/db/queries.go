package db

const (
	createUserQuery             string = `INSERT INTO %s (login, password_hash, api_key) VALUES ($1, $2, $3) RETURNING id`
	getUserQuery                string = `SELECT password_hash, api_key FROM %s WHERE login=$1`
	getUserByKeyQuery           string = `SELECT id FROM %s where api_key=$1`
	addOrderQuery               string = `INSERT INTO %s (order_number, user_id, status) VALUES ($1, $2, $3)`
	getUserIDByOrderNumberQuery string = `SELECT user_id FROM %s WHERE order_number=$1`
	getUsersOrdersQuery         string = `SELECT order_number, status, accrual, created_at FROM %s WHERE user_id=$1 ORDER BY created_at DESC`
	getBalanceQuery             string = `SELECT withdrawn, current FROM %s WHERE id=$1`
	updateBalanceQuery          string = `UPDATE %s SET current=$1, withdrawn=$2 WHERE id=$3`
	addWithdrawRecordQuery      string = `INSERT INTO %s (user_id, order_number, sum) VALUES ($1, $2, $3) RETURNING id`
	getWithdrawsHistoryQuery    string = `SELECT order_number, sum, processed_at FROM %s WHERE user_id=$1 ORDER BY processed_at DESC`
)
