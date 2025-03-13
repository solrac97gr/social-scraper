# TGStats
Breef information about how to interact with the Stats API for get info for Channels in TG.

## Request Docs
Request URL: https://api.tgstat.ru/channels/stat
Request Method: GET
Parameters:
    - token: string
 	- channel: string
Possible API response:
    - If is a Channel:
        {
    "status": "ok",
    "response": {
        "id": 118,                           # Внутренний ID канала в TGStat
        "title": "РИА Новости",              # Название канала
        "username": "@rian_ru",              # Username канала
        "peer_type": "channel",              # Тип (канал/чат)
        "participants_count": 2048184,       # Количество подписчиков канала на момент запроса
        "avg_post_reach": 541540,            # Средний охват публикации
        "adv_post_reach_12h": 475712,        # Средний рекламный охват публикации за 12 часов
        "adv_post_reach_24h": 554476,        # Средний рекламный охват публикации за 24 часа
        "adv_post_reach_48h": 580952,        # Средний рекламный охват публикации за 48 часов
        "err_percent": 26.4,                 # Процент вовлеченности подписчиков (ERR %)
        "err24_percent": 25.2,               # Процент вовлеченности подписчиков в просмотр поста за первые 24 часа (ERR24 %)
        "er_percent": 11.11,                 # Коэффициент вовлеченности подписчиков во взаимодействия с постом (реакция, пересылка, комментарий)        
        "daily_reach": 35496444,             # Cуммарный дневной охват
        "ci_index": 8737.68,                 # Индекс цитирования (ИЦ)
        "mentions_count": 171477,            # Количество упоминаний канала в других каналах
        "forwards_count": 472536,            # Количество репостов в другие каналы
        "mentioning_channels_count": 18740   # Количество каналов, упоминающих данный канал
        "posts_count": 53500,                # Общее количество неудаленных публикаций в канале
    }
}
    - If is a Chat:

## Principal Problem we need to Approach
Limited amount of request per month

## Solution
- In the file tgstats/repository/mongo.go implement a repository that will store our consultings together with a expiration_time field that will be the actual time + one month that is the frecuency of this data is requested

- in the file tgstats/tsgstats.go we will execute the request first checking if we have the data inside of the database, then checking if its not expired.

- if Data is expired make a request if not return the stored Data

- Also I wanna expose another method in the repository of mongo that let me store the amount of request I have been made:
    ```go
    {
        "request_amount":100,
        "request_not_cached":67
        "date-period":"2025-03"
    }
    ```
    This data must be stored grouped in monthly periods 

- also in the folder tgstats/config/config.go we will implement a struct with a method named load config that will take the config file (tgstats/config/config.json) and load it in the following struct
```go
type TGStatsConfig struct{
    Token string `json:"token"`
}
```
wich is the token we need to use the TG Stats API

- continue with the file tgstats/tsgstats.go we will need a method that will return the following data that was extracted from the API and recive the Channel we need to extract the data as parameter.
    - avg_post_reach
    - er_percent
we can use this Go struct for return it
```go
type TGStatsResult struct {
    AvgPostReach float32
    ERPercent   float32
}
```