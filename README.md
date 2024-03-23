# DIRTY RANDY

input league -> scrape ESPN league schedules -> output calendar events

```
dirty-randy auto-draft  f1
dr ad f1
```

output comma-delimited `.csv` events file for import into:
- google calendar

![dirty randy, brother](/dr.png)

## TOODOOZERS

- build test(s) from success output for refactor
- refactor highly reused stuff
    - create `check(err)` type util to see if `err != nil`
- figure out lift for other league support

## DONEZERS

- format event data per google calendar specs
    - Subject - "any string"
    - Start date - 05/30/2020
    - Start time - 10:00 AM
    - Location - "any string"
- create .csv file in `/csv` dir

