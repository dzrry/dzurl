# Тестовое задание в юнит Авто

# О проекте
Проект полностью соответствует требованиям тестового задания, представленного по ссылке на 
[github](https://github.com/avito-tech/auto-backend-trainee-assignment).

### Реализованные усложнения:
- Написаны тесты
- Добавлена валидация URL с проверкой корректности ссылки (была идея реализовать blacklist для недопустимых ссылок и 
слов в ссылках)
- Добавлена возможность задавать кастомные ссылки, чтобы пользователь мог сделать их человекочитаемыми. Ключ можно 
напрямую передать в теле POST-запроса
- Контейнеризация (docker и docker-compose)

### Также в проекте осознанно были сделаны следующие усложнения, не входящие в начальный список:
- Реализована возможность выбора сериализатора данных: помимо стандартного JSON доступен более быстрый и эффективный по
памяти бинарный msgpack (с его помощью достигается 15-20% экономия трафика по сравнению с JSON). Пример использования 
msgpack будет приведен ниже.
- Для хранения в БД была выбрана усложненная структура данных. Т.к. в качестве хранилища используется Redis, можно было 
воспользоваться более простой структурой для хранения данных: ключом был бы сгенерированный хэш (uuid), а значением - 
полный URL. Но было принято решение разработать архитектуру приложения таким образом, чтобы имелась возможность 
поддерживать несколько видов хранилищ: будь то реляционные/нереляционные базы данных, а не только Key-Value хранилища.

# Запуск микросервиса
```bash
$ git clone https://github.com/dzrry/dzurl.git  
$ cd dzurl
$ docker-compose build
$ docker-compose up
```
После этого API сервер готов получать HTTP запросы на порту ```8080```, но это поведение можно изменить в файле 
конфигурации ```config/config.yml```.

# Использование
API имеет всего два запроса:
- ```POST-запрос``` для сохранения короткого URL в хранилище
- ```GET-запрос``` по ключу для получения короткого URL и перенаправления на соответствующий полный URL

### Создание короткого представления заданного URL
- Для десериализации из msgpack POST-запрос может выглядеть следующим образом:
```bash
$ curl -XPOST http://localhost:8080 \
    -H 'Content-Type: application/x-msgpack' \
    --data-binary $(echo '{"url": "https://start.avito.ru/tech", "key": "start-avito"}' \
    | json2msgpack) | msgpack2json
```
Чтобы использовать десериализацию через msgpack в заголовке запроса необходимо передать ```Content-Type x-msgpack```.
Также в данном примере демонстрируется использование опционального поля ```key``` в JSON для создания кастомной ссылки. В
случае отсутствия данного параметра ключ будет сгенерирован рандомно. Это мы увидим далее.

Ответ сервера после вышеописанного запроса:
```
{"key":"start-avito","url":"https://start.avito.ru/tech","created_at":1600089608}
```
В запросе, чтобы отправить данные в человекочитаемом виде, используется инструмент под названием msgpack-tools.
Скачать его можно через Homebrew:
```bash
$ brew install msgpack-tools
```

- Для десериализации из JSON возможен следующий запрос:
```bash
$ curl -XPOST http://localhost:8080 \
    -H 'Content-Type: application/json' \
    -d '{"url": "https://www.avito.ru/krasnoyarsk"}'
```
Ответ сервера:
```bash
{"key":"btfn593pc98oejp5umq0","url":"https://www.avito.ru/krasnoyarsk","created_at":1600090788}
```
Как мы видим для короткого URL сгенерировался рандомный ключ. Принцип генерации ключа будет рассмотрен в вопросах.

### Получение и переход по сохраненному ранее короткому URL
- Пример GET-запроса:
```bash
$ curl -X GET http://localhost:8080/start-avito
```
Ответ:
```bash
<a href="https://start.avito.ru/tech">Moved Permanently</a>.
```

# Вопросы по заданию
- Каким образом генерировать хеши для коротких URL'ов? - Обычно для генерации случайных последовательностей используют 
UUID, но, по моему мнению, для коротких URL'ов такие последовательности слишком длинные (36 символов). Поэтому в
проекте используется пакет [rs/xid](https://github.com/rs/xid). Данный пакет генерирует строку символов, в 20 символов, 
причем такая строка весит на 4 байта меньше UUID.

- Сколько должны жить ссылки? - Я выбрал решение безвременного хранения ссылок, т.к. сервис будет не очень популярен. 
Но в большом сервисе могло бы произойти захламление мертвыми и одноразовыми ссылками. Возможно, стоило ввести TTL для 
каждого URL или же предусмотреть возможность создания одноразовых ссылок.

- Нужно ли перезаписывать полный URL, если такой короткий URL уже существует? - Я выбрал решение перезаписи, т.к. 
возможность удаления ссылок и, соответственно, ключей не предусмотрена.

# Документация
+ [docker-compose](#dockercompose)
+ [Dockerfile](#dockerfile)
+ [cmd](#cmd)
    + [main.go](#maingo)
        + [func main](#main)
+ [domain](#domain)
    + [domain.go](#domaingo)
        + [type Redirect](#redirect)
+ [repo](#repo)
    + [repo.go](#repogo)
       + [type RedirectRepo](#redirectrepo)
    + [redis](#redis)
        + [repo.go](#redisrepo)
            + [type redisRepo](#redisclient)
            + [func NewRepo](#newrepo)
            + [func newClient](#newclient)
            + [func (r *redisRepo) Load](#redisload)
            + [func (r *redisRepo) Store](#redisstore)
+ [service](#service)
    + [errors.go](#errors)
    + [service.go](#servicego)
        + [type RedirectService](#redirectservice)
        + [type redirectService](#redirectservicest)
        + [func NewRedirectService](#newservice)
        + [func (r *redirectService) Load](#serviceload)
        + [func (r *redirectService) Store](#servicestore)
+ [serialization](#serialization)
    + [serializer.go](#serializer)
        + [type RedirectSerializer](#redirectserializer)
    + [json](#json)
        + [serializer.go](#jsonserializer)
            + [type Redirect](#jsonredirect)
            + [func (r *Redirect) Decode](#jsondecode)   
            + [func (r *Redirect) Encode](#jsonencode)
    + [msgpack](#msgpack)
        + [serializer.go](#msgpackserializer)
            + [type Redirect](#msgpackredirect)
            + [func (r *Redirect) Decode](#msgpackdecode)   
            + [func (r *Redirect) Encode](#msgpackencode)
+ [transport](#transport)
    + [transport.go](#transportgo)
        + [type RedirectHandler](#redirecthandler)
        + [type handler](#handler)
        + [func NewHandler](#newhandler)
        + [func setupResponse](#setupresponse)
        + [func (h *handler) serializer](#handlerserializer)
        + [func (h *handler) LoadRedirect](#transportload)
        + [func (h *handler) StoreRedirect](#transportstore)
+ [config](#config)
    + [config.yml](#configyml)
    + [config.go](#configgo)
        + [type Config](#configst)
        + [type RedisConfig](#rediconfig)
        + [type ServerConfig](#serverconfig)
        + [func Read](#readconfig)
+ [mocks](#mocks)

<a name="dockercompose"></a>
# dockercompose
```yaml
version: "3.8"
services:
  server:
    build: ./
    ports:
      - "8080:8080"
    restart: always
    links:
      - redis
    depends_on:
      - redis

  redis:
    image: redis:alpine
```
Собирает, развертывает и управляет контейнером API сервера и redis. Связывает контейнеры в одну сеть, маппит порты 
сервера.

<a name="Dockerfile"></a>
# dockerfile
```dockerfile
FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /avito-auto-unit/
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -o shortener ./cmd/main.go

FROM scratch

WORKDIR /avito-auto-unit/

COPY --from=builder /avito-auto-unit/shortener shortener
COPY config/ config/

EXPOSE 8080

ENTRYPOINT ["/avito-auto-unit/shortener"]
```
Первый слой образа компилирует бинарник приложения. Помимо этого также устанавливает переменные  окружениея для 
корректной сборки go-модуля приложения. На нижнем слое на основе пустого образа копируется бинарник и файлы 
конфигурации, пробрасывается наружу порт API сервера и задается команда, исполняемая при запуске контейнера.

<a name="cmd"></a>
# cmd
Главная директория проекта.

<a name="maingo"></a>
### main.go

<a name="main"></a>
#### func main
```go
func main()
```
Инициализирует все части проекта, необходимые для работы сервера: [сервис](#service), [редис-репозиторий](#redis), 
[хэндлер](#handler), роутер, а также читает [конфигурацию из конфига](#configyml), запускает сервер.

<a name="domain"></a>
# domain

<a name="domaingo"></a>
### domain.go
Файл с главными сущностями приложения, в нашем случае одной сущностью.

<a name="redirect"></a>
#### type Redirect
```go
type Redirect struct {
	Key       string `json:"key" msgpack:"key" valid:"-"`
	URL       string `json:"url" msgpack:"url" valid:"requrl"`
	CreatedAt int64  `json:"created_at" msgpack:"created_at" valid:"-"`
}
```
Главная сущность приложения - редирект (перенаправление). Структура имеет теги для работы с [JSON](#json), 
[msgpack](msgpack) и возможностью валидации ее полей.

<a name="repo"></a>
# repo
Директория для работы с хранилищами данных.

<a name="repogo"></a>
### repo.go
Файл, определяющий общие интерфейсы для работы баз данных с конкретными сущностями.

<a name="redirectrepo"></a>
#### type RedirectRepo
```go
type RedirectRepo interface {
	Load(key string) (*domain.Redirect, error)
	Store(redirect *domain.Redirect) error
}
```
Данный интерфейс описывает функции, которые должен имплементировать каждый вид хранилища для взаимодействия с сущностью
типа [domain.Redirect](#redirect): функция сохранения структуры типа ```domain.Redirect``` в базе данных и ее 
извлечения.

<a name="redis"></a>
## redis

<a name="redisgo"></a>
### repo.go
Файл, описывающий операции для работы redis с типом [domain.Redirect](#redirect).

<a name="redisclient"></a>
#### type redisRepo
```go
type redisRepo struct {
	client *redis.Client
}
```
Тип является алиасом, инкапсулирующим логику работы с клиентом redis, предоставляемого библиотекой.

<a name="newrepo"></a>
#### func NewRepo
```go
func NewRepo (cfg *config.RedisConfig) (*redisRepo, error) 
``` 
Публичная функция инициализации редис-репозитория из файла конфигурации.

<a name="newclient"></a>
#### func newClient
```go
func newClient(addr, port, password string) (*redis.Client, error)
```
Приватная функция инициализации редис-клиента, инкапсулирующая логику создания соединения с redis. Также пингует 
полученное соединение. 

<a name="redisload"></a>
#### func (r *redisRepo) Load
```go
func (r *redisRepo) Load(key string) (*domain.Redirect, error)
```
Метод загрузки структуры типа [domain.Redirect](#redirect) из redis, реализующий интерфейс 
[RedirectRepo](#redirectrepo). Данный метод под капотом использует редис-команду ```HGETALL```, которая возвращает хэш,
хранимый по ключу.

<a name="redisstore"></a>
#### func (r *redisRepo) Store
```go
func (r *redisRepo) Store(redirect *domain.Redirect) error
```
Метод сохранения структуры типа [domain.Redirect](#redirect) в redis, реализующий интерфейс 
[RedirectRepo](#redirectrepo). Сохранение в redis происходит с помощью редис-команды ```HSET```, сохраняющей хэшмап по
выбранному ключу и перезаписывающей значение в случае существования ключа.

<a name="service"></a>
# service 
Директория с бизнес-логикой приложения.

<a name="errors"></a>
### errors.go
Файл с ошибками, возникающими при работе со структурой типа [domain.Redirect](#redirect): неверная структура редиректа 
и редирект не найден.

<a name="servicego"></a>
### service.go
Бизнес-логика приложения.

<a name="redirectservice"></a>
#### type RedirectService
```go
type RedirectService interface {
	Load(key string) (*domain.Redirect, error)
	Store(redirect *domain.Redirect) error
}
```
Данный сервис соединяет бизнес-правила приложения и логику работы с хранилищами.

<a name="redirectservicest"></a>
#### type redirectService
```go
type redirectService struct {
	redirectRepo repo.RedirectRepo
}
```
Структура для имплементации интерфейса [RedirectService](#redirectservice). Т.к. мы работаем только с базой данных, 
структура включает в себя только интерфейс [repo.RedirectRepo](#redirectrepo).

<a name="newservice"></a>
#### func NewRedirectService
```go
func NewRedirectService(redirectRepo repo.RedirectRepo) RedirectService
```
Функция инициализации сервиса по работе с сущностью редиректа.

<a name="serviceload"></a>
#### func (r *redirectService) Load
```go
func (r *redirectService) Load(key string) (*domain.Redirect, error)
```
Метод загрузки структуры типа [domain.Redirect](#redirect), реализующий интерфейс [RedirectService](#redirectservice). 
Просто пробрасывает метод [repo.RedirectRepo.Load()](#redirectrepo).

<a name="servicestore"></a>
#### func (r *redirectService) Store
```go
func (r *redirectService) Store(redirect *domain.Redirect) error
```
Метод сохранения структуры типа [domain.Redirect](#redirect), реализующий интерфейс 
[RedirectService](#redirectservice). Метод валидирует полученное значение на соответствие URL и генерирует случайный
ключ, если он не был передан.

<a name="serialization"></>
# serialization
Пакет для описания логики, по которой сериализуется и десериализуется структуры определенного типа.

<a name="serializer"></a>
### serializer.go
Файл, определяющий общие интерфейсы для сериализации/десериализации конкретной бизнес-сущности.

<a name="redirectserializer"></a>
#### type RedirectSerializer
```go
type RedirectSerializer interface {
	Decode(data []byte) (*domain.Redirect, error)
	Encode(value *domain.Redirect) ([]byte, error)
}
```
Данный интерфейс описывает функции, которые должен имплементировать каждый формат сериализации сущности типа 
[domain.Redirect](#redirect).

<a name="json"></a>
### json

<a name="jsonserializer"></a>
### serializer.go
Файл, описывающий операции работы формата ```JSON``` с типом [domain.Redirect](#redirect).

<a name="jsonredirect"></a>
#### type Redirect
```go
type Redirect struct{}
```
Структура, имплементирующая интерфейс [RedirectSerializer](#redirectserializer).

<a name="jsondecode"></a>
#### func (r *Redirect) Decode
```go
func (r *Redirect) Decode(d []byte) (*domain.Redirect, error)
```
Метод декодирования структуры типа [domain.Redirect](#redirect) из формата ```JSON```. Под капотом использует обычный 
```json.Unmarshal```.

<a name="jsonencode"></a>
#### func (r *Redirect) Encode
```go
func (r *Redirect) Encode(v *domain.Redirect) ([]byte, error)
```
Метод кодирования структуры типа [domain.Redirect](#redirect) в формат ```JSON```. Использует ```json.Marshal```.

<a name="msgpack"></a>
### msgpack

<a name="msgpackserializer"></a>
### serializer.go
Файл, описывающий операции работы формата ```msgpack``` с типом [domain.Redirect](#redirect).

<a name="msgpackredirect"></a>
#### type Redirect
```go
type Redirect struct{}
```
Структура, имплементирующая интерфейс [RedirectSerializer](#redirectserializer).

<a name="msgpackdecode"></a>
#### func (r *Redirect) Decode
```go
func (r *Redirect) Decode(d []byte) (*domain.Redirect, error)
```
Метод декодирования структуры типа [domain.Redirect](#redirect) из формата ```msgpack```. Под капотом использует 
```msgpack.Unmarshal```.

<a name="msgpackencode"></a>
#### func (r *Redirect) Encode
```go
func (r *Redirect) Encode(v *domain.Redirect) ([]byte, error)
```
Метод кодирования структуры типа [domain.Redirect](#redirect) в формат ```msgpack```.

<a name="transport"></a>
# transport

<a name="transportgo"></a>
### transport.go
Файл описывающий работу приложения на сетевом уровне.

<a name="redirecthandler"></a>
#### type RedirectHandler
```go
type RedirectHandler interface {
	LoadRedirect(http.ResponseWriter, *http.Request)
	StoreRedirect(http.ResponseWriter, *http.Request)
}
```
Данный интерфейс содержит функции для работы со структурой типа [domain.Redirect](#redirect) на уровне ```http```.

<a name="handler"></a>
#### type handler
```go
type handler struct {
	redirectService service.RedirectService
}
```
Структура, соединяющая уровень ```http``` и [service](#service) нашего приложения. Реализует интерфейс 
[RedirectHandler](#redirecthandler).

<a name="newhandler"></a>
#### func NewHandler
```go
func NewHandler(redirectService service.RedirectService) *handler
```
Функция инициализирует и возвращает указатель на структуру типа [handler](#handler).

<a name="setupresponse"></a>
#### func setupResponse
```go
func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int)
```
Функция установки заголовка в ```http-ответе``` в зависимости от переданного ```Content-Type```. Также в ответе
устанавливается ```http-код``` статуса проведенной операции.

<a name="handlerserializer"></a>
#### func (h *handler) serializer
```go
func (h *handler) serializer(contentType string) serialization.RedirectSerializer
```
В зависимости от переданного ```Content-Type``` метод возвращает структуру определенного типа, реализующую интерфейс 
[serialization.RedirectSerializer](#serializer). 

<a name="transportload"></a>
#### func (h *handler) LoadRedirect
```go
func (h *handler) LoadRedirect(w http.ResponseWriter, r *http.Request)
```
Данный метод возвращает пользователю ```http-ответ```, в теле которого лежит структура типа 
[domain.Redirect](#redirect). Под капотом с помощью функции ```URLParam()``` из пакета 
[chi](https://github.com/go-chi/chi) из тела ```GET-запроса``` достается ключ, по которому нужно найти значение в базе
данных. Далее вызывается метод [service.RedirectService.Load()](#serviceload). В случае успешного нахождения структуры
происходит перенаправление на адрес полного URL, хранящегося по указанному ключу.

<a name="transportstore"></a>
#### func (h *handler) StoreRedirect
```go
func (h *handler) StoreRedirect(w http.ResponseWriter, r *http.Request)
```
Данный метод сохраняет структуру типа [domain.Redirect](#redirect) в хранилище данных. Поля для этой структуры 
достаются из тела ```POST-запроса``` пользователя. Декодирование проводится с помощью метода 
[serialization.RedirectSerializer.Decode()](#redirectserializer). Какой конкретно сериалайзер использовать определяется
на основании полученного в запросе ```Content-Type``` с помощью метода [serializer()](#handlerserializer). После этого 
структура сохраняется в хранилище, используя метод [service.RedirectService.Store()](#servicestore). В случае успешного
сохранения структура кодируется методом [serialization.RedirectSerializer.Encode()](#redirectserializer) и, если не 
было ошибок, отправляется пользователю в теле ответа методом [setupResponse()](#setupresponse) с кодом 201.

<a name="config"></a>
# config
Пакет, содержащий конфигурационные файлы приложения.

<a name="configyml"></a>
### config.yml
```yaml
redis:
  addr: redis
  port: 6379

server:
  addr: server
  port: 8080
```
Файл, содержащий общую конфигурацию приложения.

<a name="configgo"></a>
### config.go
Файл, читающий конфиг приложения.

<a name="configst"></a>
#### type Config
```go
type Config struct {
	Redis  *RedisConfig  `yaml:"redis"`
	Server *ServerConfig `yaml:"server"`
}
```
Структура, заключающая в себе общую конфигурацию приложения.

<a name="redisconfig"></a>
#### type RedisConfig
```go
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}
```
Структура для конфигурации redis.

<a name="serverconfig"></a>
#### type ServerConfig
```go
type ServerConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}
```
Структура для конфигурации API сервера.

<a name="readconfig"></a>
#### func Read
```go
func Read(path string) (*Config, error)
```
Функция чтения [файла общей конфигурации приложения](#configyml).

<a name="mocks"></a>
# mocks
Пакет с файлами моков, автоматически сгенерированными библиотекой [mockery](https://github.com/mockery/mockery). Данные
файлы необходимы для тестов.
