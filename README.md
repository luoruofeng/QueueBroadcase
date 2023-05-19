# QueueBroadcase

## 目标
相邻节点逐一分布在*队列*中，消息通过相邻的节点*逐一进行广播*。    
节点会把信息传递给相邻的节点，相邻的节点再把信息传递给它相邻的节点，以此类推。    
如果其中一个节点下线，他相邻的节点将无法获取信息，重写将其上线，则可以恢复通讯。   

    
## 节点参数
启动多个相邻的节点    
节点的id是参数   
节点的Address是loacalhost   
节点的port是8000+id    


## 例如
```
# 开启多个窗口并运行代码，参数为id，数字需要相邻
go run . 1
go run . 2
go run . 3
go run . 4
go run . 5
```
### 结构
    Node[Id:1]  -->  Node[Id:2] --> Node[Id:3] --> Node[Id:4] --> Node[Id:5]

