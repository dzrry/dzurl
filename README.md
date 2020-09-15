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
UUID, но, по моему мнению, что для коротких URL'ов такие последовательности слишком длинные (36 символов). Поэтому в 
проекте используется пакет [rs/xid](https://github.com/rs/xid). Данный пакет генерирует строку символов, в 20
символов, причем такая строка весит на 4 байта меньше UUID.

- Сколько должны жить ссылки? - Я выбрал решение в безвременном хранении ссылок, т.к. сервис будет не очень популярен. 
Но в большом сервисе могло бы произойти захламление мертвыми и одноразовыми ссылками. Возможно, стоило ввести TTL для 
каждого URL или же предусмотреть возможность создания одноразовых ссылок.

- Нужно ли перезаписывать полный URL, если такой короткий URL уже существует. - Я выбрал решение перезаписи, т.к. 
возможность удаления ссылок и, соответственно, ключей не предусмотрена.

