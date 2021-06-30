package orderbook

import (
	"container/list"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

type OrderList struct {
	qty decimal.Decimal
	price  decimal.Decimal
	orders *list.List
}

func NewOrderList(price decimal.Decimal) *OrderList{
	return &OrderList{
		qty: decimal.Zero,
		price: price,
		orders: list.New(),
	}
}

func (ol *OrderList) Len() int {
	return ol.orders.Len()
}

func (ol *OrderList) Qty() decimal.Decimal {
	return ol.qty
}

func (ol *OrderList) Price() decimal.Decimal {
	return ol.price
}

func (ol *OrderList) Head() *list.Element {
	return ol.orders.Front()
}

func (ol *OrderList) Tail() *list.Element {
	return ol.orders.Back()
}

func (ol *OrderList) Push(o Order) *list.Element {
	ol.addQty(o.GetQty())
	return ol.orders.PushBack(o)
}

func (ol *OrderList) Remove(e *list.Element) Order {
	ol.subQty(e.Value.(Order).GetQty())
	return ol.orders.Remove(e).(Order)
}

func (ol *OrderList) addQty(qty decimal.Decimal)  {
	ol.qty = ol.qty.Add(qty)
}

func (ol *OrderList) subQty(qty decimal.Decimal)  {
	ol.qty = ol.qty.Sub(qty)
}

func (ol *OrderList) String() string {
	sb := strings.Builder{}
	iter := ol.orders.Front()
	sb.WriteString(fmt.Sprintf("price: %s, qty: %s, order size: %d, orderdetail:\n", ol.Price(), ol.Qty(), ol.Len()))
	for iter != nil {
		order := iter.Value.(Order)
		sb.WriteString(order.String() + "\n")
		iter = iter.Next()
	}
	return sb.String()
}