# О проекте
___
Проект для работы с NATS Streaming.

## Запуск
* Переходим в nats-streaming-server и запускаем сервер:

    ```
    go run nats-streaming-server.go
    ```
* Запускаем postgres:

    ```
    make init_db
    ```
* Запускаем сервер для приема сообщений. Переходим в my_service и делаем:

    ```
    go run cmd/main.go
    ```
* Запускаем сервер для отправки сообщений. Переходим в my_publisher и делаем:

    ```
   go run main.go
    ```

* Сервис доступен на 8000 порту по-умолчанию. Можно поменять в my_service/config/config.yml. 
Отправляет сохраненный из NATS Streaming заказ по uid. Модель заказа лежит в my_publisher/model.json.
    ```
   http://localhost:8000/1
    ```
## Contact

Ivan Konoplich - konoplich_i@mail.ru

Project Link: [https://github.com/IvanKonoplich/Wallet-Service.git](https://github.com/IvanKonoplich/Wallet-Service.git)

