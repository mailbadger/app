package main

/*
news-maily - Open-Source Newsletter Mail System

The MIT License (MIT)

Copyright (c) 2015 Filip Nikolovski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is furnished
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
import (
	"fmt"
	"runtime"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	error := entities.Setup()
	if error != nil {
		fmt.Println(error)
	}

	r := gin.Default()
	r.GET("/user", func(c *gin.Context) {
		user, err := entities.GetUser(1)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err,
			})
		}
		c.JSON(200, gin.H{
			"user": user,
		})
	})

	r.GET("/templates", middleware.Paginate, func(c *gin.Context) {
		pagination, ok := c.MustGet("pagination").(middleware.Pagination)
		if !ok {
			c.JSON(400, gin.H{
				"status":  "failure",
				"message": "There was an error.",
			})
		}

		entities.GetTemplates(1, &pagination)

		c.JSON(200, gin.H{
			"collection": pagination.Collection,
			"total":      pagination.Total,
			"page":       pagination.Page,
			"per_page":   pagination.PerPage,
		})
	})

	r.Run() // listen and server on 0.0.0.0:8080
}
