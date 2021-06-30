# orderbook
orderbook for trading match engine 交易撮合引擎订单簿

# Installation
```
go get github.com/san-pang/orderbook
```

# Example:  
```
import github.com/san-pang/orderbook

buyorder := newMyOrder("010", "ord1", decimal.NewFromFloat(12.34), decimal.NewFromInt(300), SIDE_BUY)
sellorder := newMyOrder("030", "ord3", decimal.NewFromFloat(12.36), decimal.NewFromInt(500), SIDE_SELL)
orderbook := NewOrderBook("000776")

// add order
orderbook.AddOrder(buyorder.GetSide(), buyorder)
orderbook.AddOrder(sellorder.GetSide(), sellorder)

// cancel order
if err := orderbook.CancelOrder(buyorder.GetSide(), buyorder.GetClOrdID()); err != nil {
  fmt.Println("cancel order failed:", err)
}

// update order
replaceOrder := *sellorder
replaceOrder.Qty = decimal.NewFromInt(1100)
if err := orderbook.UpdateOrder(replaceOrder.GetSide(), sellorder, &replaceOrder); err != nil {
  fmt.Println("update order failed:", err)
}
```

more examples can be found in test file.
