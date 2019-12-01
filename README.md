## Docker image

* [acesso/dkron](http://hub.docker.com/r/acesso/dkron)
* [acesso/dkron:v2.0.0](http://hub.docker.com/r/acesso/dkron)
* [acesso/dkron:v2.0.0-acs.0](http://hub.docker.com/r/acesso/dkron)

## Plugin binary

* [GitHub Releases](https://github.com/acesso-io/dkron-executor-nats/releases)

## Example usage

```go
package main

import (
	"fmt"

	"github.com/distribworks/dkron/v2/proto"
	"gopkg.in/resty.v1"
)

func main() {
	job := proto.Job{
		Name:        "my-task",
		Displayname: "My Task",
		Timezone:    "UTC",
		Schedule:    "@every 10s",
		Owner:       "My Team",
		OwnerEmail:  "myteam@mycompany.com",
		Disabled:    false,
		Tags:        map[string]string{},
		Metadata:    map[string]string{},
		Concurrency: "allow",
		Executor:    "nats",
		ExecutorConfig: map[string]string{
			"subject": "dkron",
			"message": "this is my message",
		},
	}

	res, err := resty.R().SetBody(job).Post("http://localhost:8080/v1/jobs")

	fmt.Println(err)
	fmt.Println(res)
}
```
