package user

import (
	"github.com/annuums/solanum/container"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func retrieveUser(c *gin.Context) {
	repo := container.DepFromGinContext[UserRepository](c, "userRepository")

	users, err := repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func addUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := &User{Name: req.Name, Email: req.Email, CreatedAt: time.Now()}
	repo := container.DepFromContext[UserRepository](c, "userRepository")
	if err := repo.Create(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}
