package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"inventory-management/backend/internal/models"
)

type ProductHandler struct {
	db *gorm.DB
}

type ProductRequest struct {
	Name     string `json:"name" binding:"required"`
	SKU      string `json:"sku" binding:"required"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
	Remark   string `json:"remark"`
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	var products []models.Product
	query := h.db.Order("id DESC")

	keyword := c.Query("keyword")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name ILIKE ? OR sku ILIKE ?", like, like)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取库存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": products})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "商品名称和 SKU 不能为空"})
		return
	}

	operatorID := currentUserID(c)
	var product models.Product
	err := h.db.Transaction(func(tx *gorm.DB) error {
		product = models.Product{
			Name:     req.Name,
			SKU:      req.SKU,
			Category: req.Category,
			Quantity: req.Quantity,
			Remark:   req.Remark,
		}
		if err := tx.Create(&product).Error; err != nil {
			return err
		}
		if req.Quantity != 0 {
			return createStockMovement(tx, product.ID, 0, req.Quantity, operatorID, "新增商品初始库存")
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "创建商品失败，请检查 SKU 是否重复"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "商品名称和 SKU 不能为空"})
		return
	}

	operatorID := currentUserID(c)
	var product models.Product
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&product, productID).Error; err != nil {
			return err
		}

		beforeQuantity := product.Quantity
		product.Name = req.Name
		product.SKU = req.SKU
		product.Category = req.Category
		product.Quantity = req.Quantity
		product.Remark = req.Remark

		if err := tx.Save(&product).Error; err != nil {
			return err
		}
		if beforeQuantity != req.Quantity {
			return createStockMovement(tx, product.ID, beforeQuantity, req.Quantity, operatorID, req.Remark)
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "更新商品失败，请检查商品是否存在或 SKU 是否重复"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.db.Delete(&models.Product{}, productID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *ProductHandler) ListStockMovements(c *gin.Context) {
	productID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var movements []models.StockMovement
	if err := h.db.Preload("Operator").
		Where("product_id = ?", productID).
		Order("id DESC").
		Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取库存变动记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": movements})
}

func createStockMovement(tx *gorm.DB, productID uint, beforeQuantity int, afterQuantity int, operatorID uint, remark string) error {
	return tx.Create(&models.StockMovement{
		ProductID:      productID,
		BeforeQuantity: beforeQuantity,
		AfterQuantity:  afterQuantity,
		ChangeQuantity: afterQuantity - beforeQuantity,
		OperatorID:     operatorID,
		Remark:         remark,
	}).Error
}

func currentUserID(c *gin.Context) uint {
	value, exists := c.Get("current_user")
	if !exists {
		return 0
	}
	user, ok := value.(models.User)
	if !ok {
		return 0
	}
	return user.ID
}
