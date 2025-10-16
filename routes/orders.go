package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"hotsauceshop/lib"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const CouponCodeMaxLength = 25
const CouponCodeMinLength = 4

func Orders(r *gin.Engine, dbPool *pgxpool.Pool, logger *slog.Logger) {
	r.GET("/api/v1/orders/shipping-options", func(c *gin.Context) {
		shippingOptions, err := lib.GetShippingOptions(dbPool)
		if err != nil {
			logger.Error(fmt.Sprintf("Error getting shipping options: %s", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Failed to get shipping options",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"shippingOptions": shippingOptions,
			},
		})
	})

	r.GET("/api/v1/coupons/:code", func(c *gin.Context) {
		couponCode := c.Param("code")
		if len(couponCode) < CouponCodeMinLength || len(couponCode) > CouponCodeMaxLength {
			logger.Error("Coupon code is too short or too long")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "ERROR",
				"message": "Invalid coupon code",
			})
			return
		}

		validCouponCode, couponCodeErr := lib.GetCouponByCode(dbPool, couponCode)
		if couponCodeErr != nil || validCouponCode == (lib.CouponCode{}) {
			logger.Error(fmt.Sprintf("GetCouponByCode error: %v", couponCodeErr))
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "ERROR",
				"message": "Invalid coupon code",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"results": gin.H{
				"couponCode": validCouponCode,
			},
		})
	})
}
