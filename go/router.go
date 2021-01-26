package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addRoutes(r *gin.Engine) {
	db, err := taskDB()
	if err != nil {
		log.Fatal(err)
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"pong": true,
		})
	})
	r.GET("/tasks", func(c *gin.Context) {
		tasks, err := getAllTasks(db)
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, tasks)
	})
	r.POST("/tasks", func(c *gin.Context) {
		var json NewTask
		err := c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var task *Task
		task = &Task{
			Name: json.Name,
		}
		err = createTask(db, task)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusAccepted, task)
	})
	r.PATCH("/tasks/:id", func(c *gin.Context) {
		oid, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var json UpdatedTask
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		task, err := updateTask(db, oid, json.Done)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusAccepted, task)
	})
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		oid, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		err = deleteTask(db, oid)
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatus(http.StatusOK)
	})
}
