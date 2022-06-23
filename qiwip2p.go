// https://developer.qiwi.com/ru/p2p-payments/
package qiwip2p

import (
	"encoding/json"
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

// Интерфейс для взаимодействия с API P2P-счетов по стандарту OAuth 2.0.
type API struct {
	// Публичный ключ для авторизации при выставлении счетов через форму.
	PublicKey string
	// Секретный ключ для авторизации запросов к API
	SecretKey string
	// Функция, которая используется для http запросов к Qiwi P2P API
	HttpDoRequest func(method string, url string, headers map[string]string, body []byte) (respStatusCode int, respBody []byte, err error)
}

// Создаёт новый интерфейс для взаимодействия с API P2P-счетов Qiwi.
// Если вы хотите кастомизировать выполнение HTTP запросов к API, создавайте объект API самостоятельно,
// в поле HttpDoRequest можно указать свою функцию для выполнения HTTP запросов.
func NewAPI(publicKey string, secretKey string) *API {
	return &API{
		PublicKey:     publicKey,
		SecretKey:     secretKey,
		HttpDoRequest: DefaultHttpDoRequest,
	}
}

var defaultFasthttpClient = &fasthttp.Client{
	NoDefaultUserAgentHeader:      true,
	DisableHeaderNamesNormalizing: true,
	DisablePathNormalizing:        true,
}

func DefaultHttpDoRequest(method string, url string, headers map[string]string, body []byte) (respStatusCode int, respBody []byte, err error) {
	req := &fasthttp.Request{}
	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.SetBody(body)

	resp := &fasthttp.Response{}

	err = defaultFasthttpClient.Do(req, resp)
	if err != nil {
		return 0, nil, fmt.Errorf("DefaultHttpDoRequest: %w", err)
	}

	respStatusCode = resp.StatusCode()
	respBody = resp.Body()

	return
}

func (api *API) makeAPICall(httpMethod string, pathCompound string, billMeta *BillMetadata) (*BillMetadata, error) {
	requestURL := "https://api.qiwi.com/partner/bill/v1/bills/" + pathCompound

	var body []byte

	if billMeta != nil {
		jsoniterCfg := jsoniter.Config{
			EscapeHTML:                    false,
			OnlyTaggedField:               true,
			ObjectFieldMustBeSimpleString: true,
			CaseSensitive:                 true,
		}.Froze()

		var err error
		body, err = jsoniterCfg.Marshal(billMeta)
		if err != nil {
			return nil, fmt.Errorf("makeAPICall: %w", err)
		}
	}

	respStatusCode, respBody, err := api.HttpDoRequest(httpMethod, requestURL, map[string]string{
		"Authorization": "Bearer " + api.SecretKey,
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}, body)
	if err != nil {
		return nil, fmt.Errorf("makeAPICall: %w", err)
	}

	if respStatusCode < 200 || respStatusCode >= 300 {
		return nil, fmt.Errorf("makeAPICall: qiwi p2p api error (status code %v): %w", respStatusCode, errors.New(string(respBody)))
	}

	responseData := &BillMetadata{}
	err = json.Unmarshal(respBody, responseData)
	if err != nil {
		return nil, fmt.Errorf("makeAPICall: %w", err)
	}

	return responseData, nil
}

// Выставление счета
//
// https://developer.qiwi.com/ru/p2p-payments/#create
func (api *API) CreateBill(billID BillID, billMeta *BillMetadata) (*BillMetadata, error) {
	respBillMeta, err := api.makeAPICall("PUT", string(billID), billMeta)
	if err != nil {
		return nil, fmt.Errorf("CreateBill: %w", err)
	}
	return respBillMeta, nil
}

// Проверка статуса перевода по счету
//
// https://developer.qiwi.com/ru/p2p-payments/#invoice-status
func (api *API) GetBill(billID BillID) (*BillMetadata, error) {
	respBillMeta, err := api.makeAPICall("GET", string(billID), nil)
	if err != nil {
		return nil, fmt.Errorf("GetBill: %w", err)
	}
	return respBillMeta, nil
}

// Отмена неоплаченного счета
//
// https://developer.qiwi.com/ru/p2p-payments/#cancel
func (api *API) CancelBill(billID BillID) (*BillMetadata, error) {
	respBillMeta, err := api.makeAPICall("POST", string(billID)+"/reject", nil)
	if err != nil {
		return nil, fmt.Errorf("CancelBill: %w", err)
	}
	return respBillMeta, nil
}
