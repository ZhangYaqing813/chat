## Server LoadConfig 方法设计

### 目标
```text
能够实现配置文件加载解析。
配置文件一共分为4个部分，
redis、mysql、server、client 
其中 redis、mysql、server 为同一个struct(serverConfig)，提供给server 使用
clientConfig 为单独的一个，给client 使用，主要是配置连接服务端使用。
```
### 实现
· 方法一 
```text
配置文件直接放在代码包中的独立文件中，配置好后，由相关方法一次加载解析完成后供服务使用
因此在代码编译之前需要把配置文件配置完成，每次修改配置文件都需要进行新的编译
```

· 方法二
```text
配置文件单独存放一个目录，在服务运行启动的之后指定配置文件的目录和配置文件名，这样需要
在主程序添加一个读取配置文路径的参数，配置文件修改后只需要重启服务即可
```

### 思路
```text
读取配置文件整体思路分为三个部分
```

· 第一步
```text
代码中添加相关的配置代码，加载配置文件。
主要是读取配置文件，然后根据换行将配置文件存放到切片中，便于后面处理
    
```
```go


```
· 第二步
```text
解析配置代码，解析加载后的配置文件
需要解析两个配置：
1、解析配置文件中的字段，根据ini文件特点读取到各个配置节点，进行保存
2、解析代码中与配置文件对应的嵌套结构体（结构体）获取，获取结构体中各个子结构体的信息，对应就是
    配置文件中的节点信息。
3、整体是以 配置文件的切片为循环，找到配置文件中的节点信息后，在进行找嵌套结构体中的子结构体，缓存起来，
    结束循环。
```
· 第三步 
```text
对应的字段赋值
1、先解析每一行的配置文件，拆分成key(字段名)和value(字段值)
2、根据上面获取到的子结构体的信息进行for 遍历，获取到子结构体中的字段信息，然后进行key对比，
    如果匹配上则跳出循环，同时缓存获取到的子结构体中的字段信息
3、根据获取到的子结构体中的字段新进行类型判断，然后进行赋值。
```

