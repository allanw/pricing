package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type OrderItem struct {
	ProductId int `json:"product_id"`
	Quantity  int
}

type Order struct {
	Order struct {
		Id    int
		Items []OrderItem
	}
}

type OrderItemInfo struct {
	ProductId int
	Price     int
	Vat       float32
	Quantity  int
}

type OrderInfo struct {
	OrderId    int
	TotalPrice int
	TotalVat   float32
	OrderItems []OrderItemInfo
}

type Price struct {
	ProductId int `json:"product_id"`
	Price     int
	VatBand   string `json:"vat_band"`
}

type PricingInfo struct {
	Prices   []Price
	VatBands map[string]float32 `json:"vat_bands"`
}

func main() {
	http.HandleFunc("/order", order)

	http.ListenAndServe(":8080", nil)
}

func WriteOrderToCsvFile(order_info *OrderInfo, path string) {
	_, err := os.Stat(path)
	if err != nil {
		// file doesn't exist, create it
		file, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		writer.Write([]string{"order_id", "total_price", "total_vat", "product_id", "product_price", "quantity"})

		for _, item := range order_info.OrderItems {
			writer.Write([]string{strconv.Itoa(order_info.OrderId),
				strconv.Itoa(order_info.TotalPrice),
				strconv.FormatFloat(float64(order_info.TotalVat), 'f', 1, 32),
				strconv.Itoa(item.ProductId),
				strconv.Itoa(item.Price),
				strconv.Itoa(item.Quantity)})
		}
	} else {
		// file already exists
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()

		for _, item := range order_info.OrderItems {
			writer.Write([]string{strconv.Itoa(order_info.OrderId),
				strconv.Itoa(order_info.TotalPrice),
				strconv.FormatFloat(float64(order_info.TotalVat), 'f', 1, 32),
				strconv.Itoa(item.ProductId),
				strconv.Itoa(item.Price),
				strconv.Itoa(item.Quantity)})
		}
	}
}

func order(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var order Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		panic(err)
	}

	raw, err := ioutil.ReadFile("pricing.json")

	var pricing_data PricingInfo
	json.Unmarshal(raw, &pricing_data)

	price_map := make(map[int]Price)
	vat_map := make(map[string]float32)

	for _, price := range pricing_data.Prices {
		price_map[price.ProductId] = Price{Price: price.Price, VatBand: price.VatBand}
	}

	for band_name, vat_rate := range pricing_data.VatBands {
		vat_map[band_name] = vat_rate
	}

	order_id := order.Order.Id
	var total_order_price int
	var total_order_vat float32

	var order_items []OrderItemInfo

	for _, o := range order.Order.Items {
		order_item_price := price_map[o.ProductId].Price
		order_item_vat := float32(price_map[o.ProductId].Price) * float32(o.Quantity) * vat_map[price_map[o.ProductId].VatBand]
		order_items = append(order_items, OrderItemInfo{ProductId: o.ProductId, Price: order_item_price, Vat: order_item_vat, Quantity: o.Quantity})
		total_order_price += price_map[o.ProductId].Price * o.Quantity
		total_order_vat += order_item_vat
	}

	order_info := OrderInfo{OrderId: order_id, TotalPrice: total_order_price, TotalVat: total_order_vat, OrderItems: order_items}

	res, err := json.MarshalIndent(order_info, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", res)

	WriteOrderToCsvFile(&order_info, "orders.csv")

}
