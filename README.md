# Распределенный вычислитель арифметических выражений

#### Загружаем проект с git

```sh
git clone https://github.com/anaskozyr/distributed-calculator
```


#### Переходим в каталог проекта

```sh
cd distributed-calculator
```

#### Делаем сборку проекта и запуск сервисов командой

```sh
sudo docker compose up --build
```

#### Запуск GUI 

```sh
firefox frontend/index.html
```

Либо если вы используете хром 

```sh
google-chrome frontend/index.html
```

##### Остановка контейнеров производится командой

```sh
sudo docker-compose down
```

### Принципиальная схема в виде текстовой диаграммы

![Схема проекта](schema.png)

В этой схеме:

  ####### Фронтенд (Front-end) обеспечивает пользовательский интерфейс для взаимодействия с системой.
  ####### Бэкенд (Back-end) включает в себя оркестратор и агентов, которые обрабатывают запросы и вычисления.
  ####### Оркестратор принимает запросы от фронтенда, управляет вычислениями и делегирует задачи агентам.
  ####### Агенты выполняют вычисления и отправляют результаты обратно оркестратору.

Огромное спасибо за Ваше тестирование!

С глубоким уважением, Анастасия. :) 
