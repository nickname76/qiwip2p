# Qiwi P2P-счета

Библиотека для работы с Qiwi P2P-счетами. Поддерживает создание ссылок для создания счёта на стороне клиента и непосредственную работу с API P2P-счетов.

Документация P2P-счетов: https://developer.qiwi.com/ru/p2p-payments/

Документация этой библиотеки: https://pkg.go.dev/github.com/nickname76/qiwip2p

*Пожалуйсте, постаьте **звезду** этому репозиторию, если вам пригодилась данная библиотека.*

## Пример использования

```Go
package main

import (
	"log"
	"time"

	"github.com/nickname76lib/qiwip2p"
)

func main() {
	api := qiwip2p.NewAPI("PUBLIC_KEY", "SECRET_KEY")
	bill, err := api.CreateBill("test12345", &qiwip2p.BillMetadata{
		Comment: "Test Comment",
		Amount: &qiwip2p.BillMetadataAmount{
			Value:    "99.99",
			Currency: qiwip2p.CurrencyRUB,
		},
		CustomFields: map[string]string{
			qiwip2p.CustomFieldOptionThemeCode: "THEME_CODE",
			"test_custom_field":                "123456",
		},
		ExpirationDateTime: qiwip2p.FormatBillMetadataDateTime(time.Now().Add(time.Hour)),
	})
	if err != nil {
		log.Fatalf("Ошибка: %w", err)
	}

	log.Println(bill)
}

```
