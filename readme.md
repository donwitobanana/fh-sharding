# FH Sharding

## Assumptions

- at application startup workerpools are created in the number set by PARTITIONS_NUMBER variable
- each workerpool is given `partition ID` which maps to a storage partition of the same ID
- number of initial workers is set by WORKERS_INIT_NUMBER variable
- 10 seconds before initial workers timeout set by WORKERS_TIMEOUT variable new workers are added
up to WORKERS_MAX_NUMBER, to make sure at least one worker is up between stopping old workers and reinstatiating them
- each batch of processed messages is saved into single file under the path of given format `<storage_dir>/<partition_id>/<worker_id>-<batch_number>` where `<storage_dir>` is directory pointed by STORAGE_PATH variable
- batch size is set by WORKERS_BATCH_SIZE variable

## Testing

to start service locally, first override env variables in `.env` file, then run:
```bash
make run
``` 

while service is up, run:
```bash
make test
```
it will run smoke test script resulting in service writing processed messages to the storage pointed by STORAGE_PATH variable.

To reset storage run: 
```bash
make reset-storage
```