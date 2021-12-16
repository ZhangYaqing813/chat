### Gateway
- Gateway主要实现与消息路由转发，根据自定义message 类型进行不同的数据处理。
    例如：login(用户登录)的message，怎根据msg.type 进行相应的处理。
- 实现后消息处理，比如一些消息的redis 缓存，聊天回话的存储等。
- Gateway 目前实现方式 采用 `switch case ` 方式，
  
- 