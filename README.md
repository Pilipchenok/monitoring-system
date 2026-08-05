# Go Monitoring System

Распределенная клиент-серверная система на Go для сбора, хранения и визуализации системных метрик компьютера в реальном времени.

---

## Возможности

- Сбор реальных системных метрик (CPU, RAM, Disk, Network) с помощью библиотеки gopsutil.
- REST API для приема данных от агентов и отдачи истории для интерфейса.
- Надежное хранение исторических данных в PostgreSQL.
- Автоматическая инициализация и миграция структуры базы данных при первом запуске.
- Интерактивный Live-дашборд на чистом HTML/JS с использованием Chart.js (обновление данных без перезагрузки страницы).
- Полная изоляция сервисов внутри виртуальной сети Docker.
- Удобное управление жизненным циклом проекта через Makefile.

## Архитектура

- Agent (Фоновый процесс, собирающий метрики ОС с заданной частотой и отправляющий их по HTTP).
- Collector (HTTP-сервер, обрабатывающий входящие метрики и отдающий их для дашборда). Включает слои:
  - handler (Маршрутизация и валидация HTTP-запросов/ответов).
  - service (Бизнес-логика).
  - repository (Взаимодействие с базой данных).
- PostgreSQL (Реляционная база данных для хранения метрик).
- Frontend (Статичный HTML-файл, который браузер запрашивает у Коллектора. Использует Fetch API (Polling) для динамического обновления графиков).

## Как работает запрос

OS -> Agent (gopsutil) -> HTTP POST -> Collector -> PostgreSQL ->

-> Collector -> HTTP GET (Polling) -> Browser (Chart.js)

## Структура проекта

```text
monitoring-system/
├── cmd/
│   ├── agent/
│   │   └── main.go
│   └── collector/
│       └── main.go
├── configs/
│   ├── agent_config.yaml
│   └── collector_config.yaml
├── internal/
│   ├── agent/
│   ├── collector/
│   │   ├── handler/
│   │   ├── repository/
│   │   └── service/
│   └── model/
├── migrations/
│   └── 001_init.sql
├── static/
│   └── index.html
├── docker-compose.yml
├── Makefile
└── README.md
```

## Конфигурация

Система настраивается через YAML-файлы в папке configs/.

Пример `configs/agent_config.yaml`:

```YAML
collector_url: "http://collector:8080"
poll_interval: 2s
report_interval: 5s
```

Пример `configs/collector_config.yaml`:

```YAML
port: 8080
database_url: "postgres://user:pass@postgres:5434/monitoring?sslmode=disable"
```

## Работа с системой

1. Клонировать репозиторий

```bash
git clone https://github.com/Pilipchenok/monitoring-system.git
cd monitoring-system

```

2. Запустить всю инфраструктуру одной командой

```bash
make docker-up
```

3. Открыть Live-дашборд

Перейдите в браузере по адресу:

```text
http://localhost:8080/
```

Вы увидите графики потребления ресурсов, которые будут автоматически обновляться каждые 5 секунд.

4. Остановка системы

```bash
make docker-down
```

---

## Тестирование

Проект покрыт unit-тестами с использованием паттерна моков (Mocking) для изоляции слоев.

```bash
make test
make race
make vet
```

---

## Что было реализовано

- Полноценный End-to-End пайплайн работы с данными (от железа до графиков в браузере).
- Внедрение зависимостей (Dependency Injection) через интерфейсы на Go.
- Использование стандартного роутера `http.ServeMux` (Go 1.22+) для REST API.
- Настройка `Depends_on` и политик перезапуска (`restart: always`) для решения проблемы Race Condition при старте контейнеров в Docker.
- Механизм Database Initialization (Entrypoint скрипты) для PostgreSQL.
- Написание Vanilla JS скриптов для асинхронного обновления DOM дерева.

---

## Технологии

* **Backend:** Go (1.22+), `net/http`, `database/sql`, `shirou/gopsutil`
* **Database:** PostgreSQL
* **Frontend:** HTML, JavaScript (Fetch API), Chart.js
* **DevOps:** Docker, Docker Compose, Makefile, YAML

---

## Репозиторий

GitHub: https://github.com/Pilipchenok/monitoring-system.git
