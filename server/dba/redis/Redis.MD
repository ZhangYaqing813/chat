## Redis 缓存信息的操作和问题
```text
由来：
    起初使用redis 时，是为了缓存在线用户列表，是一个类似于key-value 模式的
    数据结构，key 为在线用户名，value 为在线用户的链接内存地址（*net.conn）
    但是发现 在序列化内存地址的时候出现问题，后改用 UserOlineIntoRedis
    结构体 依然有问题（前后折腾了两天的时间，尝试过hex.EncodeToString,未能成功可能是操作方式不对）
    因此放弃，存入redis，改用map[string]net.conn 在内存中维护。
    为了不放弃对redis的练习，准备将发送的信息存入redis。其问题依旧很多
    
目的：
    1、练习使用redis 缓存、读取数据
    2、熟悉和思考 数据存储方式对上下文的影响
    3、发现问题
```
### 操作
  ```text
 目的：
    主要是将server 端转发的信息进行redis 缓存
    
```
- 字段设计
 ```text
redis 缓存的key 的命名方式为 key:userAReFromUserB/userASendToUserB/userAUnreadMessage field:时间戳： value：message
    key:
        userAReFromUserB 	    该user 接收到的信息。
        userASendToUser B  	该user 发送的信息
        userAUnreadMessage 	该user 未读的信息

    field:
        信息中的时间戳，需要是否处理待定
```

- Redis 的写入
```text
 func (R *RedisOpt) AddMessage(key string,messages msg.Messages,fieldName string) (code int, err error)
    参数：
    key 对应在线用户的三个key（userAReFromUserB/userASendToUserB/userAUnreadMessage） 中的某一个
    message  server 端需要转发的信息
    fieldName 字段名 
key 的值使用调用者决定的，当server 查看完信息后转发时决定将信息存储于那个字段中
例如：
    server端转发时发现目标user 不在线，则将信息缓存至目标user的unreadMessage 中。
    当server 接收到目标user上线后取出对应信息，转发后进行删除。同时进行自身的userAReFromUserB 进行缓存。
```
- Redis 的读取
```text
func (R *RedisOpt) GetMessage(userName string,fieldName []string) (messages []string, err error)
Redis 读取分两种情况：
    情况1：
        只读取user 的userAUnreadMessage 信息，此时很简单。只需要取出返回即可
        后续处理：
            需要将 userAUnreadMessage 清空，然后转存至 userAReFromUserB 中
    
    情况2：
        相对很复杂，且会影响后续数据处理，单说读取，不说后续处理
        1）、用户获取自己发送的信息，此种情况简单。
        2）、用户获取自己已经接受的信息，此种情况简单。
        3）、用户获取自己完整的聊天记录，既对话式的，很复杂（指的后续数据处理）
    
```
- Redis Delete
```text
func (R *RedisOpt) Delete(key string,fieldName string) (code int, err error)
使用的条件有限，仅限于在将未读信息转发完成后 删除目标用户的未读信息
```

### 问题思考
- 思考一 
```text
Redis 是采用的是Hash 数据方式，结构简单明了，但是问题在于hash 是一个无序的，在数据取出时是否有必要将所有数据一次取出，
且取出的数据为无序的，因此需要进一步做排序处理，然后client端看到的才是一个有序的。且这种还需要分多种情况
    情况1：
        只看到自己发送的有序消息
    情况二：
        只看到自己接收的有序信息
    情况3：
        如何看到有序的对话信息，两种信息的组合
```
- 思考二
```text
后续的数据处理是否对整体转发的效率产生影响。如何进行消息的分页处理
```
- 思考三
