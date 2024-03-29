package controllers2

import (
	"encoding/json"
	"go-seckill/internal/db"
	"go-seckill/internal/db/shop_orm"
	"go-seckill/internal/redisconf"
	"go-seckill/internal/utils"
	"strconv"

	"github.com/valyala/fasthttp"
)

// 一次创建一个PurchaseLimit实例, 并将其添加到purchase_limits table中
func CreatePurchaseLimit(ctx *fasthttp.RequestCtx) {
	p := new(shop_orm.PurchaseLimit)
	err := json.Unmarshal(ctx.Request.Body(), p)
	if err != nil {
		logger.Warnf("Unmarshal PurchaseLimit error happen: %v", err)
		utils.ResponseWithJson(ctx, 400, utils.CommonResponse{
			Code: 8400,
			Msg:  "解析PurchaseLimit时出现错误",
			Data: nil,
		})
		return
	}
	logger.Infof("解析后的PurchaseLimit是: %+v", p)

	// 首先查看, PurchaseLimit的product_id是否存在于purchase_limits表格中
	// if p.IfPurchaseLimitExist() {
	// 	logger.Warnf("PurchaseLimit已有相同ID的在表格中")
	// 	utils.ResponseWithJson(ctx, 404, utils.CommonResponse{
	// 		Code: 8404,
	// 		Msg:  "欲添加的PurchaseLimit已经存在于数据库中",
	// 		Data: nil,
	// 	})
	// 	return
	// }
	tx := db.Conn2.Begin()
	err = p.CreatePurchaseLimit(tx)
	if err != nil {
		logger.Warnf("当添加PurchaseLimit时, 错误: %v", err)
		tx.Rollback()
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "当添加PurchaseLimit时, 出错: " + err.Error(),
			Data: nil,
		})
		return
	}
	err = tx.Commit().Error
	if err != nil {
		// tx.Rollback()
		logger.Infof("CreatePurchaseLimit tx commit error: %v", err)
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "添加PurchaseLimit, 执行事务时失败",
			Data: nil,
		})
		return
	}
	// 更新runtime中的PurchaseLimitMap
	err = LoadGoodPurchaseLimit()
	if err != nil {
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "更新PurchaseLimitMap变量时失败",
			Data: nil,
		})
		return
	}
	logger.Infof("添加PurchaseLimit成功")
	utils.ResponseWithJson(ctx, 200, utils.CommonResponse{
		Code: 8200,
		Msg:  "添加PurchaseLimit成功",
		Data: nil,
	})
}

// 获取所有商品的purchase_limit
func QueryPurchaseLimits(ctx *fasthttp.RequestCtx) {
	// 首先, 获取id
	id := string(ctx.QueryArgs().Peek("id"))
	// log.Println(ctx.QueryArgs().Peek("id"))
	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			logger.Warnf("传入的数据类型非整数型: %T\n", idInt)
			utils.ResponseWithJson(ctx, 200, utils.CommonResponse{
				Code: 8200,
				Msg:  "传入的数据类型非整数型",
				Data: nil,
			})
			return
		}
		// 如果存在的话, 返回Purchase_limit这个对象
		p := new(shop_orm.PurchaseLimit)
		p.ProductID = idInt
		purchaseLimit := p.QueryPurchaseLimitByProductID()
		logger.Infof("purchaseLimit query succuesful")
		utils.ResponseWithJson(ctx, 200, utils.CommonResponse{
			Code: 8200,
			Msg:  "query PurchaseLimit successful",
			Data: purchaseLimit,
		})
		return
	}
	// 如果存在的话, 返回Purchase_limit这个对象
	p := new(shop_orm.PurchaseLimit)
	purchaseLimits := p.QueryPurchaseLimits()
	logger.Infof("purchaseLimit query succuesful")
	utils.ResponseWithJson(ctx, 200, utils.CommonResponse{
		Code: 8200,
		Msg:  "query PurchaseLimit successful",
		Data: purchaseLimits,
	})
}

func UpdatePurchaseLimit(ctx *fasthttp.RequestCtx) {
	p := new(shop_orm.PurchaseLimit)
	err := json.Unmarshal(ctx.Request.Body(), p)
	if err != nil {
		logger.Warnf("While unmarshal request.body(), error: %v", err)
		utils.ResponseWithJson(ctx, 400, utils.CommonResponse{
			Code: 8400,
			Msg:  "While unmarshal request.body(), error",
			Data: nil,
		})
		return
	}
	logger.Infof("unmarshal request.body() success")

	// 如果存在的话, 返回Purchase_limit这个对象
	tx := db.Conn2.Begin()
	err = p.UpdatePurchaseLimit(tx)
	if err != nil {
		logger.Warnf("UpdatePurchaseLimit transaction error: %v", err)
		tx.Rollback()
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "UpdatePurchaseLimit error " + err.Error(),
			Data: nil,
		})
		return
	}
	err = tx.Commit().Error
	if err != nil {
		// tx.Rollback()
		logger.Warnf("UpdatePurchaseLimit transaction commit error: %v", err)
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "UpdatePurchaseLimit 事务提交失败",
			Data: nil,
		})
		return
	}
	// 更新runtime中的PurchaseLimitMap
	err = LoadGoodPurchaseLimit()
	if err != nil {
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "更新PurchaseLimitMap变量时失败",
			Data: nil,
		})
		return
	}
	logger.Infof("UpdatePurchaseLimit transaction commit successful")
	utils.ResponseWithJson(ctx, 200, utils.CommonResponse{
		Code: 8200,
		Msg:  "UpdatePurchaseLimit 事务提交成功",
		Data: nil,
	})
}

func DeletePurchaseLimit(ctx *fasthttp.RequestCtx) {
	// delete purchaseLimit
	p := new(shop_orm.PurchaseLimit)
	err := json.Unmarshal(ctx.Request.Body(), p)
	if err != nil {
		logger.Warnf("While unmarshal request.body(), error: %v", err)
		utils.ResponseWithJson(ctx, 400, utils.CommonResponse{
			Code: 8400,
			Msg:  "While unmarshal request.body(), error",
			Data: nil,
		})
		return
	}
	logger.Infof("DeletePurchaseLimit success")

	// 如果存在的话, 返回Purchase_limit这个对象
	tx := db.Conn2.Begin()
	err = p.DeletePurchaseLimit(tx)
	if err != nil {
		logger.Warnf("DeletePurchaseLimit transaction error: %v", err)
		tx.Rollback()
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "DeletePurchaseLimit error " + err.Error(),
			Data: nil,
		})
		return
	}
	err = tx.Commit().Error
	if err != nil {
		// tx.Rollback()
		logger.Warnf("DeletePurchaseLimit transaction commit error: %v", err)
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "DeletePurchaseLimit transaction commit error",
			Data: nil,
		})
		return
	}
	// 更新runtime中的purchaseLimitMap
	err = LoadGoodPurchaseLimit()
	if err != nil {
		utils.ResponseWithJson(ctx, 500, utils.CommonResponse{
			Code: 8500,
			Msg:  "更新PurchaseLimitMap变量时失败",
			Data: nil,
		})
		return
	}
	logger.Infof("DeletePurchaseLimit transaction commit successful")
	utils.ResponseWithJson(ctx, 200, utils.CommonResponse{
		Code: 8200,
		Msg:  "DeletePurchaseLimit transaction commit successful",
		Data: nil,
	})
}

// SyncGoodsLimit ...
// 更新商品限制计划
// 例如, 在更新MySQL的限制购买条件后, 若要将商品购买限制同步到app中, 只需要调用goodsLimit这个接口就可以
func LoadGoodPurchaseLimit() error {

	// 加载limit限制计划
	err := redisconf.LoadLimits()
	if err != nil {
		logger.Warnf("SyncGoodsLimit: 加载limit变量到全局变量purchaseLimitMap时出现错误 %v", err)
		return err
	}
	logger.Infof("SyncGoodsLimit: 加载limit变量到全局变量purchaseLimit成功")
	return nil
}
