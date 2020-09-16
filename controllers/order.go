package controllers

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/valyala/fasthttp"
	"go_redis/jsonStruct"
	"go_redis/mysql"
	"go_redis/mysql/shop/goods"
	"go_redis/redis_config"
	"go_redis/utils"
	"log"
	"net/http"
	"strconv"
)

func errorHandle(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	http.Error(w, err.Error(), code)
}

//var cancelBuyLock sync.Mutex

// 处理用户要购买某种商品时, 提交的参数: userId, productId, productNum 的参数的处理呀
// 使用application/json的方式
func test(w http.ResponseWriter, r *http.Request) {
}

func Buy(ctx *fasthttp.RequestCtx) {
	//// 请求方法限定为post
	//if ctx.Request.Header.IsPost() == false {
	//	ctx.Response.Header.Set("Allow", fasthttp.MethodPost)
	//	ctx.Error("request method must be post", 405)
	//	return
	//}
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	errorHandle(w, errors.New("请求方法不合法!"), 405)
	//	return
	//}

	// 使用了easyjson, 据说可以提高marshal, unmarshal的效率
	buyReqPointer := new(jsonStruct.BuyReq)
	err := buyReqPointer.UnmarshalJSON(ctx.PostBody())
	//err := json.Unmarshal(ctx.PostBody(), buyReqPointer)
	if err != nil {
		log.Printf("decode buy request error: %v", err)
		utils.ResponseWithJson(ctx, fasthttp.StatusInternalServerError, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "服务器内部错误: 无法解析客户端发送的body",
			Data: nil,
		})
		//ctx.Error("decode json body error", 500)
		return
	}

	// 一些数据校验部分, 校验用户id, productId, productNum
	u := new(User)
	u.userID = buyReqPointer.UserId
	// 判断productId和productNum是否合法
	ok, err := u.CanBuyIt(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
	if err != nil {
		log.Printf("user: %s CanBuyIt error: %s\n", u.userID, err.Error())
		//content, err := c.MarshalJSON()
		//content, err := jsonStruct.CommonResp(c)
		//if err != nil {
		//	log.Printf("%v\n", err)
		//	_ = utils.ResponseWithJson(ctx, fasthttp.StatusInternalServerError, jsonStruct.CommonResponse{
		//		Code: 8500,
		//		Msg:  "服务器内部错误: struct > []byte 时出现错误",
		//		Data: nil,
		//	})
		//	return
		//}
		utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
			Code: 8005,
			Msg:  "您购买的商品数量已达到上限或者缺货;" + err.Error(),
			Data: nil,
		})
		return
	}
	if ok {
		// 生成订单信息
		orderNum, err := u.orderGenerator(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
		if err != nil {
			utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
				Code: 8002,
				Msg:  "生成订单过程中出现错误:" + err.Error(),
				Data: nil,
			})
			//content, err := c.MarshalJSON()
			////content, err := jsonStruct.CommonResp(c)
			//if err != nil {
			//	log.Println(err)
			//	ctx.Error("store num is not enough", 500)
			//	return
			//	//errorHandle(w, errors.New(err.Error()), 500)
			//}
			//ctx.SetContentType("application/json")
			//ctx.SetBody(content)
			////w.Header().Set("Content-Type", "application/json")
			////w.Write(content)
			return
		}

		// 给用户的已经购买的商品hash表里面的值添加数量
		err = u.Bought(buyReqPointer.ProductId, buyReqPointer.PurchaseNum)
		if err != nil {
			utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
				Code: 8004,
				Msg:  "给用户的已经购买的商品hash表单productId添加数量时发生错误;" + err.Error(),
				Data: nil,
			})
			//content, err := c.MarshalJSON()
			////content, err := jsonStruct.CommonResp(c)
			//if err != nil {
			//	//errorHandle(w, errors.New(err.Error()), 500)
			//	log.Println()
			//	ctx.Error("add bought list error", 500)
			//	return
			//}
			//ctx.SetContentType("application/json")
			//ctx.SetBody(content)
			////w.Header().Set("Content-Type", "application/json")
			////w.Write(content)
			return
		}

		//w.Header().Set("application/json", "json")
		utils.ResponseWithJson(ctx, fasthttp.StatusOK, jsonStruct.CommonResponse{
			Code: 8001,
			Msg:  "操作成功",
			Data: jsonStruct.OrderResponse{
				UserId:      buyReqPointer.UserId,
				PurchaseNum: buyReqPointer.PurchaseNum,
				ProductId:   buyReqPointer.ProductId,
				OrderNum:    orderNum,
			},
		})
		//content, err := c.MarshalJSON()
		////content, err := jsonStruct.CommonResp(c)
		//if err != nil {
		//	log.Println(err)
		//	ctx.Error("json marshal error", 500)
		//	return
		//	//errorHandle(w, errors.New(err.Error()), 500)
		//}
		//ctx.SetContentType("application/json")
		////w.Header().Set("Content-Type", "application/json")
		////w.Write(content)
		//ctx.SetBody(content)
		return
	}
}

// redis收到后台的请求, 用户取消了订单, 需要用到的参数有: userId, productId, purchaseNum,  redis直接操作用户的: user:[userId]:bought 里面key为productId的, 赋值为0
// 这个接口必须由后台调用, 因为我没有做数据校验
func CancelBuy(ctx *fasthttp.RequestCtx) {
	//if ctx.Request.Header.IsPost() == false {
	//	ctx.Request.Header.Set("Allow", http.MethodPost)
	//	ctx.Error("request method is not supported", 405)
	//	return
	//}
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	errorHandle(w, errors.New("请求方式不合法!"), 405)
	//	return
	//}

	// 解析: /cancelBuy接口传过来的四个参数, userId, productId, purchaseNum, orderId
	cancelBuyReqPointer := new(jsonStruct.CancelBuyReq)
	err := json.Unmarshal(ctx.Request.Body(), cancelBuyReqPointer)
	if err != nil {
		log.Printf("%v\n", err)
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "解析body到json格式时出现错误",
			Data: nil,
		})
		//ctx.Error("decode request body error", 500)
		return
	}

	//cancelBuyReqPointer, err := decodeCancelBuyReq(r.Body)
	//if err!=nil {
	//	errorHandle(w, errors.New("reqBody解析到struct时出错!"), 500)
	//	return
	//}
	u := new(User)
	u.userID = cancelBuyReqPointer.UserId
	err = u.CancelBuy(cancelBuyReqPointer.OrderNum)
	if err != nil {
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8006,
			Msg:  fmt.Sprintf("用户: %s 取消订单: %s 时出现错误", cancelBuyReqPointer.UserId, cancelBuyReqPointer.OrderNum),
			Data: nil,
		})
		//c := jsonStruct.CommonResponse{
		//	Code: 8006,
		//	Msg:  "取消订单时失败!",
		//	Data: nil,
		//}
		//content, err := jsonStruct.CommonResp(c)
		//if err != nil {
		//	log.Println("encode resp body to []byte error", err)
		//	ctx.Error("encode resp body to []byte error", 500)
		//	return
		//}
		//ctx.SetContentType("application/json")
		//ctx.SetBody(content)
		////w.Header().Set("Content-Type", "application/json")
		////w.Write(content)
		return
	}
	utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
		Code: 8007,
		Msg:  fmt.Sprintf("用户: %s 取消订单: %s 成功", cancelBuyReqPointer.UserId, cancelBuyReqPointer.OrderNum),
		Data: nil,
	})
	return
	//content, err := jsonStruct.CommonResp(c)
	//if err != nil {
	//	log.Println(err)
	//	ctx.Error("encode resp body to []byte error", 500)
	//	return
	//	//errorHandle(w, errors.New(err.Error()), 500)
	//}
	//ctx.SetContentType("application/json")
	//ctx.SetBody(content)
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(content)
}

// 调用这个函数, 立刻同步
// (redis中存在的商品(一般情况下, 这个时候mysql中也是存在对应的产品的), redis中的数据同步到mysql), 将redis中已变更的商品数据, 同步到mysql中
// 用途: 更新redis中的商品数据到mysql中
func SyncGoodsFromRedis2Mysql(ctx *fasthttp.RequestCtx) {
	redisconn := redis_config.Pool.Get()
	defer redisconn.Close()
	// 首先, 将redis中存在的商品信息同步到mysql中
	reply, err := redis.Strings(redisconn.Do("keys", "store:*"))
	if err != nil {
		log.Println(err)
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "更新redis中现有的商品信息到mysql中出现错误",
			Data: nil,
		})
		//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
		return
	}
	type Goods struct {
		ProductName string `redis_config:"productName"`
		ProductId   int    `redis_config:"productId"`
		StoreNum    int    `redis_config:"storeNum"`
	}
	goodsListRedis := make([]*Goods, 0)
	for _, v := range reply {
		log.Println(v)
		goodsMap, err := redis.Values(redisconn.Do("hgetall", v))
		if err != nil {
			log.Println(err)
			utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
				Code: 8500,
				Msg:  "获取hmap中的键值对时出现了错误",
				Data: nil,
			})
			//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
			return
		}
		//log.Println(goodsMap)
		g := new(Goods)
		err = redis.ScanStruct(goodsMap, g)
		if err != nil {
			log.Println("redis_config scanStruct error: ", err)
			utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
				Code: 8500,
				Msg:  "redis_config scanStruct error",
				Data: nil,
			})
			//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
			return
		}
		log.Println(g)
		goodsListRedis = append(goodsListRedis, g)
	}
	// 开始一个mysql事务
	tx, err := mysql.Conn.Begin()
	if err != nil {
		log.Println(err)
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "",
			Data: nil,
		})
		//ctx.Error(err.Error(), 500)
		return
	}
	// 这里必须使用事务, 不能这么一条一条的搞
	for _, v := range goodsListRedis {
		_, err := tx.Exec("update goods set product_name=?, inventory=? where product_id=?", v.ProductName, v.StoreNum, v.ProductId)
		if err != nil {
			err1 := tx.Rollback()
			if err1 != nil {
				log.Println(err)
				return
			}
			log.Println(err)
			utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
				Code: 8500,
				Msg:  "",
				Data: nil,
			})
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "提交mysql事务时出现错误",
			Data: nil,
		})
		//ctx.Error(err.Error(), 500)
		return
	}
	utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
		Code: 8001,
		Msg:  "同步redis信息到mysql成功",
		Data: nil,
	})
	return
	//respJson, err := jsonStruct.CommonResp(jsonStruct.CommonResponse{
	//	Code: 8001,
	//	Msg:  "处理成功",
	//	Data: nil,
	//})
	//if err != nil {
	//	errLog(ctx, err, err.Error(), 500)
	//	return
	//}
	//ctx.Response.SetStatusCode(200)
	//ctx.Response.SetBody(respJson)
	//ctx.Response.Header.Set("Content-Type", "application/json")
}

// (mysql中存在 && redis中不存在)的商品数据到redis, 这个接口的用处是: mysql中新添加的商品数据, 需要同步到redis中, 同时保证redis中已存在的商品数据不变
// 用途: Mysql中添加了新的商品数据,把它同步到redis中
func SyncGoodsFromMysql2Redis(ctx *fasthttp.RequestCtx) {
	redisconn := redis_config.Pool.Get()
	defer redisconn.Close()
	// 在现有的MySQL表格中找到所有的商品数据, 比对redis的productList, 如果发现有商品不存在于redis中, 就把它添加进去
	storeList, err := redis.Strings(redisconn.Do("keys", "store:*"))
	if err != nil {
		log.Println(err)
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "获取redis中已经存在的商品信息出现错误",
			Data: nil,
		})
		//ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	storeIDlist := make([]string, 0, 128) // 分离redis中商品的ID出来, 到单独的store id list
	for _, v := range storeList {
		storeIDlist = append(storeIDlist, v[6:])
	}
	log.Println(storeIDlist) // redis中存在的商品信息
	goodsList, err := goods.QueryGoods()
	if err != nil {
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	for _, v := range goodsList {
		_, ok := utils.FindElement(storeIDlist, strconv.Itoa(v.ProductId))
		if !ok {
			// 给redis中添加相关商品数据
			err = redisconn.Send("hmset", "store:"+strconv.Itoa(v.ProductId), "productName", v.ProductName, "productId", v.ProductId, "storeNum", v.Inventory)
			if err != nil {
				log.Printf("%+v创建hash `store:%d`失败", err, v.ProductId)
				// 这里有风险, 万一给redis添加信息的时候出现错误, 那就凉凉
				utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
					Code: 8500,
					Msg:  "mysql to redis_config error",
					Data: nil,
				})
				//ctx.Error("给redis添加更新的产品数据出现错误", 500)
				return
			}
		}
	}
	utils.ResponseWithJson(ctx, 200, jsonStruct.CommonResponse{
		Code: 8001,
		Msg:  "将mysql中新添加的数据缓存到redis中成功",
		Data: nil,
	})
	return
	//respJson, err := jsonStruct.CommonResp(jsonStruct.CommonResponse{
	//	Code: 8001,
	//	Msg:  "处理成功",
	//	Data: nil,
	//})
	//if err != nil {
	//	errLog(ctx, err, err.Error(), 500)
	//	return
	//}
	//ctx.Response.SetStatusCode(200)
	//ctx.Response.SetBody(respJson)
	//ctx.Response.Header.Set("Content-Type", "application/json")
}

// 展示商品清单
func GoodsList(ctx *fasthttp.RequestCtx) {
	redisconn := redis_config.Pool.Get()
	defer redisconn.Close()

	reply, err := redis.Strings(redisconn.Do("keys", "store:*"))
	if err != nil {
		log.Println(err)
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "从redis中获取商品信息失败",
			Data: nil,
		})
		//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
		return
	}
	type good struct {
		ProductName string `redis_config:"productName"`
		ProductId   int    `redis_config:"productId"`
		StoreNum    int    `redis_config:"storeNum"`
	}
	goodsList := make([]*good, 0)
	for _, v := range reply {
		log.Println(v)
		goodsMap, err := redis.Values(redisconn.Do("hgetall", v))
		if err != nil {
			log.Println(err)
			utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
				Code: 8500,
				Msg:  "获取商品key: value时出现错误",
				Data: nil,
			})
			//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
			return
		}
		//log.Println(goodsMap)
		g := new(good)
		err = redis.ScanStruct(goodsMap, g)
		if err != nil {
			log.Println("redis_config scanStruct error: ", err)
			utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
				Code: 8500,
				Msg:  "redis_config scanStruct error",
				Data: nil,
			})
			//ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
			return
		}
		log.Println(g)
		goodsList = append(goodsList, g)
	}
	utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
		Code: 8001,
		Msg:  "获取商品清单成功",
		Data: goodsList,
	})
	return
	//response := jsonStruct.CommonResponse{
	//	Code: 8001,
	//	Msg:  "success",
	//	Data: goodsList,
	//}
	//err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response)
	//if err != nil {
	//	log.Println(err)
	//	ctx.Error("内部处理错误", fasthttp.StatusInternalServerError)
	//	return
	//}
	//ctx.Response.Header.Set("Content-Type", "application/json")
}

// 更新商品限制计划
// 例如, 在更新MySQL的限制购买条件后, 若要将商品购买限制同步到app中, 只需要调用goodsLimit这个接口就可以
func SyncGoodsLimit(ctx *fasthttp.RequestCtx) {
	// 加载limit限制计划
	err := LoadLimit()
	if err != nil {
		log.Println(err)
		utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
			Code: 8500,
			Msg:  "加载mysql中限制购买的数据到全局变量purchaseLimit时出现错误",
			Data: nil,
		})
		return
		//ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
	utils.ResponseWithJson(ctx, 500, jsonStruct.CommonResponse{
		Code: 8001,
		Msg:  "加载mysql中限制购买的数据到全局变量purchaseLimit",
		Data: nil,
	})
	return
	//response := jsonStruct.CommonResponse{
	//	Code: 8001,
	//	Msg:  "success",
	//	Data: nil,
	//}
	//err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response)
	//if err != nil {
	//	ctx.Error("internel error", fasthttp.StatusInternalServerError)
	//}
	//ctx.Response.Header.Set("Content-Type", "application/json")
}
