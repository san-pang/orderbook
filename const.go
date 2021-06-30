package orderbook

/*
	买卖方向，订单簿只设置买入和卖出的委托方向，其它的委托方向例如申购、赎回，需要自己转换成买入和卖出
 */
type OrderSide string
const SIDE_BUY OrderSide = "1"  // 买入
const SIDE_SELL OrderSide = "2"  // 卖出

