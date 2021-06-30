package orderbook

import "errors"

var (
	ERR_ORDER_NOT_EXISTS = errors.New("order not found")
	ERR_SIDE_INVALID = errors.New("side is invalid")
	ERR_DULPLICATE_ORDER = errors.New("order already exists")
	ERR_CANNOT_UPDATE_ORDER = errors.New("price or orderid is not same, not allow to update")
	ERR_PRICELEVEL_NOTFOUND = errors.New("price level not found")
)
