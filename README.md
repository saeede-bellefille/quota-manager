Quota Manager Project
====

Quota project contains an HTTP server which saves inputs from user.

Running Project
---
Quota project needs a config file and redis to be run.

Config file should be available in `config.yaml` (in current directory)

Redis should be accessible in address which is defined in `config.yaml`


Project Description
---
Each user has two limitations:
- Number of requests per minute (defined in config with `request_per_minute` key)
- Total size of data per month (defined in config with `size_per_month` key)

Each request is JSON in following format
```json
{
    "id": 100,
    "userId": 1,
    "data": "My Data"
}
```
`id` is used for uniqeness (each id is proccessed only once).

`data` is saved according to quota defined for `userId`.

Notes
---
- Data is *NOT* actually saved, it is only a mock (just sleeps one second an prints data)
- Project runs multiple workers (defined in `worker_pool_size`) which runs mock saving
- All data is queued in a worker queue (defined in `queue_size`), which means if that much data is being proccessed, new data waits until they are finished
- All quotas are checked in redis, which enables us to run multiple instances of project (with same redis)
