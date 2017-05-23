### console.dir(selected_node)
#### what is this?
```
jm.node{
    children(array)
    data(Object)
    parent(Object) //parent info
    _data(Object)
    __proto__(Object)
}
```
#### inside data(Object)
```
irection:1
expanded:true
id:"c24ddde2b510ef39" // What ID?? PARENT ID!
index:1
isroot:false
```
- **parent = root --> isroot? true (parent: undefined)**
- nomal case: parent id + isroot? false

#### consoleで確認した後
##### selected_nodeをsocket.send()で送る-> NULLになってしまう?
###### Go lang Websocket use JSON
```golang
var JSON = Codec{jsonMarshal, jsonUnmarshal}
// JSON is a codec to send/receive JSON data in a frame from a WebSocket connection.

import "websocket"

type T struct {
    Msg string
    Count int
}
// receive JSON type T
var data T
websocket.JSON.Receive(ws, &data)
// send JSON type T
websocket.JSON.Send(ws, data)
```
##### client.go:
websocket.Conn

```go
func echo(conn *websocket.Conn) {
    for {
        m := msg{}
        err := conn.ReadJSON(&m)
        if err != nil {
            fmt.Println("Error reading json.", err)
        }
        fmt.Printf("Got message: %#v\n", m)
        if err = conn.WriteJSON(m); err != nil {
            fmt.Println(err)
        }
    }
}
```
###### HP docs:
```
func ReadJSON
func ReadJSON(c *Conn, v interface{}) error
ReadJSON is deprecated, use c.ReadJSON instead.

func WriteJSON
func WriteJSON(c *Conn, v interface{}) error
WriteJSON is deprecated, use c.WriteJSON instead.
```

Call the connection's WriteMessage and ReadMessage methods to send and receive messages as a **slice of bytes**. This snippet of code shows how to echo messages using these methods:

```javascript
for {
    messageType, p, err := conn.ReadMessage()
    if err != nil {
        return
    }
    if err = conn.WriteMessage(messageType, p); err != nil {
        return err
    }
}
```

In above snippet of code, p is a []byte and messageType is an int with value websocket.BinaryMessage or websocket.TextMessage.

### Convert to JSON Object:
Then, use the JavaScript built-in function `JSON.parse()` to convert the string into a JavaScript object:
```javascript
var obj = JSON.parse(text);
```


