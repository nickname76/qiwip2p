package qiwip2p

import (
	"time"

	"github.com/google/go-querystring/query"
)

// Необязательные параметры для выставления сч1та.
type OplataCreateOptions struct {
	// Идентификатор выставляемого сч1та в вашей системе. Может быть длиной <=200.
	// Он должен быть уникальным, и генерироваться на вашей стороне любым способом.
	// Идентификатором может быть любая уникальная последовательность букв или цифр.
	// Также разрешено использование символа подчеркивания (_) и дефиса (-).
	BillID BillID `url:"billId,omitempty"`
	// Сумма, на которую выставляется счёт, округленная в меньшую сторону до 2 десятичных знаков.
	// Должна быть не больше 6 знаков до запятой и не больше 2 знаков после запятой.
	// Пример: `Amount: "123456.78"`
	Amount string `url:"amount,omitempty"`
	// Номер телефона пользователя (в международном формате).
	Phone string `url:"phone,omitempty"`
	// E-mail пользователя.
	Email string `url:"email,omitempty"`
	// Идентификатор пользователя в вашей системе.
	Account string `url:"account,omitempty"`
	// Комментарий к счету. Может быть длиной <=255.
	Comment string `url:"comment,omitempty"`
	// Дополнительные данные счета. Может быть длиной суммарно <=255.
	// Может также включать поля (см. константы этого модуля):
	// CustomFieldOptionPaySourcesFilter (с возможными значениями: PaySourceQiwi, PaySourceCard, PaySourceQiwiAndCard)
	// и CustomFieldOptionThemeCode
	CustomFields map[string]string `url:"customFields,omitempty"`
	// Дата и время по МСК (UTC+3), до которого счёт будет актуален. Формат - ГГГГ-ММ-ДДTччмм, см. ConvertTimeToLifetimeValue().
	// Если перевод по счету не будет произведен до этой даты,
	// ему присваивается финальный статус `EXPIRED` и последующий перевод станет невозможен.
	// **Внимание! По истечении 45 суток от даты выставления счёт автоматически будет переведен в финальный статус**
	Lifetime string `url:"lifetime,omitempty"`
	// URL для переадресации на ваш сайт в случае успешного перевода
	SuccessURL string `url:"successUrl,omitempty"`
}

// Создаёт ссылку при переходе по которой отображается форма с выбором способа перевода.
// publicKey - публичный ключ Qiwi P2p API. Не может быть пустым.
// При использовании этого способа нельзя гарантировать, что
// все счета выставлены вами, в отличие от выставления счёта по API.
//
// https://developer.qiwi.com/ru/p2p-payments/#http-invoice
func (api *API) CreateOplataURL(options *OplataCreateOptions) string {
	// Ошибка не может возникать в этом вызове
	v, err := query.Values(options)
	if err != nil {
		panic(err)
	}

	v.Add("publicKey", api.PublicKey)

	queryStr := v.Encode()

	return "https://oplata.qiwi.com/create?" + queryStr
}

// Форматирует time.Time в значение для `OplataCreateOptions.Lifetime`
func FormatTimeToLifetime(t time.Time) string {
	return t.Format("2006-01-02T150405")
}

// Парсит значение `OplataCreateOptions.Lifetime` в time.Time
func ParseLifetime(str string) (time.Time, error) {
	return time.Parse("2006-01-02T150405", str)
}
