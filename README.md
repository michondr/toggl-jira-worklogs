# move your records from [Toggl](https://toggl.com/) to [Jira](https://www.atlassian.com/software/jira)

make your .env with all necessary tokens:
```shell
cp .env.dist .env
```

build with
```shell
go build .
```

run
```shell
./toggl-jira-worklogs -h
```

command runs by default with current day and fetches records from toggl from start to end of that day. 

all records which have some time tracked in that time window will be fetched.

if the record already exists in jira worklogs it will not be inserted

```text
$ ./toggl-jira-worklogs 
from  2024-08-26T00:00:00+02:00
until 2024-08-27T00:00:00+02:00

will process 10 toggl entries

Issue ID  Time	
FF-5457   15m	  duplicate of https://fc.atlassian.net/browse/FF-5457?focusedWorklogId=940
FF-6213   15m	  duplicate of https://fc.atlassian.net/browse/FF-6213?focusedWorklogId=938
FF-5026   30m	  duplicate of https://fc.atlassian.net/browse/FF-5026?focusedWorklogId=935
FF-6202   15m	  duplicate of https://fc.atlassian.net/browse/FF-6202?focusedWorklogId=934
FF-6188   45m	  duplicate of https://fc.atlassian.net/browse/FF-6188?focusedWorklogId=937
FF-6127   15m	  added: https://fc.atlassian.net/browse/FF-6127?focusedWorklogId=939
FF-6177   15m     added: https://fc.atlassian.net/browse/FF-6177?focusedWorklogId=931
FF-5085   1h 30m  added: https://fc.atlassian.net/browse/FF-5085?focusedWorklogId=930
FF-3123   4h 30m  added: https://fc.atlassian.net/browse/FF-3123?focusedWorklogId=943
FF-3123   45m	  added: https://fc.atlassian.net/browse/FF-3123?focusedWorklogId=942

```

### things to improve
* more adjustability?