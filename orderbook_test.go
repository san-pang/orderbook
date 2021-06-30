package orderbook

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

type myOrder struct {
	ApplID string
	ClOrdID string
	Price decimal.Decimal
	Qty decimal.Decimal
	Side OrderSide
}

func newMyOrder(applID, clordID string, price, qty decimal.Decimal, side OrderSide) *myOrder {
	return &myOrder{
		ApplID:  applID,
		ClOrdID: clordID,
		Price:   price,
		Qty:     qty,
		Side: side,
	}
}

func (o *myOrder)GetApplID() string {
	return o.ApplID
}

func (o *myOrder)GetClOrdID() string {
	return o.ClOrdID
}

func (o *myOrder)GetSide() OrderSide {
	return o.Side
}

func (o *myOrder)GetPrice() decimal.Decimal {
	return o.Price
}

func (o *myOrder)GetQty() decimal.Decimal {
	return o.Qty
}

func (o *myOrder)String() string {
	return fmt.Sprintf("applID=%s, clordID=%s, price=%s, qty=%s, side=%s", o.ApplID, o.ClOrdID, o.Price, o.Qty, o.Side)
}

func TestOrderBook_AddOrder(t *testing.T) {
	order1 := newMyOrder("010", "ord1", decimal.NewFromFloat(12.34), decimal.NewFromInt(300), SIDE_BUY)
	order2 := newMyOrder("020", "ord2", decimal.NewFromFloat(12.35), decimal.NewFromInt(400), SIDE_BUY)
	order3 := newMyOrder("030", "ord3", decimal.NewFromFloat(12.36), decimal.NewFromInt(500), SIDE_SELL)
	order4 := newMyOrder("040", "ord4", decimal.NewFromFloat(12.34), decimal.NewFromInt(600), SIDE_BUY)
	orderbook := NewOrderBook("000776")
	orderbook.AddOrder(order1.GetSide(), order1)
	orderbook.AddOrder(order2.GetSide(), order2)
	orderbook.AddOrder(order3.GetSide(), order3)
	orderbook.AddOrder(order4.GetSide(), order4)
	fmt.Println("买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("卖出方向订单簿：\n", orderbook.Sells().String())
	if orderbook.Buys().orderSize != 3 {
		t.Error("买方向订单数量错误")
	}
	if orderbook.Sells().orderSize != 1 {
		t.Error("卖方向订单数量错误")
	}
	if orderbook.Buys().depth != 2 {
		t.Error("买方向价格档位错误")
	}
	if orderbook.Sells().depth != 1 {
		t.Error("买方向价格档位错误")
	}
	if !orderbook.Buys().qty.Equals(decimal.NewFromInt(300 + 400 + 600)) {
		t.Error("买方向委托总数量错误")
	}
	if !orderbook.Sells().qty.Equals(decimal.NewFromInt(500)) {
		t.Error("卖方向委托总数量错误")
	}
}

func TestOrderBook_CancelOrder(t *testing.T) {
	order1 := newMyOrder("010", "ord1", decimal.NewFromFloat(12.34), decimal.NewFromInt(300), SIDE_BUY)
	order2 := newMyOrder("020", "ord2", decimal.NewFromFloat(12.35), decimal.NewFromInt(400), SIDE_BUY)
	order3 := newMyOrder("030", "ord3", decimal.NewFromFloat(12.36), decimal.NewFromInt(500), SIDE_SELL)
	order4 := newMyOrder("040", "ord4", decimal.NewFromFloat(12.34), decimal.NewFromInt(600), SIDE_BUY)
	order5 := newMyOrder("050", "ord5", decimal.NewFromFloat(12.35), decimal.NewFromInt(700), SIDE_SELL)
	order6 := newMyOrder("060", "ord6", decimal.NewFromFloat(12.36), decimal.NewFromInt(800), SIDE_SELL)
	orderbook := NewOrderBook("000776")
	orderbook.AddOrder(order1.GetSide(), order1)
	orderbook.AddOrder(order2.GetSide(), order2)
	orderbook.AddOrder(order3.GetSide(), order3)
	orderbook.AddOrder(order4.GetSide(), order4)
	orderbook.AddOrder(order5.GetSide(), order5)
	orderbook.AddOrder(order6.GetSide(), order6)
	fmt.Println("买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("卖出方向订单簿：\n", orderbook.Sells().String())
	if err := orderbook.CancelOrder(SIDE_BUY, "NotExistOrder"); err != ERR_ORDER_NOT_EXISTS {
		t.Error("取消不存在的订单应该报订单不存在")
	}
	if err := orderbook.CancelOrder(SIDE_SELL, "NotExistOrder"); err != ERR_ORDER_NOT_EXISTS {
		t.Error("取消不存在的订单应该报订单不存在")
	}
	if err := orderbook.CancelOrder(order2.GetSide(), order2.GetClOrdID()); err != nil {
		t.Error("取消存在的订单error应该等于nil, 实际上=", err)
	}
	if err := orderbook.CancelOrder(order6.GetSide(), order6.GetClOrdID()); err != nil {
		t.Error("取消存在的订单error应该等于nil, 实际上=", err)
	}
	if err := orderbook.CancelOrder(order5.GetSide(), order5.GetClOrdID()); err != nil {
		t.Error("取消存在的订单error应该等于nil, 实际上=", err)
	}
	if orderbook.Buys().orderSize != 2 {
		t.Error("买方向订单数量错误")
	}
	if orderbook.Sells().orderSize != 1 {
		t.Error("卖方向订单数量错误")
	}
	if orderbook.Buys().depth != 1 {
		t.Error("买方向价格档位错误")
	}
	if orderbook.Sells().depth != 1 {
		t.Error("卖方向价格档位错误")
	}
	if !orderbook.Buys().qty.Equals(decimal.NewFromInt(300 + 600)) {
		t.Error("买方向委托总数量错误")
	}
	if !orderbook.Sells().qty.Equals(decimal.NewFromInt(500)) {
		t.Error("卖方向委托总数量错误")
	}
	fmt.Println("撤单后的买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("撤单后的卖出方向订单簿：\n", orderbook.Sells().String())
}

func TestOrderBook_UpdateOrder(t *testing.T) {
	order1 := newMyOrder("010", "ord1", decimal.NewFromFloat(12.34), decimal.NewFromInt(300), SIDE_BUY)
	order2 := newMyOrder("020", "ord2", decimal.NewFromFloat(12.35), decimal.NewFromInt(400), SIDE_BUY)
	order3 := newMyOrder("030", "ord3", decimal.NewFromFloat(12.36), decimal.NewFromInt(500), SIDE_SELL)
	order4 := newMyOrder("040", "ord4", decimal.NewFromFloat(12.34), decimal.NewFromInt(600), SIDE_BUY)
	order5 := newMyOrder("050", "ord5", decimal.NewFromFloat(12.35), decimal.NewFromInt(700), SIDE_SELL)
	order6 := newMyOrder("060", "ord6", decimal.NewFromFloat(12.37), decimal.NewFromInt(800), SIDE_SELL)
	orderbook := NewOrderBook("000776")
	orderbook.AddOrder(order1.GetSide(), order1)
	orderbook.AddOrder(order2.GetSide(), order2)
	orderbook.AddOrder(order3.GetSide(), order3)
	orderbook.AddOrder(order4.GetSide(), order4)
	orderbook.AddOrder(order5.GetSide(), order5)
	orderbook.AddOrder(order6.GetSide(), order6)
	fmt.Println("买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("卖出方向订单簿：\n", orderbook.Sells().String())
	if err := orderbook.UpdateOrder(SIDE_BUY, newMyOrder("", "", decimal.NewFromFloat(11.22), decimal.NewFromInt(200), SIDE_BUY), newMyOrder("", "", decimal.NewFromFloat(11.22), decimal.NewFromInt(200), SIDE_BUY)); err != ERR_ORDER_NOT_EXISTS {
		t.Error("更新不存在的订单应该报订单不存在")
	}
	if err := orderbook.UpdateOrder(SIDE_SELL, newMyOrder("", "", decimal.NewFromFloat(11.22), decimal.NewFromInt(200), SIDE_SELL), newMyOrder("", "", decimal.NewFromFloat(11.22), decimal.NewFromInt(200), SIDE_SELL)); err != ERR_ORDER_NOT_EXISTS {
		t.Error("更新不存在的订单应该报订单不存在")
	}
	ord := *order1
	ord.ClOrdID = "ord11"
	if err := orderbook.UpdateOrder(ord.GetSide(), order1, &ord); err != ERR_CANNOT_UPDATE_ORDER {
		t.Error("更新订单时，如果新旧订单的编号不一致，error与预期不符")
	}

	ord2 := *order1
	ord2.Price = decimal.NewFromFloat(88.11)
	if err := orderbook.UpdateOrder(ord2.GetSide(), order1, &ord2); err != ERR_CANNOT_UPDATE_ORDER {
		t.Error("更新订单时，如果新旧订单的价格不一致，error与预期不符")
	}

	ord3 := *order1
	ord3.Qty = decimal.NewFromInt(1100)
	if err := orderbook.UpdateOrder(ord3.GetSide(), order1, &ord3); err != nil {
		t.Error("更新正常的订单error应该等于nil, 实际上=", err)
	}

	if orderbook.Buys().orderSize != 3 {
		t.Error("买方向订单数量错误")
	}
	if orderbook.Sells().orderSize != 3 {
		t.Error("卖方向订单数量错误")
	}
	if orderbook.Buys().depth != 2 {
		t.Error("买方向价格档位错误")
	}
	if orderbook.Sells().depth != 3 {
		t.Error("卖方向价格档位错误")
	}
	if !orderbook.Buys().qty.Equals(decimal.NewFromInt(1100 + 400 + 600)) {
		t.Error("买方向委托总数量错误")
	}
	if !orderbook.Sells().qty.Equals(decimal.NewFromInt(500 + 700 + 800)) {
		t.Error("卖方向委托总数量错误")
	}
	fmt.Println("撤单后的买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("撤单后的卖出方向订单簿：\n", orderbook.Sells().String())
}

func TestOrderBook_LessThan(t *testing.T) {
	order1 := newMyOrder("010", "ord1", decimal.NewFromFloat(12.34), decimal.NewFromInt(300), SIDE_BUY)
	order2 := newMyOrder("020", "ord2", decimal.NewFromFloat(12.35), decimal.NewFromInt(400), SIDE_BUY)
	order4 := newMyOrder("040", "ord4", decimal.NewFromFloat(12.34), decimal.NewFromInt(600), SIDE_BUY)

	order3 := newMyOrder("030", "ord3", decimal.NewFromFloat(12.36), decimal.NewFromInt(500), SIDE_SELL)
	order5 := newMyOrder("050", "ord5", decimal.NewFromFloat(12.35), decimal.NewFromInt(700), SIDE_SELL)
	order6 := newMyOrder("060", "ord6", decimal.NewFromFloat(12.37), decimal.NewFromInt(800), SIDE_SELL)
	orderbook := NewOrderBook("000776")
	orderbook.AddOrder(order1.GetSide(), order1)
	orderbook.AddOrder(order2.GetSide(), order2)
	orderbook.AddOrder(order3.GetSide(), order3)
	orderbook.AddOrder(order4.GetSide(), order4)
	orderbook.AddOrder(order5.GetSide(), order5)
	orderbook.AddOrder(order6.GetSide(), order6)
	fmt.Println("买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("卖出方向订单簿：\n", orderbook.Sells().String())

	fmt.Println("价格小于12.35的买入订单：")
	order_buy := orderbook.Buys().LessThan(decimal.NewFromFloat(12.35))
	for order_buy != nil {
		fmt.Println(order_buy)
		order_buy = orderbook.Buys().LessThan(order_buy.Price())
	}

	fmt.Println("价格小于12.37的卖出订单：")
	order_sell := orderbook.Sells().LessThan(decimal.NewFromFloat(12.37))
	for order_sell != nil {
		fmt.Println(order_sell)
		order_sell = orderbook.Sells().LessThan(order_sell.Price())
	}
}

func TestOrderBook_LessThanEqual(t *testing.T) {
	order1 := newMyOrder("010", "ord1", decimal.NewFromFloat(12.34), decimal.NewFromInt(300), SIDE_BUY)
	order2 := newMyOrder("020", "ord2", decimal.NewFromFloat(12.35), decimal.NewFromInt(400), SIDE_BUY)
	order4 := newMyOrder("040", "ord4", decimal.NewFromFloat(12.34), decimal.NewFromInt(600), SIDE_BUY)

	order3 := newMyOrder("030", "ord3", decimal.NewFromFloat(12.36), decimal.NewFromInt(500), SIDE_SELL)
	order5 := newMyOrder("050", "ord5", decimal.NewFromFloat(12.35), decimal.NewFromInt(700), SIDE_SELL)
	order6 := newMyOrder("060", "ord6", decimal.NewFromFloat(12.37), decimal.NewFromInt(800), SIDE_SELL)
	orderbook := NewOrderBook("000776")
	orderbook.AddOrder(order1.GetSide(), order1)
	orderbook.AddOrder(order2.GetSide(), order2)
	orderbook.AddOrder(order3.GetSide(), order3)
	orderbook.AddOrder(order4.GetSide(), order4)
	orderbook.AddOrder(order5.GetSide(), order5)
	orderbook.AddOrder(order6.GetSide(), order6)
	fmt.Println("买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("卖出方向订单簿：\n", orderbook.Sells().String())

	fmt.Println("价格小于等于12.34的买入订单：")
	order_buy := orderbook.Buys().LessEqualThan(decimal.NewFromFloat(12.34))
	for order_buy != nil {
		fmt.Println(order_buy)
		order_buy = orderbook.Buys().LessThan(order_buy.Price())
	}

	fmt.Println("价格小于等于12.36的卖出订单：")
	order_sell := orderbook.Sells().LessEqualThan(decimal.NewFromFloat(12.36))
	for order_sell != nil {
		fmt.Println(order_sell)
		order_sell = orderbook.Sells().LessThan(order_sell.Price())
	}
}

func TestOrderBook_GreaterThan(t *testing.T) {
	order1 := newMyOrder("010", "ord1", decimal.NewFromFloat(12.34), decimal.NewFromInt(300), SIDE_BUY)
	order2 := newMyOrder("020", "ord2", decimal.NewFromFloat(12.35), decimal.NewFromInt(400), SIDE_BUY)
	order4 := newMyOrder("040", "ord4", decimal.NewFromFloat(12.34), decimal.NewFromInt(600), SIDE_BUY)

	order3 := newMyOrder("030", "ord3", decimal.NewFromFloat(12.36), decimal.NewFromInt(500), SIDE_SELL)
	order5 := newMyOrder("050", "ord5", decimal.NewFromFloat(12.35), decimal.NewFromInt(700), SIDE_SELL)
	order6 := newMyOrder("060", "ord6", decimal.NewFromFloat(12.37), decimal.NewFromInt(800), SIDE_SELL)
	orderbook := NewOrderBook("000776")
	orderbook.AddOrder(order1.GetSide(), order1)
	orderbook.AddOrder(order2.GetSide(), order2)
	orderbook.AddOrder(order3.GetSide(), order3)
	orderbook.AddOrder(order4.GetSide(), order4)
	orderbook.AddOrder(order5.GetSide(), order5)
	orderbook.AddOrder(order6.GetSide(), order6)
	fmt.Println("买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("卖出方向订单簿：\n", orderbook.Sells().String())

	fmt.Println("价格大于12.34的买入订单：")
	order_buy := orderbook.Buys().GreaterThan(decimal.NewFromFloat(12.34))
	for order_buy != nil {
		fmt.Println(order_buy)
		order_buy = orderbook.Buys().GreaterThan(order_buy.Price())
	}

	fmt.Println("价格大于12.35的卖出订单：")
	order_sell := orderbook.Sells().GreaterThan(decimal.NewFromFloat(12.35))
	for order_sell != nil {
		fmt.Println(order_sell)
		order_sell = orderbook.Sells().GreaterThan(order_sell.Price())
	}
}

func TestOrderBook_GreaterEqualThan(t *testing.T) {
	order1 := newMyOrder("010", "ord1", decimal.NewFromFloat(12.34), decimal.NewFromInt(300), SIDE_BUY)
	order2 := newMyOrder("020", "ord2", decimal.NewFromFloat(12.35), decimal.NewFromInt(400), SIDE_BUY)
	order4 := newMyOrder("040", "ord4", decimal.NewFromFloat(12.34), decimal.NewFromInt(600), SIDE_BUY)

	order3 := newMyOrder("030", "ord3", decimal.NewFromFloat(12.36), decimal.NewFromInt(500), SIDE_SELL)
	order5 := newMyOrder("050", "ord5", decimal.NewFromFloat(12.35), decimal.NewFromInt(700), SIDE_SELL)
	order6 := newMyOrder("060", "ord6", decimal.NewFromFloat(12.37), decimal.NewFromInt(800), SIDE_SELL)
	orderbook := NewOrderBook("000776")
	orderbook.AddOrder(order1.GetSide(), order1)
	orderbook.AddOrder(order2.GetSide(), order2)
	orderbook.AddOrder(order3.GetSide(), order3)
	orderbook.AddOrder(order4.GetSide(), order4)
	orderbook.AddOrder(order5.GetSide(), order5)
	orderbook.AddOrder(order6.GetSide(), order6)
	fmt.Println("买入方向订单簿：\n", orderbook.Buys().String())
	fmt.Println("卖出方向订单簿：\n", orderbook.Sells().String())

	fmt.Println("价格大于等于12.34的买入订单：")
	order_buy := orderbook.Buys().GreaterEqualThan(decimal.NewFromFloat(12.34))
	for order_buy != nil {
		fmt.Println(order_buy)
		order_buy = orderbook.Buys().GreaterThan(order_buy.Price())
	}

	fmt.Println("价格大于等于12.36的卖出订单：")
	order_sell := orderbook.Sells().GreaterEqualThan(decimal.NewFromFloat(12.36))
	for order_sell != nil {
		fmt.Println(order_sell)
		order_sell = orderbook.Sells().GreaterThan(order_sell.Price())
	}
}