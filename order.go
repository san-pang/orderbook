package orderbook

import (
	"github.com/shopspring/decimal"
)

type Order interface {
	GetApplID() string  //业务标识
	GetClOrdID() string //订单号
	GetPrice() decimal.Decimal  //委托价格
	GetQty() decimal.Decimal  //委托数量
	String() string // 订单的字符串打印数据
}