package orderbook

import (
	"container/list"
	"fmt"
)

type OrderBook struct {
	buyOrders map[string]*list.Element
	sellOrders map[string]*list.Element
	sells *OrderTree
	buys *OrderTree
	securityID string
}

func NewOrderBook(securityID string) *OrderBook {
	return &OrderBook{
		buyOrders: make(map[string]*list.Element),
		sellOrders: make(map[string]*list.Element),
		sells:  NewOrderTree(),
		buys:   NewOrderTree(),
		securityID: securityID,
	}
}

func (ob *OrderBook) CancelOrder(side OrderSide, orderID string) error {
	if side == SIDE_BUY {
		if e, found := ob.buyOrders[orderID]; found {
			delete(ob.buyOrders, orderID)
			ob.buys.Remove(e)
			return nil
		}
		return ERR_ORDER_NOT_EXISTS
	} else if side == SIDE_SELL {
		if e, found := ob.sellOrders[orderID]; found {
			delete(ob.sellOrders, orderID)
			ob.sells.Remove(e)
			return nil
		}
		return ERR_ORDER_NOT_EXISTS
	}
	return ERR_SIDE_INVALID
}

func (ob *OrderBook) GetOrder(side OrderSide, orderID string) (*list.Element, error) {
	if side == SIDE_BUY {
		if e, found := ob.buyOrders[orderID]; found {
			return e, nil
		}
		return nil, ERR_ORDER_NOT_EXISTS
	} else if side == SIDE_SELL {
		if e, found := ob.sellOrders[orderID]; found {
			return e, nil
		}
		return nil, ERR_ORDER_NOT_EXISTS
	}
	return nil, ERR_SIDE_INVALID
}

func (ob *OrderBook) UpdateOrder(side OrderSide, oldOrder Order, newOrder Order) error {
	fmt.Println(side, oldOrder, newOrder)
	// 这里是基于订单编号和订单价格不变的前提才能更新，因为涉及到订单列表、层次、订单数量的统计更新
	if oldOrder.GetClOrdID() != newOrder.GetClOrdID() || !oldOrder.GetPrice().Equals(newOrder.GetPrice()) {
		return ERR_CANNOT_UPDATE_ORDER
	}
	if side == SIDE_BUY {
		e, found := ob.buyOrders[oldOrder.GetClOrdID()]
		if !found {
			return ERR_ORDER_NOT_EXISTS
		}
		orderlist, found := ob.buys.prices[newOrder.GetPrice().String()]
		if !found {
			return ERR_PRICELEVEL_NOTFOUND
		}
		e.Value = newOrder
		ob.buys.addQty(newOrder.GetQty().Sub(oldOrder.GetQty()))
		orderlist.addQty(newOrder.GetQty().Sub(oldOrder.GetQty()))
		return nil
	} else if side == SIDE_SELL {
		e, found := ob.sellOrders[oldOrder.GetClOrdID()]
		if !found {
			return ERR_ORDER_NOT_EXISTS
		}
		orderlist, found := ob.sells.prices[newOrder.GetPrice().String()]
		if !found {
			return ERR_PRICELEVEL_NOTFOUND
		}
		e.Value = newOrder
		ob.sells.addQty(newOrder.GetQty().Sub(oldOrder.GetQty()))
		orderlist.addQty(newOrder.GetQty().Sub(oldOrder.GetQty()))
		return nil
	}
	return ERR_SIDE_INVALID
}

func (ob *OrderBook) AddOrder(side OrderSide, order Order) error {
	if side == SIDE_BUY {
		if _, found := ob.buyOrders[order.GetClOrdID()]; !found {
			ord := ob.buys.Add(order)
			ob.buyOrders[order.GetClOrdID()] = ord
			return nil
		}
		return ERR_DULPLICATE_ORDER
	} else if side == SIDE_SELL {
		if _, found := ob.sellOrders[order.GetClOrdID()]; !found {
			ord := ob.sells.Add(order)
			ob.sellOrders[order.GetClOrdID()] = ord
			return nil
		}
		return ERR_DULPLICATE_ORDER
	}
	return ERR_SIDE_INVALID

}

func (ob *OrderBook) Buys() *OrderTree {
	return ob.buys
}

func (ob *OrderBook) Sells() *OrderTree {
	return ob.sells
}

func (ob *OrderBook) SecurityID() string {
	return ob.securityID
}