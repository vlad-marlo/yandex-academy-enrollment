# Лавка
Проект представляет из себя REST API-сервис, который будет регистрировать курьеров, добавлять новые заказы и 
расчитывать рейтинг курьеров.

## Требования к сервису

В сервисе должны быть реализованы:

1) REST API сервиса;
2) расчет рейтинга курьеров;
3) rate limiter для сервиса;

## 1. REST API

В качестве базовой функциональности сервиса необходимо реализовать 7 базовых методов.

Для всех методов в случае корректного ответа ожидается ответ `HTTP 200 OK`.

### POST /couriers
Для загрузки списка курьеров в систему запланирован описанный ниже интерфейс.

Обработчик принимает на вход список в формате json с данными о курьерах и графиком их работы.

Курьеры работают только в заранее определенных районах, а также различаются по типу: пеший, велокурьер и 
курьер на автомобиле. От типа зависит объем заказов, которые перевозит курьер.
Районы задаются целыми положительными числами. График работы задается списком строк формата `HH:MM-HH:MM`.

### GET /couriers/{courier_id}

Возвращает информацию о курьере.

### GET /couriers

Возвращает информацию о всех курьерах.

У метода есть параметры `offset` и `limit`, чтобы обеспечить постраничную выдачу.
Если:
* `offset` или `limit` не передаются, по умолчанию нужно считать, что `offset = 0`, `limit = 1`;
* офферов по заданным `offset` и `limit` не найдено, нужно возвращать пустой список `couriers`.

### POST /orders

Принимает на вход список с данными о заказах в формате json. У заказа отображаются характеристики — вес, район, 
время доставки и цена.
Время доставки - строка в формате HH:MM-HH:MM, где HH - часы (от 0 до 23) и MM - минуты (от 0 до 59). Примеры: “09:00-11:00”, “12:00-23:00”, “00:00-23:59”.


### GET /orders/{order_id}

Возвращает информацию о заказе по его идентификатору, а также дополнительную информацию: вес заказа, район доставки, 
промежутки времени, в которые удобно принять заказ.

### GET /orders

Возвращает информацию о всех заказах, а также их дополнительную информацию: вес заказа, район доставки, промежутки времени, в которые удобно принять заказ.

У метода есть параметры `offset` и `limit`, чтобы обеспечить постраничную выдачу.
Если:
* `offset` или `limit` не передаются, по умолчанию нужно считать, что `offset = 0`, `limit = 1`;
* офферов по заданным `offset` и `limit` не найдено, нужно возвращать пустой список `orders`.

### POST /orders/complete

Принимает массив объектов, состоящий из трех полей: id курьера, id заказа и время выполнения заказа, после отмечает, что заказ выполнен.

Если заказ:
* не найден, был назначен на другого курьера или не назначен совсем — следует вернуть ошибку `HTTP 400 Bad Request`.
* выполнен успешно — следует выводить `HTTP 200 OK` и идентификатор завершенного заказа.

Обработчик должен быть идемпотентным.

## 2. Рейтинг курьеров

Команда сервиса решила начать учет заработной платы и рейтинго курьеров.
Для этого необходимо реализовать новый метод `GET /couriers/meta-info/{courier_id}`.

Параметры метода:
* `start_date` - дата начала отсчета рейтинга
* `end_date` - дата конца отсчета рейтинга.

Примером значения параметров может быть `2023-01-20`. В задании можно полагаться на то, что все заказы и даты для 
расчетов имеют одну и ту же фиксированную временную зону - UTC.

Метод должен возвращать заработанные курьером деньги за заказы и его рейтинг.

**Заработок рассчитывается по формуле:**

Заработок рассчитывается как сумма оплаты за каждый завершенный развоз в период с `start_date` (включая) до 
`end_date` (исключая):

`sum = ∑(cost * C)`

`C`  — коэффициент, зависящий от типа курьера:
* пеший — 2
* велокурьер — 3
* авто — 4

Если курьер не завершил ни одного развоза, то рассчитывать и возвращать заработок не нужно.

**Рейтинг рассчитывается по формуле:**

Рейтинг рассчитывается следующим образом:
((число всех выполненных заказов с `start_date` по `end_date`) / (Количество часов между `start_date` и `end_date`)) * C
C - коэффициент, зависящий от типа курьера:
* пеший = 3
* велокурьер = 2
* авто - 1

Если курьер не завершил ни одного развоза, то рассчитывать и возвращать рейтинг не нужно.

## 3. Rate limiter

Каждый большой сервис с API, открытым из интернета, должен ограничивать количество входящих запросов.
Для этого используется rate limiter. Вам нужно реализовать такое решение для разрабатываемого сервиса.

Для решения задачи можно написать собственную реализацию или использовать известное готовое решение.
Сервис должен ограничивать нагрузку в 10 RPS на каждый эндпоинт. Если допустимое количество запросов превышено, сервис должен отвечать кодом 429.
