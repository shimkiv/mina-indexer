package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/figment-networks/indexing-engine/store/jsonquery"
	"github.com/figment-networks/mina-indexer/archive-proxy/queries"
)

func main() {
	conn, err := initConn(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	router := gin.Default()

	router.GET("/blocks", func(c *gin.Context) {
		var after int
		var err error
		var limit int

		if val := c.Query("after"); val != "" {
			after, err = strconv.Atoi(val)
			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": err})
				return
			}
		}

		if val := c.Query("limit"); val != "" {
			limit, err = strconv.Atoi(val)
			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": err})
				return
			}
		}
		if limit <= 0 {
			limit = 100
		}

		renderQuery(c, conn, "array", queries.Blocks, after, limit)
	})

	router.GET("/blocks/:hash", func(c *gin.Context) {
		renderQuery(c, conn, "object", queries.Block, c.Param("hash"))
	})

	router.GET("/blocks/:hash/user_commands", func(c *gin.Context) {
		renderQuery(c, conn, "array", queries.UserCommands, c.Param("hash"))
	})

	router.GET("/blocks/:hash/internal_commands", func(c *gin.Context) {
		renderQuery(c, conn, "array", queries.InternalCommands, c.Param("hash"))
	})

	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "3088"
	}

	if err := router.Run(":" + listenPort); err != nil {
		log.Fatal(err)
	}
}

func initConn(str string) (*gorm.DB, error) {
	conn, err := gorm.Open("postgres", str)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func renderQuery(ctx *gin.Context, conn *gorm.DB, mode string, query string, args ...interface{}) {
	var result []byte
	var err error

	if mode == "object" {
		result, err = jsonquery.MustObject(conn, jsonquery.Prepare(query), args)
	} else {
		result, err = jsonquery.MustArray(conn, jsonquery.Prepare(query), args)
	}

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(200, "application/json", result)
}
