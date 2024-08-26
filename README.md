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

if the record already exists in jira worklogs it will not be returned

### things to improve
* more adjustability?