# envus (env ultra simple)
File based, local environment capable property system.

Register your property key, and envus will find the value using three tiers of definition.

Tier 1: A configured toml / json file containing prop definitions
Tier 2: Local environment that is toggleable
Tier 3: A default value passed in while registering variable

```bash
go get github.com/jabgibson/envus
```

envus by default looks for a properties file named .envus, this is of course
configurable by explicitly setting it.

```toml
# brutus.toml

[[prop]]
key = "NAME"
value = "Brutus"

[[prop]]
key = "vice"
value = "Mountain Dew"
```

Using a toml file such as above defined boris.toml the following code will parse it ans provide its
data as property values


```go
package main

import "github.com/jabgibson/envus"
import "fmt"
import "os"

func init() {
	os.Setenv("COLOR", "blue") // setting environment variable since toml file won't include this value
}

func main() {
	var name string
	var color string
	var secret string
	var vice string
	
	envus.Register("NAME", &name, "Boris")
	envus.Register("COLOR", &color, "orange")
	envus.Register("secret", &secret, "The world is a vampire")
	envus.Register("VICE", &vice, "kombucha")
	envus.UseLocalEnvironment() // Will use local environment variables before resorting to default values [tier 2]
	envus.OverrideEnvusFilename("brutus.toml") // the default filename if .envus if you don't care to override
	envus.Load()
	
	fmt.Println(name) // will be 'Brutus'
	fmt.Println(color) // will be 'blue'
	fmt.Println(secret) // will be 'The world is a vampire'
	fmt.Println(vice) // will be 'Mountain Dew'
}
```