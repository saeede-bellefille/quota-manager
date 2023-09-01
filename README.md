Quota Manager Project
====

Quota project contains an HTTP server which saves inputs from the user.

Running Project
---
Quota project needs a config file and redis to be run.

The config file should be available in `config.yaml` (in current directory)

Redis should be accessible in the address which is defined in `config.yaml`


Project Description
---
Each user has two limitations:
- Number of requests per minute (defined in config with `request_per_minute` key)
- Total size of data per month (defined in config with `size_per_month` key)

Each request is JSON in following format (which is sent to `/run` on HTTP)
```json
{
    "id": 100,
    "userId": 1,
    "data": "My Data"
}
```
`id` is used for uniqueness (each id is processed only once).

`data` is saved according to quota defined for `userId`.

Notes
---
- Data is *NOT* actually saved, it is only a mock (it just sleeps for one second and prints data).
- The project runs multiple workers (defined in `worker_pool_size`) which runs mock saving
- All data is queued in a worker queue (defined in `queue_size`), which means if that much data is being processed, new data waits until they are finished
- All quotas are checked in redis, which enables us to run multiple instances of project (with same redis)
