
### ![](https://raw.github.com/nic0lae/JerryMouse/master/logo.png)

### Build API Server
```go
import "github.com/nic0lae/JerryMouse"

// Define Handler
func requestHandler(rh JsonHandler) Handle() {

}

// Run Server
apiServer = JerryMouse.ApiServers.HttpApi()
apiServer.Run(":8080", [JerryMouse.ApiServers.JsonHandler])

```
