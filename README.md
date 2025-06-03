# 3xui-metrics-collector

Сервис для сбора метрик с панели управления 3xui и их экспорта в формате Prometheus.

## Описание

Этот сервис собирает метрики с панели управления 3xui и экспортирует их в формате Prometheus. Метрики включают:
- Количество клиентов
- Общий трафик
- Статус сервиса
- И другие метрики, доступные через API 3xui

## Требования

- Go 1.21 или выше
- Доступ к панели управления 3xui
- Prometheus (для сбора метрик)

## Сборка проекта

1. Клонируйте репозиторий:
```bash
git clone https://github.com/andrejmatveev/3xui-metrics-collector.git
cd 3xui-metrics-collector
```

2. Соберите проект:
```bash
go build -o 3xui-metrics-collector
```

## Конфигурация

Сервис использует следующие переменные окружения:

- `X_UI_HOST` - хост панели управления 3xui
- `X_UI_PORT` - порт панели управления 3xui
- `X_UI_BASEPATH` - базовый путь API (по умолчанию пустой)
- `X_UI_USERNAME` - имя пользователя для доступа к API
- `X_UI_PASSWORD` - пароль для доступа к API

## Настройка systemd сервиса

1. Создайте файл конфигурации `/etc/3xui-metrics-collector/config.env`:
```bash
X_UI_HOST=your_xui_host
X_UI_PORT=your_xui_port
X_UI_BASEPATH=your_base_path
X_UI_USERNAME=your_username
X_UI_PASSWORD=your_password
```

2. Создайте файл сервиса `/etc/systemd/system/3xui-metrics-collector.service`:
```ini
[Unit]
Description=3xui Metrics Collector
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/3xui-metrics-collector
ExecStart=/opt/3xui-metrics-collector/3xui-metrics-collector
EnvironmentFile=/etc/3xui-metrics-collector/config.env
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

3. Создайте необходимые директории и переместите файлы:
```bash
sudo mkdir -p /opt/3xui-metrics-collector
sudo mkdir -p /etc/3xui-metrics-collector
sudo cp 3xui-metrics-collector /opt/3xui-metrics-collector/
sudo cp config.env /etc/3xui-metrics-collector/
```

4. Установите правильные разрешения:
```bash
sudo chown -R root:root /opt/3xui-metrics-collector
sudo chown -R root:root /etc/3xui-metrics-collector
sudo chmod 755 /opt/3xui-metrics-collector/3xui-metrics-collector
sudo chmod 600 /etc/3xui-metrics-collector/config.env
```

5. Запустите сервис:
```bash
sudo systemctl daemon-reload
sudo systemctl enable 3xui-metrics-collector
sudo systemctl start 3xui-metrics-collector
```

## Управление сервисом

- Проверить статус:
```bash
sudo systemctl status 3xui-metrics-collector
```

- Остановить сервис:
```bash
sudo systemctl stop 3xui-metrics-collector
```

- Перезапустить сервис:
```bash
sudo systemctl restart 3xui-metrics-collector
```

- Посмотреть логи:
```bash
sudo journalctl -u 3xui-metrics-collector -f
```

## Метрики

Сервис экспортирует метрики в формате Prometheus на порту 2112. Доступные метрики:

- `xui_clients_total` - общее количество клиентов
- `xui_traffic_total_bytes` - общий трафик в байтах
- `xui_service_status` - статус сервиса (1 - активен, 0 - неактивен)

## Лицензия

MIT 