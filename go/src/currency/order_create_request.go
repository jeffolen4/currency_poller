package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type OrderCreateRequest struct {
	Order OrderCreate `json:"order"`
}

type OrderCreate struct {
	Units        string `json:"units"`
	Instrument   string `json:"instrument"`
	TimeInForce  string `json:"timeInForce"`
	Type         string `json:"type"`
	PositionFill string `json:"positionFill"`
}

type OrderCreateResponse struct {
	OrderCreateRequest OrderCreateRequest `json:"orderCreateRequest"`
	OrderCancelTransaction OrderTransaction `json:"orderCancelTransaction"`
	OrderFillTransaction OrderTransaction `json:"orderFillTransaction"`
	OrderCreateTransaction OrderTransaction `json:"orderCreateTransaction"`
}

type OrderTransaction struct {
	Id string `json:"id"`
	Time string `json:"time"`
}

func main() {
	// Replace with your own values
	apiKey := "YOUR_API_KEY"
	accountID := "YOUR_ACCOUNT_ID"
	instrument := "EUR_USD"

	orderCreateResponse, err := createOrder(apiKey, accountID, instrument, "1000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Order created: %v\n", orderCreateResponse.OrderCreateTransaction.Id)
}

func createOrder(apiKey string, accountID string, instrument string, units string) (*OrderCreateResponse, error) {
	orderCreateRequest := OrderCreateRequest{
		Order: OrderCreate{
			Units:        units,
			Instrument:   instrument,
			TimeInForce:  "FOK",
			Type:         "MARKET",
			PositionFill: "DEFAULT",
		},
	}

	payload, err := json.Marshal(orderCreateRequest)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api-fxtrade.oanda.com/v3/accounts/"+accountID+"/orders", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var orderCreateResponse OrderCreateResponse

	err = json.NewDecoder(resp.Body).Decode(&orderCreateResponse)
	if err != nil {
		return nil, err
	}

	return &orderCreateResponse, nil
}