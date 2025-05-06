package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

/////// GET REQS

func (h *Handler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	uname, err := c.Cookie("uname")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return

	}
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	getUserReq := &GetUsersReq{
		UserName: uname,
		Token:    token,
	}

	users, err := h.Service.getUsers(ctx, getUserReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) GetUserByUsername(c *gin.Context) {
	ctx := c.Request.Context()
	var u GetUserByUsernameReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Service.getUserByUsername(ctx, u.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
func (h *Handler) GetUserByQuery(c *gin.Context) {
	ctx := c.Request.Context()
	query := c.Request.URL.Query()

	users, err := h.Service.getUsersByQuery(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

///// POST REQS

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("uname", "", -1, "/", "", true, true)
	c.SetCookie("token", "", -1, "/", "", true, true)

	c.JSON(http.StatusResetContent, nil)
}

func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var body CreateUserReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//if body.Password != body.PasswordConfirm {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
	//	return
	//}

	_, err := h.Service.register(ctx, &body)
	if err != nil {
		// TODO: Match Errors if internal server or uname exists
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true})
}

func (h *Handler) CreateUsers(c *gin.Context) {
	ctx := c.Request.Context()
	var body CreateUsersReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := h.Service.createUsers(ctx, &body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true})
}

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var body LoginUserReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.login(ctx, &body)
	if err != nil {
		// TODO: Match Errors if internal or wrong credential
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", res.Token, 0, "", "", true, true)
	c.SetCookie("uname", res.UserName, 0, "", "", true, true)

	res.Token = ""

	c.JSON(http.StatusOK, gin.H{"user": res})
}

/////////// PUT REQS

func (h *Handler) ChangePwd(c *gin.Context) {
	ctx := c.Request.Context()
	var body ChangePasswordReq
	uname, err := c.Cookie("uname")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body.UserName = uname

	err = h.Service.changePassword(ctx, &body, token)
	if err != nil {
		switch err.Error() {
		case "incorrect":
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		case "invalid":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true})

}
