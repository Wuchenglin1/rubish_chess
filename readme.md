# 这是一个非常辣鸡的象棋

# 基本功能没有实现

## 说在前面

***不能下象棋，只能看***::sob:

### 主要实现的思路

在本地只需要加载图片文件和数组位置，主要的逻辑运算放在服务端进行(技术太菜，没来得及实现，太多时间花在了无用功上，一个功能没实现)

:star:使用websocket建立长连接

:star:使用ebiten做渲染

:sob:服务端的逻辑没有写完，不能完全展示出来

:x:没有实现grpc

:x:没有分成一个一个微服务实现





### 项目结构

├─.idea
├─client	`客户端`
│  ├─chess	`这是象棋主文件包`
│  ├─file2bytes		`这是一个把图片转换成二进制的包`
│  │  └─file2byteslice-master	
│  │      └─cmd
│  │          └─file2byteslice
│  ├─img	`存放img图片//当然我用 embed 包把全部png图片和img图片内嵌在代码中`
│  ├─main	`项目的开端[打包的exe文件]`
│  └─tool	`一些工具`
└─server	`服务端实现`
    ├─api	
    ├─dao
    ├─img	
    ├─main
    ├─model
    ├─service
    └─tool