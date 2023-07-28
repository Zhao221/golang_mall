# golang_mall
关于电子商城的项目，相关的技术栈go、gin、gorm、mysql、redis


**涉及到的模块**

**~.用户模块**

**~.商品模块**

**~.收藏夹模块**

**~.订单模块**

**~.购物车模块**

**~.地址模块**

**~.支付模块**

**~.秒杀专场**

**项目亮点**

1.用户登录使用jwt

2.用户密码使用hash加密，金额使用AES对称加密

3.使用Redis set数据结构实现好友共同关注功能

4.使用Redis bitmap数据结构实现用户签到功能

5.商品秒杀：使用hash存储数据，+分布式锁实现（还不完整待实现）

6.使用kafka消息队列给老客户推荐新上货物（待做）

7.支付功能逻辑较复杂：设计到订单，用户，老板，商品等多个模块


