package shop_orm

import (
	"fmt"

	"gorm.io/gorm"
)

type Good struct {
	SelfDefine
	ProductCategory string `gorm:"" json:"productCategory"`
	ProductName     string `gorm:"" json:"productName"`
	Inventory       int    `gorm:"" json:"inventory"`
	Price           int    `gorm:"" json:"price"`
}

// 添加记录到goods table
func (g *Good) CreateGood(tx *gorm.DB) error {
	// check good info
	if err := g.CheckGoodParams(); err != nil {
		return err
	}
	// check if good exist
	if g.IfGoodExistByProductCategoryAndProductName() {
		return fmt.Errorf("创建商品时, 已存在")
	}
	if err := tx.Create(g).Error; err != nil {
		return fmt.Errorf("创建商品时, 出现错误 %s", err.Error())
	}
	return nil
}

// query good by productCategory && productName
func (g *Good) QueryGoodsByProductCategoryAndProductName(tx *gorm.DB) []*Good {
	var results []*Good
	conn.Model(&Good{}).Where("product_category=? AND product_name=?", g.ProductCategory, g.ProductName).Find(&results)
	return results
}

// 查找记录 goods table
// 把所有deleted_at=null的实例显示出来
func (g *Good) QueryGoods() ([]*Good, error) {
	var goods = make([]*Good, 0, 128)
	conn.Model(&Good{}).Find(&goods)
	return goods, nil
}

// 已知good的product_id和product_name, 更新good的inventory和product_name
func (g *Good) UpdateGood(tx *gorm.DB) error {
	if err := g.CheckGoodParams(); err != nil {
		return fmt.Errorf("更新商品信息时, 传入的参数有误")
	}
	// if !g.IfGoodExistByID() {
	// 	return fmt.Errorf("更新商品信息不存在")
	// }
	if err := tx.Model(&Good{}).Where("id=?", g.ID).Updates(Good{ProductCategory: g.ProductCategory, ProductName: g.ProductName, Inventory: g.Inventory, Price: g.Price}).Error; err != nil {
		return fmt.Errorf("更新商品时, 出现错误 %s", err.Error())
	}
	return nil
}

// increase or decrease product number by productID and number
func (g *Good) UpdateGoodByProductIDandPurchaseNum(tx *gorm.DB, productID int, purchaseNum int) error {
	// if err := tx.Raw("UPDATE seckill.goods SET inventory = inventory - ? WHERE id = ?", purchaseNum, productID).Error; err != nil {
	// 	return fmt.Errorf("更新商品信息时, 发生错误: %v", err)
	// }
	logger.Infoln("UPDATE goods SET inventory=inventory-? WHERE id=?", purchaseNum, productID)
	if err := tx.Exec("UPDATE goods SET inventory=? WHERE id=?", gorm.Expr("inventory+?", purchaseNum), productID).Error; err != nil {
		return fmt.Errorf("更新商品信息时, 发生错误: %v", err)
	}
	// if err := tx.Model(&Good{}).Where("id = ?", productID).Update("inventory", fmt.Sprintf("inventory-%d", purchaseNum)).Error; err!=nil {
	// 	return fmt.Errorf("更新商品信息时, 发生错误: %v", err)
	// }
	return nil
}

// 根据product_id确定唯一的一个商品, 然后删除它
func (g *Good) DeleteGood(tx *gorm.DB) error {
	if !g.IfGoodExistByID() {
		return fmt.Errorf("删除的商品不存在")
	}
	// if err := tx.Model(&Good{}).Where("id=?", g.ID).Update("deleted_at", time.Now()).Error; err != nil {
	// 	return err
	// }
	if err := tx.Model(&Good{}).Delete(g).Error; err != nil {
		return err
	}
	return nil
}

// 查找某个值为productId的商品是否存在
// productID 和 productName 都必须唯一
func (g *Good) IfGoodExistByProductCategoryAndProductName() bool {
	var result Good
	conn.Model(&Good{}).Where("product_category=? AND product_name=?", g.ProductCategory, g.ProductName).First(&result)
	if result.ProductName != "" || result.ProductCategory != "" {
		// if result.ProductName == g.ProductName || result.ProductCategory == g.ProductCategory
		logger.Warning("商品信息已存在")
		return true
	}
	return false
}

func (g *Good) IfGoodExistByID() bool {
	var result Good
	conn.Model(&Good{}).Where("id=?", g.ID).First(&result)
	if result.ID == g.ID {
		logger.Warning("商品信息已存在")
		return true
	}
	return false
}

func (g *Good) CheckGoodParams() error {
	if g.ProductCategory != "" && g.ProductName != "" && g.Inventory > 0 && g.Price > 0 {
		return nil
	}
	return fmt.Errorf("创建商品时, 检查商品参数时出错")
}
