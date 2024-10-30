package handlers

import (
	"delivery/constants"
	"delivery/entities"
	"delivery/logger"
	htp "delivery/pkg/http"
	jwta "delivery/pkg/jwt"
	"delivery/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SendCode(c *gin.Context) {
	var req entities.SendCodeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.handleResponse(c, htp.BadRequest, err.Error())
		return
	}

	err = req.Validate()
	if err != nil {
		h.handleResponse(c, htp.InvalidArgument, err.Error())
		return
	}

	smscode, err := utils.GenerateVerificationCode()
	if err != nil {
		h.log.Error("kod generatsiya qilishda xatolik", logger.Error(err))
	}

	err = h.redis.Set(c, req.PhoneNumber, smscode, 2*time.Minute).Err()
	if err != nil {
		h.log.Error("Redisda kodni saqlashda xatolik1", logger.Error(err))
	}

	h.handleResponse(c, htp.OK, "Telefon raqamingizga 6 xonali kod yuborildi")
}


func (h *Handler) Registration(c *gin.Context) {
	var req entities.RegistrReq
	err := c.ShouldBindJSON(&req)

	if err != nil {
		h.handleResponse(c, htp.BadRequest, err.Error())
		return
	}

	err = req.Validate()
	if err != nil {
		h.handleResponse(c, htp.BadRequest, err.Error())
		return
	}

	resp, err := h.adminController.Registration(c, req)
	if err != nil {
		h.handleResponse(c, htp.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, htp.OK, resp)
}

func (h *Handler) InsertUserLocation(c *gin.Context) {
	var location entities.UserLocation
	err := c.ShouldBindJSON(&location)
	if err != nil {
		h.handleResponse(c, htp.BadRequest, err.Error())
	}
	userId, err := jwta.ExtractUserIDFromToken(c, []byte(h.cfg.JWTSecretKey))
	if err != nil {
		h.handleResponse(c, htp.InternalServerError, "Invalid or expired token")
		return
	}
	location.UserID = userId
	err = h.adminController.InsertUserLocation(c, location)
	if err != nil {
		h.handleResponse(c, htp.InternalServerError, err.Error())
		return
	}
	h.handleResponse(c, htp.OK, constants.Success)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req entities.UserProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, htp.BadRequest, err.Error())
		return
	}
	userId, err := jwta.ExtractUserIDFromToken(c, []byte(h.cfg.JWTSecretKey))
	if err != nil {
		h.handleResponse(c, StatusFromError(err), err.Error())
		return
	}
	req.ID=userId
	req.UpdatedBy = userId
	req.UpdatedAt = time.Now()
	if err := h.adminController.UpdateUserProfile(c.Request.Context(), req); err != nil {
		h.handleResponse(c, htp.InternalServerError, err.Error())
		return
	}
	h.handleResponse(c, htp.OK, "Profile updated successfully!")
}

func (h *Handler) GetProfile(c *gin.Context) {
	userId, err := jwta.ExtractUserIDFromToken(c, []byte(h.cfg.JWTSecretKey))
	if err != nil {
		h.handleResponse(c, htp.InternalServerError, "Invalid or expired token")
		return
	}

	data, err := h.adminController.GetUserProfile(c, userId)
	if err != nil {
		h.handleResponse(c, htp.InternalServerError, constants.InternelServError)
		return
	}
	h.handleResponse(c, htp.OK, data)
}

func (h *Handler) GetUserLocation(c *gin.Context) {
	userId, err := jwta.ExtractUserIDFromToken(c, []byte(h.cfg.JWTSecretKey))
	if err != nil {
		h.handleResponse(c, htp.InternalServerError, "Invalid or expired token")
		return
	}
	data, err := h.adminController.GetUserLocation(c, userId)
	if err != nil {
		h.handleResponse(c, htp.InternalServerError, constants.InternelServError)
		return
	}
	h.handleResponse(c, htp.OK, data)
}



// func (h *Handler) Login(c *gin.Context) {
// 	loginReq := entities.LoginReq{}
// 	err := c.ShouldBindJSON(&loginReq)
// 	if err != nil {
// 		h.handleResponse(c, htp.BadRequest, err.Error())
// 		return
// 	}

// 	err = loginReq.Validate()
// 	if err != nil {
// 		h.handleResponse(c, htp.InvalidArgument, err.Error())
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.ContextTimeoutDuration)
// 	defer cancel()

// 	resp, err := h.adminController.Login(
// 		ctx,
// 		loginReq,
// 	)
// 	if err != nil {
// 		h.handleResponse(c, StatusFromError(err), err.Error())
// 		return
// 	}

// 	h.handleResponse(c, htp.OK, resp)
// }
