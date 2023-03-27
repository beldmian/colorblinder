# Colorblinder filterer
Сервис фильтрации mpeg-dash стримов и VoD контента для людей с особенностями зрения

## Usage

### Startup
Для запуска приложения доступно несколько вариантов

#### Binary
Собираем приложение через `make build`, затем запускаем бинарник `./main`,
получаем запущенный backend-сервис на порте, указанном в `config.yaml` файле

#### Docker
Запускаем приложение через `docker compose up`, получаем запущенные backend и frontend сервисы

#### Kubernetes
Разворачиввем в кластере kubernetes конфиги из директории `deploy/kubernetes` (например через команду `kubectl apply -f ./deploy/kubernetes`).
Получаем запущенные сервисы бекэнда и мониторинга с предустановленым дашбордом в grafana.

### API scheme

- `POST /create_filter`
   - Request
   ```json
   {
      "rgba_overlay": [128, 0, 128, 50],
      "start_second": 0,
      "is_photosensitive": true
   }
   ```
   - Response
   ```json
   {
      "id": "b1642bb4-e7fd-42ac-93dd-9334b4b74e35"
   }
   ```
- `POST /start_stream`
   - Request
   ```json
   {
      "filter_id": "b1642bb4-e7fd-42ac-93dd-9334b4b74e35",
      "stream_url": "https://dash.akamaized.net/dash264/TestCasesHD/2b/qualcomm/1/MultiResMPEG2.mpd"
   }
   ```
   - Response
   ```json
   {
      "new_url": "/stream/b1642bb4-e7fd-42ac-93dd-9334b4b74e35/file.mpd"
   }
   ```

### Frontend

## User flow
1. Пользователь отправляет запрос с настройками и получает uuid настроенного плеера
2. Пользователь идет на uuid настроенного плеера с url-адресом mpd файла и получает адрес нового mpd файла
   1. В процессе создается процесс ffmpeg'а с заданными пользователем настройками, который делает новый mpd-стрим
3. Пользователь просматривает видео с заданными настройками
4. Через некоторое количество времени после прекращения запросов процесс ffmpeg'а и все его файлы удаляются

ffmpeg -i video.mp4 -vf photosensitivity=bypass=1,metadata=print:file=photosensitivity-analysis.txt -f null /dev/null