package accrual

const (
	getUnprocessedOrdersQuery string = "SELECT order_number, status FROM orders WHERE status IN ($1, $2)"
	setProcessingStatusQuery  string = "UPDATE orders SET status=$1 WHERE order_number=$2"
	setProcessedStatusQuery   string = "UPDATE orders SET status=$1, accrual =$2 WHERE order_number=$3 RETURNING user_id"
	addPointsToUser           string = "UPDATE users SET current=users.current+$1 WHERE id=$2"
)
