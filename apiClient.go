package apiclient

import (
	"fmt"
)

type DebugStr struct {
	LastRequestMethod string
	LastResponseData  []byte
	LastSentData      []byte
}

func (D *DebugStr) FormatError() string {
	return fmt.Sprintf("method: %s, rd: %s, sd: %s",
		D.LastRequestMethod,
		D.LastResponseData,
		D.LastSentData,
	)
}

type ExchangeAPI struct {
	ApiKey    string
	ApiSecret string
	AccountID string

	Root string

	Depth int

	DebugStr
}

//APIClient interface for all exchange api
type APIClient interface {
	/* Инициализирует поля реализации. В параметрах принимает приватные значения для
	 * доступа к api. Если конкретная реализация не использует какой-то параметр, при
	 * вызове ф-ии он остается пустым.
	 * logger надо вызвать log.SetOutput(logger), чтобы все log.* функции писали куда надо
	 */
	Init(apiKey, apiSecret, accountID string, depth int) error

	/* Получает список пар и балансы валют, переводит в формат apiclient и заполняет
	 * соотв. поля реализации. Возвращает карту балансов реализации и ошибку.
	 */
	GetBalances() (map[string]Balance, error)

	/* Ставит ордер на покупку.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   amount   кол-во
	 *   price    цена
	 *
	 * возвращает структуру без заполненных полей *Executed если ордер еще не исполнился
	 * в крайнем случае заполняет только ID
	 */
	Buy(symbol string, amount, price float64) (MakedOrder, error)

	/* Ставит лимитный ордер на продажу.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   amount   кол-во
	 *   price    цена
	 *
	 * возвращает структуру без заполненных полей *Executed если ордер еще не исполнился
	 * в крайнем случае заполняет только ID
	 */
	Sell(symbol string, amount, price float64) (MakedOrder, error)

	/* Отменяет ордер.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   id       id ордера в строковом формате биржы
	 */
	CancelOrder(symbol, id string) error

	/* Отменяет все открыте ордера на заданной паре.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 */
	CancelAll(symbol string) error

	/* Получает информацию об ордере.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   id       id ордера в строковом формате биржы
	 */
	GetOrderStatus(symbol, id string) (MakedOrder, error)

	/* Получает список ВСЕХ открытых оредров на данном аккаунте для выбранной пары, заполняет структуру
	 * в формате apiclient и возвращает эту структуру. Вызывается часто.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 */
	GetMyOpenOrders(symbol string) ([]MakedOrder, error)

	/* Создает запрос на вывод средств.
	 *
	 * arguments:
	 *   asset     id валюты в строковом формате: <raw-line, uppercase> (BTC)
	 *   address   blockchain address
	 *   chain     blockchain network
	 *   amount    кол-во
	 *
	 * returns:
	 *   id    id запроса на вывод в строковом формате биржы
	 *   err   ошибка
	 */
	Withdraw(asset, address, chain string, amount float64) (id string, err error)

	// Функции ниже в том числе должны работать, если задать пустые настройки api

	/* Получает список всех возможных text-id пар в формате apiclient на бирже.
	 * Сортиурет, так что сначала идут USDT пары, затем BTC, затем ETH, после все остальные.
	 */
	GetAllSymbols() ([]string, error)

	/* Получает кол-во знаков после запятой для price и amount, заполняет структуру в формате apiclient
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 */
	GetDecs(symbol string) (Decimals, error)

	/* Получает биржевой стакан с глубиной, заданной в Init, заполняет структуру в формате apiclient
	 * и сортирует Asks по возрастанию цены и Bids по убыванию цены.
	 * Ф-ия вызывается часто, перед каждой сделкой для анализа orderBook.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 * 		в вид который принимает биржа, необходимо привести самому в этой функции
	 *
	 * требуется отсортировать Asks по возрастанию цены
	 * Bids по убыванию цены
	 */
	GetOrderBook(symbol string) (OrderBook, error)

	/* Получает цену последней сделки на бирже.
	 * Ф-ия вызывается часто.
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 * 		в вид который принимает биржа, необходимо привести самому в этой функции
	 */
	GetLastPrice(symbol string) (float64, error)

	/* Получает 100 последних объемных свеч и свеч цены, заполняет структуру в формате apiclient
	 * и сортирует по времени (от болле старого к более новому)
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта в котрой цена, справа валюта в которой кол-во)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 *   candlePeriod  период свеч в минутах, типа int
	 */
	GetKLine(symbol string, candlePeriod int) (KLine, error)

	/* Получает 100 последних сделок, заполняет структуру в формате apiclient
	 * и сортирует ее по времени (от болле старого к более новому)
	 *
	 * arguments:
	 *   symbol   text-id пары в формате apiclient
	 * 		(например "BTC_ETH" слева валюта основная, справа валюта которую покупаем при вызове buy)
	 *	 	в вид который принимает биржа, необходимо привести самому в этой функции
	 */
	GetTradeHistory(symbol string) ([]Trade, error)
}

//Side type for enum about order buy or sell types or something like that
type Side string

//Status type for enum about order status
type Status string

//Color type for volume candle color
type Color string

//constants about Status and Side
const (
	Buy             Side   = "BUY"
	Sell            Side   = "SELL"
	Filled          Status = "FILLED"
	NotFilled       Status = "NOT_FILLED"
	PartiallyFilled Status = "PARTIALLY_FILLED"
	Red             Color  = "rgba(255, 82, 82, 0.5)"
	Green           Color  = "rgba(0, 150, 136, 0.5)"
)

//Balance help struct for APIClient
type Balance struct {
	Free   float64 `json:"free"`   // available balance for use in new orders
	Locked float64 `json:"locked"` // locked balance in orders or withdrawals
}

//OrderBook help struct for APIClient
type OrderBook struct {
	Asks []Order `json:"asks"` // asks.Price > any bids.Price
	Bids []Order `json:"bids"`
}

//Order help struct for APIClient
type Order struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

//MakedOrder help struct for APIClient
type MakedOrder struct {
	ID string `json:"id"`
	//  Status Should be one of apiclient.Status constants(Filled, NotFilled, PartiallyFilled)
	Status      Status  `json:"status"`
	LeftAmount  float64 `json:"leftAmount"`
	RightAmount float64 `json:"rightAmount"`

	LeftAmountExecuted  float64 `json:"leftAmountExecuted"`
	RightAmountExecuted float64 `json:"rightAmountExecuted"`

	Commission    float64 `json:"commission"`
	Price         float64 `json:"price"`
	PriceExecuted float64 `json:"priceExecuted"` // Real deal price
	//  Side Should be one of apiclient.Side constants(Buy, Sell)
	Side Side `json:"side"`
}

//PriceCandle help struct for APIClient
type PriceCandle struct {
	Time  int64   `json:"time"`  // UNIX time in seconds (10 digits)
	Open  float64 `json:"open"`  // open price for period
	Close float64 `json:"close"` // close price for period
	High  float64 `json:"high"`  // high price for period
	Low   float64 `json:"low"`   // low price for period
}

//VolumeCandle help struct for APIClient
type VolumeCandle struct {
	Time  int64   `json:"time"`  // UNIX time in seconds (10 digits)
	Value float64 `json:"value"` // volume for period
	Color Color   `json:"color"` // apiclient.Green if Close > Open, apiclient.Red if Close < Open
}

//KLine help struct for APIClient
type KLine struct {
	PriceCandles  []PriceCandle  `json:"priceCandles"`
	VolumeCandles []VolumeCandle `json:"volumeCandles"`
}

type Trade struct {
	Time   int64   `json:"time"` // UNIX time in seconds (10 digits)
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
	Side   Side    `json:"side"` // if side not defined, = Buy if current price > last price, = Sell if current price < last price
}

type Decimals struct {
	PriceDecimal  int `json:"priceDecimal"`
	AmountDecimal int `json:"amountDecimal"`
}
