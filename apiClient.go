package apiclient

import "io"

//APIClient interface for all excahnge api
type APIClient interface {
	/* Инициализирует поля реализации. В параметрах принимает приватные значения для
	 * доступа к api. Если конкретная реализация не использует какой-то параметр, при
	 * вызове ф-ии он остается пустым.
	 * logger надо вызвать log.SetOutput(logger), чтобы все log.* функции писали куда надо
	 */
	Init(accountID string, apiKey string, apiSecret string) error

	/* Получает список пар и балансы валют, переводит в формат apiclient и заполняет
	 * соотв. поля реализации. Возвращает указатель на таблицу балансов реализации.
	 *
	 * Функция должна быть вызвана первой, до вызова любых функций интерфейса, кроме
	 * функций-свойств (Proc_XX). Затем ф-ия вызывается после каждого выполненного
	 * ордера для проверки балансов. Следует избегать блокировок.
	 */
	GetBalances() (*map[string]Balance, error)

	/* Получает список всех открытх сделок на бирже(не только наших), заполняет структуру в формате apiclient
	 * и возвращает указтель на эту структуру. Ф-ия вызывается часто, перед каждой
	 * сделкой для анализа orderBook. Следует избегать блокировок.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 * 		в вид который принимает биржа, необходимо привести самому в этой функции
	 *
	 * должен отсортировать asks по возрастанию цены
	 * bids по убыванию цены
	 */
	GetOrderBook(symbol string) (*OrderBook, error)

	/* Ставит ордер на продажу.
	 * теперь не вызывает getOrderInfo
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   amount   кол-во
	 *   price    цена
	 *
	 * возвращает структуру без заполненных полей *Executed если ордер еще не исполнился
	 * в крайнем случае заполняет только ID
	 */
	Sell(symbol string, amount float64, price float64) (*MakedOrder, error)

	/* Ставит ордер на покупку.
	 * теперь не вызывает getOrderInfo
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   amount   кол-во
	 *   price    цена
	 *
	 * возвращает структуру без заполненных полей *Executed если ордер еще не исполнился
	 * в крайнем случае заполняет только ID
	 */
	Buy(symbol string, amount float64, price float64) (*MakedOrder, error)

	/* Отменяет ордер.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   id       id ордера в строковом формате биржы
	 */
	CancelOrder(symbol string, id string) error

	/* Получает информацию об ордере.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   id       id ордера в строковом формате биржы
	 */
	GetOrderStatus(id string, symbol string) (*MakedOrder, error)

	/* Получает список всех открытых сделок на данном аккаунте, заполняет структуру
	 * в формате apiclient и возвращает указтель на эту структуру. Вызывается часто,
	 * следует избегать блокировок.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 */
	GetMyOpenOrders(symbol string) (*[]MakedOrder, error)

	/* Получает список всех сделок на данном аккаунте, заполняет структуру
	 * в формате apiclient и возвращает указтель на эту структуру. Вызывается часто,
	 * следует избегать блокировок.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 */
	GetMyAllOrders(symbol string) (*[]MakedOrder, error)

	/* Создает запрос на вывод средств.
	 *
	 * arguments:
	 *   asset     id валюты в строковом формате: <raw-line, uppercase> (BTC)
	 *   address   blockchain address
	 *   amount    кол-во
	 *
	 * returns:
	 *   id    id запроса на вывод в строковом формате биржы
	 *   err   ошибка
	 */
	Withdraw(asset string, address string, amount float64) (id string, err error)

	/* Функции доступа к приватным полям реализации. Используются абстрактными
	 * функциями во внешних модулях (xutils) для упрощения доступа к данным.
	 * Возвращают соответствующие значения если поле представлено и пустую
	 * строку если данное поле не используется.
	 *
	 * Пример реализации:
	 *   func (A *XX_API) Prop_Root() string  { return A.root   }
	 *   func (A *XX_API) Prop_Key()  string  { return A.apiKey }
	 *   func (A *XX_API) Prop_ID()  string  { return A.apiID }
	 *   func (A *XX_API) Prop_Logger()  io.Writer  { return A.logger }
	 *   func (A *XX_API) Prop_Secret() string { return ""       } // нет в XX_API{}
	 */
	PropRoot() string
	PropKey() string
	PropSecret() string
	PropID() string
	PropLogger() io.Writer
}

//Status type for enum about order status
type Status string

//Side type for enum about order buy or sell types or something like that
type Side string

//constants about Status and Side
const (
	Buy             Side   = "BUY"
	Sell            Side   = "SELL"
	Filled          Status = "FILLED"
	NotFilled       Status = "NotFilled"
	PartiallyFilled Status = "PartiallyFilled"
)

//Balance help struct for APIClient
type Balance struct {
	Free   float64 `json:"free"`   //available balance for use in new orders
	Locked float64 `json:"locked"` //locked balance in orders or withdrawals
}

//OrderBook help struct for APIClient
type OrderBook struct {
	Asks []Order `json:"asks"` //asks.Price > any bids.Price
	Bids []Order `json:"bids"`
}

//Order help struct for APIClient
type Order struct {
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

//MakedOrder help struct for APIClient
type MakedOrder struct {
	ID string `json:"id"`
	//  Status Should be one of apiclient.Status constants(Filled, NotFilled, PartiallyFilled)
	Status      Status  `json:"status"`
	Closed      bool    `json:"closed"`
	LeftAmount  float64 `json:"leftAmount"`
	RightAmount float64 `json:"rightAmount"`

	LeftAmountExecuted  float64 `json:"leftAmountExecuted"`
	RightAmountExecuted float64 `json:"rightAmountExecuted"`
	//Commission is factically not used
	Commission   float64 `json:"commission"`
	Rate         float64 `json:"rate"`
	RateExecuted float64 `json:"rateExecuted"`
	//  Side Should be one of apiclient.Side constants(Buy, Sell)
	Side Side `json:"side"`
}
