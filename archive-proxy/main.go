package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/figment-networks/indexing-engine/store/jsonquery"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/figment-networks/mina-indexer/archive-proxy/queries"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	log.Println("connecting to database...")
	conn, err := initConn(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.LogMode(os.Getenv("TRACE_SQL") == "1")

	router := gin.Default()

	router.GET("/", handleInfo(conn))
	router.GET("/status", handleStatus(conn))
	router.GET("/blocks", handleBlocks(conn))
	router.GET("/blocks/:hash", handleBlock(conn))
	router.GET("/blocks/:hash/user_commands", handleUserCommands(conn))
	router.GET("/blocks/:hash/internal_commands", handleInternalCommands(conn))
	router.GET("/public_keys", handlePublicKeys(conn))
	router.GET("/public_keys/:id", handlePublicKey(conn))
	router.GET("/staking_ledger", handleStakingLedger())

	listenAddr := os.Getenv("PORT")
	if listenAddr == "" {
		listenAddr = "3088"
	}
	listenAddr = "0.0.0.0:" + listenAddr

	log.Println("starting server on", listenAddr)
	if err := router.Run(listenAddr); err != nil {
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
		result, err = jsonquery.Object(conn, jsonquery.Prepare(query), args)
		if result == nil && err == nil {
			ctx.AbortWithStatusJSON(404, gin.H{"error": "record not found"})
			return
		}
	} else {
		result, err = jsonquery.MustArray(conn, jsonquery.Prepare(query), args)
	}

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(200, "application/json", result)
}

type blocksParams struct {
	StartHeight uint `form:"start_height"`
	Limit       uint `form:"limit"`
}

func handleInfo(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderQuery(c, conn, "object", queries.Info)
	}
}

func handleBlocks(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := blocksParams{}
		if err := c.Bind(&params); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err})
			return
		}

		if params.Limit == 0 {
			params.Limit = 100
		}
		if params.Limit > 1000 {
			params.Limit = 1000
		}

		renderQuery(c, conn, "array", queries.Blocks, params.StartHeight, params.Limit)
	}
}

func handleBlock(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderQuery(c, conn, "object", queries.Block, c.Param("hash"))
	}
}

func handleUserCommands(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderQuery(c, conn, "array", queries.UserCommands, c.Param("hash"))
	}
}

func handleInternalCommands(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderQuery(c, conn, "array", queries.InternalCommands, c.Param("hash"))
	}
}

func handlePublicKeys(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderQuery(c, conn, "object", queries.PublicKeys)
	}
}

func handlePublicKey(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderQuery(c, conn, "object", queries.PublicKey, c.Param("id"))
	}
}

func handleStatus(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := conn.Exec("SELECT 1").Error
		if err != nil {
			log.Println("db connection test erro:", err)
			c.AbortWithStatusJSON(500, gin.H{"healthy": false})
			return
		}
		c.JSON(200, gin.H{"healthy": true})
	}
}

// TODO: move ledger dumping into background, run the command every 60s
func handleStakingLedger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ledgerType := "current"
		if c.Query("type") == "next" {
			ledgerType = "next"
		}

		ledgerbuf := bytes.NewBuffer(nil)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		cmd := exec.CommandContext(ctx, "coda", "advanced", "dump-staking-ledger", ledgerType, "-json")
		cmd.Stderr = os.Stderr
		cmd.Stdout = ledgerbuf

		if err := cmd.Run(); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Data(200, "application/json", ledgerbuf.Bytes())
	}
}
