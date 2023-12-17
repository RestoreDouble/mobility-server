package controller

import (
	"math/rand"
	"mobility-server/ent/customer"
	"mobility-server/src/database"
	"mobility-server/src/service"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var usersOTP = make(map[int]string)

func generateOTP() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return strconv.Itoa(rng.Intn(999999) + 1)
}

func GetOtp(c *gin.Context) {
	phoneNumber := c.Query("phone")

	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	phone, err := strconv.Atoi(phoneNumber)
	if err != nil {
		// ... handle error
		panic(err)
	}

	otp := generateOTP()
	usersOTP[phone] = otp

	// Send OTP via SMS
	c.JSON(http.StatusOK, gin.H{"message": "sent OTP"})
}

func VerifyOTP(c *gin.Context) {
	phoneNumber := c.Query("phone")
	clientOTP := c.Query("otp")

	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	phone, err := strconv.Atoi(phoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		panic(err)
	}

	if clientOTP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP is required"})
		return
	}

	userOTP, exists := usersOTP[phone]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "No OTP found for the given phone number"})
		return
	}

	if clientOTP != userOTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// OTP validated successfully
	delete(usersOTP, phone)

	var res []struct {
		IsNew bool `json:"is_new"`
	}

	queryErr := database.Client.Customer.Query().
		Where(customer.Phone(phone)).
		Select(customer.FieldIsNew).
		Scan(c, &res)

	if queryErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
		return
	}

	claims := &service.UserClaims{
		Phone: phone,
		IsNew: true,
	}

	if len(res) > 0 {
		// Create a new token with claims
		claims.IsNew = res[0].IsNew
	} else {
		_, createErr := database.Client.Customer.Create().
			SetPhone(phone).
			SetUpdatedAt(time.Now()).
			Save(c)

		if createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": createErr.Error()})
			return
		}
	}

	tokenString, tokenErr := service.GenerateAccessToken(claims)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tokenErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_new": claims.IsNew,
		"token":  tokenString,
	})
}
