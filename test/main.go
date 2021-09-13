package main

import (
	"go-seckill/test/pressuremaker"
	"sync"
)

// 这个包对已经写成的功能模块进行压力测试和功能测试
// 如果对err信息感兴趣的话, 可以单独写一个分析error信息的函数
func main() {
	// logger.Infoln("Register Users")
	// pressuremaker.RegisterUsers()

	// generate token and store in token.db

	// logger.Println("Start generate Token")
	// pressuremaker.InitSqlite()
	// pressuremaker.CreateToken()

	// logger.Println("concurrent generate Token")
	// pressuremaker.GetTokenListConcurrent()

	// add test Product
	// pressuremaker.AddProducts()

	// log.Println("Start test")
	test()
}

// concurrent run CreateOrder() method, then run statistics
func test() {
	tokenList := pressuremaker.GetTokenListFromSqlite()

	var w sync.WaitGroup
	// 时间统计channel
	timeStatistics := make(chan float64, pressuremaker.ConcurrentNum)
	errChan := make(chan error, pressuremaker.ConcurrentNum)

	start := 0
	end := start + pressuremaker.ConcurrentNum
	for i := start; i < end; i++ {
		o := pressuremaker.Order{
			Token:       tokenList[i],
			ProductID:   4,
			PurchaseNum: 1,
		}
		w.Add(1)
		// 会将所有的error发送给errChan这个channel, 方便之后统计 }
		// every time duration send to timeStatistics
		go o.CreateOrder(&w, timeStatistics, errChan)
	}

	w.Wait()
	// 关闭时间统计channel, 开始我们的计算!
	close(timeStatistics)
	// 关闭错误统计channel, 开始我们的计算!
	close(errChan)

	// 把统计到的时间节点放置到一个slice中, 写需要计算的函数方法
	timeStatisticsList := make([]float64, 0, pressuremaker.ConcurrentNum)
	for t := range timeStatistics {
		timeStatisticsList = append(timeStatisticsList, t)
	}
	// logger.Infoln(timeStatisticsList, len(timeStatisticsList))
	pressuremaker.PlayTimeStatisticsList(timeStatisticsList)

	// 把统计到的错误信息放置到一个slice中, 写出自己需要的函数方法
	errStatisticsList := make([]error, 0, pressuremaker.ConcurrentNum)
	for e := range errChan {
		errStatisticsList = append(errStatisticsList, e)
	}
	// 正式运行的情况下, 把下面的这个替换为你自己写的错误统计函数
	logger.Infoln("errChan info: ", errStatisticsList)
}
