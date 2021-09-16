> 这个包的用途是: 压力测试
>
> 单元测试在每个包的内部, 例如: 数据库模块位于: internal/mysql/shop_orm, 相关的单元测试文件也在同样的目录

方式: 通过模拟 10000 个并发连接, 向服务器请求数据, 通过统计每个连接从开始到结束经历的时间长度, 来测量系统的处理能力

使用方法:

```go
// 这个包对已经写成的功能模块进行压力测试和功能测试
// 如果对err信息感兴趣的话, 可以单独写一个分析error信息的函数
func main() {
  // 首先, 注册用户, pressuremaker包里面已经写好了批量注册测试用户的函数
  // 去除下方注释, 直接调用即可
	// logger.Infoln("Register Users")
	// pressuremaker.RegisterUsers()

	// generate token and store in token.db
  // 每个用户的token都是不同的, 将它们存在在sqlite3中(我最喜欢的数据库😁)
  // 去除下方注释, 直接调用即可
	// logger.Println("Start generate Token")
	// pressuremaker.InitSqlite()
	// pressuremaker.CreateToken()

  // 并发的生成token并存储在sqlite3, 和上方任选其一
  // 这个有些小问题, *sql.DB有时候会处理不了很大的并发连接, 直接报错
  // 或者我改变获取token的方式, 通过建立worker pool的方式来获取token
  // 建议使用上方的单线程生成用户token的方法, mysql加了索引后, 速度挺快的, 几秒钟的事情
	// logger.Println("concurrent generate Token")
	// pressuremaker.GetTokenListConcurrent()

  // 生成测试商品, 去除下方注释, 直接调用即可
	// add test Product
	// pressuremaker.AddProducts()

	// log.Println("Start test")
	test()
}
```
