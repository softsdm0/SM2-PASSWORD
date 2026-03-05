# session

## 前后端交互流程

### 交换密钥阶段

前端使用 sm2 公钥加密{'sm4Key': 'xxxxxxxxx','msg': '使用 sm4 加密的 ok.'}转为 string。 发送给后端/public/api/session/id

后端接收后使用 sm2 私钥解密，得到 sm4Key 和 msg，再用 sm4Key 解密 msg 看看内容是否为`ok.`如果不是就返回错误码，如果是就把 sm4Key 存到 redis 并且生成一个 sessionId 作为 redis 的 key`password_session_xxxxxxxxx`再把 sessionId 通过 sm4Key 加密返回给前端。

### 交换数据阶段

前端在每次请求前使用 sm4Key 加密 body 传给后端，并且在请求头带上`password_sessionid: xxxxxxx`

后端收到请求后使用读取请求头里面的`password_sessionid`去数据库里面查对应的 sm4Key，使用 sm4Key 解密 body
