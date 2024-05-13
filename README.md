Запуск только Citus
```bash
$ docker-compose up --scale citus-worker=2 -d citus-manager citus-master citus-worker
```

Запуск
```bash
$ docker-compose up --scale citus-worker=2 -d
```

Стоп
```bash
$ docker-compose stop
```

Пересборка сервиса
```bash
$ docker-compose up --build --no-deps app
```

