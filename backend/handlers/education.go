package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Syzua/degree-auth/backend/middleware"
	"github.com/Syzua/degree-auth/backend/models"
	"github.com/Syzua/degree-auth/backend/services"
)

func AddEducation(c *gin.Context) {
	var edu models.Education
	if err := c.ShouldBindJSON(&edu); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Code: 400, Message: "invalid request"})
		return
	}

	args := [][]byte{
		[]byte(edu.CertNo),
		[]byte(edu.Name),
		[]byte(edu.StudentID),
		[]byte(edu.School),
		[]byte(edu.Major),
		[]byte(edu.Degree),
		[]byte(edu.GraduationDate),
	}

	_, err := services.InvokeChaincode("AddEducation", args, "Org1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Code: 0, Message: "education record added", Data: edu})
}

func UpdateEducation(c *gin.Context) {
	var edu models.Education
	if err := c.ShouldBindJSON(&edu); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Code: 400, Message: "invalid request"})
		return
	}

	args := [][]byte{
		[]byte(edu.CertNo),
		[]byte(edu.Name),
		[]byte(edu.StudentID),
		[]byte(edu.School),
		[]byte(edu.Major),
		[]byte(edu.Degree),
		[]byte(edu.GraduationDate),
	}

	_, err := services.InvokeChaincode("UpdateEducation", args, "Org1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Code: 0, Message: "education record updated", Data: edu})
}

func QueryEducationByID(c *gin.Context) {
	certNo := c.Param("certNo")
	if certNo == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Code: 400, Message: "certNo is required"})
		return
	}

	result, err := services.QueryChaincode("QueryEducationByID", [][]byte{[]byte(certNo)}, "Org1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	var edu models.Education
	json.Unmarshal(result, &edu)

	c.JSON(http.StatusOK, models.APIResponse{Code: 0, Message: "success", Data: edu})
}

func VerifyEducation(c *gin.Context) {
	certNo := c.Query("certNo")
	name := c.Query("name")
	if certNo == "" || name == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Code: 400, Message: "certNo and name are required"})
		return
	}

	result, err := services.QueryChaincode("VerifyByCertNoAndName", [][]byte{[]byte(certNo), []byte(name)}, "Org1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	var verified bool
	json.Unmarshal(result, &verified)

	c.JSON(http.StatusOK, models.APIResponse{Code: 0, Message: "success", Data: gin.H{"verified": verified}})
}

func GetHistory(c *gin.Context) {
	certNo := c.Param("certNo")
	if certNo == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Code: 400, Message: "certNo is required"})
		return
	}

	result, err := services.QueryChaincode("GetHistoryByID", [][]byte{[]byte(certNo)}, "Org1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	var records []interface{}
	json.Unmarshal(result, &records)

	c.JSON(http.StatusOK, models.APIResponse{Code: 0, Message: "success", Data: records})
}

func AuthorizeViewer(c *gin.Context) {
	certNo := c.Param("certNo")
	viewerID := c.PostForm("viewerID")
	if certNo == "" || viewerID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Code: 400, Message: "certNo and viewerID are required"})
		return
	}

	_, err := services.InvokeChaincode("AuthorizeViewer", [][]byte{[]byte(certNo), []byte(viewerID)}, "Org1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Code: 0, Message: "viewer authorized"})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	users := map[string]struct {
		Pass string
		Role string
	}{
		"admin":    {"admin123", "university"},
		"employer": {"emp123", "employer"},
		"student":  {"stu123", "student"},
	}

	user, ok := users[username]
	if !ok || user.Pass != password {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Code: 401, Message: "invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Code: 500, Message: "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Code: 0, Message: "login success", Data: gin.H{"token": token, "role": user.Role}})
}
