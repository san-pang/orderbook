package orderbook

import (
	"container/list"
	"fmt"
	rbtx "github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
	"strings"
)

type OrderTree struct {
	priceTree *rbtx.RedBlackTreeExtended
	prices    map[string]*OrderList
	qty    decimal.Decimal
	orderSize int
	depth     int
}

func NewOrderTree() *OrderTree {
	return &OrderTree{
		priceTree: &rbtx.RedBlackTreeExtended{
			Tree: rbt.NewWith(rbtComparator),
		},
		prices: make(map[string]*OrderList),
	}
}

func rbtComparator(a, b interface{}) int {
	return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
}

func (ot *OrderTree) Len() int {
	return ot.orderSize
}

func (ot *OrderTree) Depth() int {
	return ot.depth
}

func (ot *OrderTree) Qty() decimal.Decimal {
	return ot.qty
}

func (ot *OrderTree) addQty(qty decimal.Decimal) {
	ot.qty = ot.qty.Add(qty)
}

func (ot *OrderTree) subQty(qty decimal.Decimal) {
	ot.qty = ot.qty.Sub(qty)
}

func (ot *OrderTree) Add(o Order) *list.Element {
	price := o.GetPrice()
	strPrice := price.String()
	priceQueue, ok := ot.prices[strPrice]
	if !ok {
		priceQueue = NewOrderList(o.GetPrice())
		ot.prices[strPrice] = priceQueue
		ot.priceTree.Put(price, priceQueue)
		ot.depth++
	}
	ot.orderSize++
	ot.addQty(o.GetQty())
	return priceQueue.Push(o)
}

func (ot *OrderTree) Remove(e *list.Element) Order {
	price := e.Value.(Order).GetPrice()
	strPrice := price.String()
	priceQueue := ot.prices[strPrice]
	o := priceQueue.Remove(e)
	if priceQueue.Len() == 0 {
		delete(ot.prices, strPrice)
		ot.priceTree.Remove(price)
		ot.depth--
	}
	ot.orderSize--
	ot.subQty(o.GetQty())
	return o
}

func (ot *OrderTree) MaxPriceOrderList() *OrderList {
	if ot.depth > 0 {
		if value, found := ot.priceTree.GetMax(); found {
			return value.(*OrderList)
		}
	}
	return nil
}

func (ot *OrderTree) MinPriceOrderList() *OrderList {
	if ot.depth > 0 {
		if value, found := ot.priceTree.GetMin(); found {
			return value.(*OrderList)
		}
	}
	return nil
}

func (ot *OrderTree) LessThan(price decimal.Decimal) *OrderList {
	tree := ot.priceTree.Tree
	node := tree.Root

	var floor *rbt.Node
	for node != nil {
		if tree.Comparator(price, node.Key) > 0 {
			floor = node
			node = node.Right
		} else {
			node = node.Left
		}
	}

	if floor != nil {
		return floor.Value.(*OrderList)
	}

	return nil
}

func (ot *OrderTree) LessEqualThan(price decimal.Decimal) *OrderList {
	tree := ot.priceTree.Tree
	node := tree.Root

	var floor *rbt.Node
	for node != nil {
		if tree.Comparator(price, node.Key) >= 0 {
			floor = node
			node = node.Right
		} else {
			node = node.Left
		}
	}

	if floor != nil {
		return floor.Value.(*OrderList)
	}

	return nil
}

func (ot *OrderTree) GreaterThan(price decimal.Decimal) *OrderList {
	tree := ot.priceTree.Tree
	node := tree.Root

	var ceiling *rbt.Node
	for node != nil {
		if tree.Comparator(price, node.Key) < 0 {
			ceiling = node
			node = node.Left
		} else {
			node = node.Right
		}
	}

	if ceiling != nil {
		return ceiling.Value.(*OrderList)
	}

	return nil
}

func (ot *OrderTree) GreaterEqualThan(price decimal.Decimal) *OrderList {
	tree := ot.priceTree.Tree
	node := tree.Root

	var ceiling *rbt.Node
	for node != nil {
		if tree.Comparator(price, node.Key) <= 0 {
			ceiling = node
			node = node.Left
		} else {
			node = node.Right
		}
	}

	if ceiling != nil {
		return ceiling.Value.(*OrderList)
	}

	return nil
}

func (ot *OrderTree) Orders() (orders []*list.Element) {
	for _, price := range ot.prices {
		iter := price.Head()
		for iter != nil {
			orders = append(orders, iter)
			iter = iter.Next()
		}
	}
	return
}

func (ot *OrderTree) OrdersByPrice(price decimal.Decimal) (orders []*list.Element) {
	if orderlist, found := ot.prices[price.String()]; found {
		iter := orderlist.Head()
		for iter != nil {
			orders = append(orders, iter)
			iter = iter.Next()
		}
	}
	return
}

func (ot *OrderTree) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("汇总数据：\n qty: %s, orderSize: %d, depth: %d, priceLevel：\n ", ot.Qty(), ot.orderSize, ot.Depth()))
	level := ot.MaxPriceOrderList()
	for level != nil {
		sb.WriteString(level.String() + "\n")
		level = ot.LessThan(level.Price())
	}

	return sb.String()
}