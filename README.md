# MyShortLink
 
MyShortLink — простой и удобный сокращатель ссылок, созданный с использованием Go. Если вам нужно сократить длинные URL для сообщений в социальных сетях или вы просто предпочитаете более управляемые ссылки, MyShortLink предоставляет бесперебойный опыт без необходимости регистрации. Тем не менее, если вы хотите отслеживать и управлять вашими сокращёнными ссылками, регистрация аккаунта доступна для хранения ваших ссылок в учетной записи пользователя.

Доступен по ссылке: **[myshl.ru](https://myshl.ru/)**

# Функции
- Быстрое сокращение ссылок: Сокращайте URL-адреса немедленно без необходимости регистрации – всего за несколько кликов.
- Учетные записи пользователей: Опциональная регистрация для лёгкого управления и отслеживания ваших сокращённых URL.

# Установка, эксплуатация и внесение изменений
## Процесс установки для Ubuntu

1. Подключение к VPS

```console
ssh your_user@your_vps_address
```

2. Установка Git (если еще не установлен)

```console
sudo apt update && sudo apt install git -y
```

3. Клонирование репозитория

```console
cd /path/to/your/folder
git clone https://github.com/paych3ck/my-short-link.git
```

4. Запуск проекта
```console
cd my-short-link
go build
./my-short-link
```

## Использование как сервис
1. Создание Unit-файла
```console
sudo nano /etc/systemd/system/my-short-link.service
```

2. Содержимое Unit-файла

```bash
[Unit]
Description=My Short Link Application
After=network.target

[Service]
User=<username>
WorkingDirectory=/path/to/your/folder/my-short-link
ExecStart=/path/to/your/folder/my-short-link/my-short-link
Restart=always
# Другие опции, например, для установки переменных окружения:
# Environment="VAR1=value1" "VAR2=value2"

[Install]
WantedBy=multi-user.target
```

- Description: краткое описание службы.
- After: указывает, когда служба должна быть запущена. В этом случае после запуска сети.
- User: пользователь, от имени которого будет запущена служба.
- WorkingDirectory: рабочая директория для службы.
- ExecStart: команда для запуска приложения.
- Restart: политика перезапуска приложения, always означает, что приложение будет перезапущено при любом его завершении.
- WantedBy: целевая группа, к которой будет привязана служба.

2. Запуск службы
```console
sudo systemctl start my-short-link.service
```

3. Автоматический запуск при загрузке системы
```console
sudo systemctl enable my-short-link.service
```

4. Проверка статуса службы
```console
sudo systemctl status my-short-link.service
```

5. Остановка службы
```console
sudo systemctl stop my-short-link.service
```

6. Отключение автозагрузки службы
```console
sudo systemctl disable my-short-link.service
```

7. Перезапуск службы для применения изменений после редактирования unit-файла
```console
sudo systemctl restart my-short-link.service
```

## Внесение изменений
1. Остановка службы
```console
sudo systemctl stop my-short-link.service
```

2. Переход в директорию проекта
```console
cd /root/my-short-link
```

3. Внесение изменений в код
```console
nano your_file.go
```

4. Пересборка приложения
```console
go build
```

5. Перезапуск службы
```console
sudo systemctl start my-short-link.service
```