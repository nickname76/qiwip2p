package qiwip2p

import "time"

// Идентификатор счёта
type BillID string

// Идентификатор в сервисе приёма платежей
type SiteID string

// Данные о счёте
type BillMetadata struct {
	// Ваш идентификатор в сервисе приема платежей для физических лиц p2p.qiwi.com.
	// НЕ ИСПОЛЬЗОВАТЬ С CreateBill().
	SiteID SiteID `json:"siteId,omitempty"`
	// Идентификатор выставляемого счёта в вашей системе. Может быть длиной <=200.
	// Он должен быть уникальным, и генерироваться на вашей стороне любым способом.
	// Идентификатором может быть любая уникальная последовательность букв или цифр.
	// Также разрешено использование символа подчеркивания (_) и дефиса (-).
	// НЕ ИСПОЛЬЗОВАТЬ С CreateBill().
	BillID BillID `json:"billId,omitempty"`
	// Ссылка на созданную форму. Перенаправьте пользователя
	// по этой ссылке для оплаты счета или используйте
	// библиотеку Popup, чтобы открыть форму во всплывающем окне.
	// https://developer.qiwi.com/ru/p2p-payments/#popup
	// НЕ ИСПОЛЬЗОВАТЬ С CreateBill().
	PayURL string `json:"payUrl,omitempty"`
	// Идентификаторы пользователя.
	Customer *CustomerIdentificators `json:"customer,omitempty"`
	// Комментарий к счету.
	Comment string `json:"comment,omitempty"`
	// Информация о сумме счета.
	Amount *BillMetadataAmount `json:"amount,omitempty"`
	// Информация о статусе счета.
	// НЕ ИСПОЛЬЗОВАТЬ С CreateBill().
	Status *BillMetadataStatus `json:"status,omitempty"`
	// Объект строковых дополнительных параметров. Возможные элементы: paySourcesFilter, themeCode
	// Может также включать поля (см. константы этого модуля):
	// CustomFieldOptionPaySourcesFilter (с возможными значениями: PaySourceQiwi, PaySourceCard, PaySourceQiwiAndCard)
	// и CustomFieldOptionThemeCode.
	CustomFields map[string]string `json:"customFields,omitempty"`
	// Системная дата создания счета. Формат даты: `ГГГГ-ММ-ДДTчч:мм:сс±чч:мм`.
	// НЕ ИСПОЛЬЗОВАТЬ С CreateBill().
	CreationDateTime string `json:"creationDateTime,omitempty"`
	// Срок действия созданной формы для перевода. Формат даты: `ГГГГ-ММ-ДДTчч:мм:сс±чч:мм`.
	ExpirationDateTime string `json:"expirationDateTime,omitempty"`
}

// Идентификаторы пользователя
type CustomerIdentificators struct {
	// Номер телефона пользователя (в международном формате).
	Phone string `json:"phone,omitempty"`
	// E-mail пользователя.
	Email string `json:"email,omitempty"`
	// Идентификатор пользователя в вашей системе.
	Account string `json:"account,omitempty"`
}

// Валюта счёта
type Currency string

const (
	// Российский рубль
	CurrencyRUB Currency = "RUB"
	// Казахстанский тенге
	CurrencyKZT Currency = "KZT"
)

// Сумма счёта
type BillMetadataAmount struct {
	// Сумма счета, округленная до 2 знаков после запятой в меньшую сторону.
	Value string `json:"value,omitempty"`
	// Валюта суммы счета (Alpha-3 ISO 4217 код). Доступные значения: CurrencyRUB, CurrencyKZT (см. константы модуля)
	Currency Currency `json:"currency,omitempty"`
}

// Статус оплаты счета.
type BillMetadataStatusValue string

const (
	// Счет выставлен, ожидает оплаты.
	BillMetadataStatusValueWaiting BillMetadataStatusValue = "WAITING"
	// Счет оплачен. Финальный.
	BillMetadataStatusValuePaid BillMetadataStatusValue = "PAID"
	// Время жизни счета истекло. Счет не оплачен. Финальный.
	BillMetadataStatusValueExpired BillMetadataStatusValue = "EXPIRED"
	// Счет отклонен. Финальный.
	BillMetadataStatusValueRejected BillMetadataStatusValue = "REJECTED"
)

type BillMetadataStatus struct {
	// Текущий статус счета.
	Value BillMetadataStatusValue `json:"value,omitempty"`
	// Дата обновления статуса. Формат даты: `ГГГГ-ММ-ДДTчч:мм:сс±чч:мм`.
	ChangedDateTime string `json:"changedDateTime,omitempty"`
}

const (
	// При открытии формы будут отображаться только указанные
	// способы перевода, если они доступны. Возможные значения см. ниже.
	CustomFieldOptionPaySourcesFilter = "paySourcesFilter"
	// QIWI Кошелек.
	PaySourceQiwi = "qw"
	// Банковская карта.
	PaySourceCard = "card"
	// QIWI Кошелек и банковская карта.
	PaySourceQiwiAndCard = PaySourceQiwi + "," + PaySourceCard
)

const (
	// Код персонализации вашей формы. Может быть длиной <=255.
	// https://developer.qiwi.com/ru/p2p-payments/#custom
	CustomFieldOptionThemeCode = "themeCode"
)

// Форматирует time.Time в значение формата для полей с datetime в BillMetadata
func FormatBillMetadataDateTime(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

// Парсит значение поля с datetime из BillMetedata
func ParseBillMetadataDateTime(changedDateTime string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, changedDateTime)
}
