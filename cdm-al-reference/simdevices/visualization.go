/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package simdevices

import (
	"embed"
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gorilla/websocket"
)

var (
	//go:embed static/*
	staticFiles embed.FS

	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan []simulatedDeviceInfo)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for demo
		},
	}
)

func startDeviceVisualization(serverAddress string) {
	router := gin.New()
	router.Use(logger.SetLogger(logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
		return log.Logger
	})))

	// Serve static index file
	router.GET("/", func(c *gin.Context) {
		content, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading file")
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(content))
	})

	// WebSocket endpoint
	router.GET("/ws", handleWebSocket)

	// Start the broadcaster
	go handleBroadcast()

	// Start the webserver
	go runServer(router, serverAddress)
}

func runServer(router *gin.Engine, serverAddress string) {
	log.Info().Msgf("Starting visualization server at %s", serverAddress)
	if err := router.Run(serverAddress); err != nil {
		log.Fatal().Err(err).Msg("Starting of virtualization server failed")
	}
}

func handleWebSocket(context *gin.Context) {
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Err(err).Msg("Upgrade error")
		return
	}
	defer ws.Close()

	clients[ws] = true

	// Send initial device list
	sendDeviceList(ws)

	// Keep connection alive
	for {
		messageType, _, err := ws.ReadMessage()
		if err != nil {
			log.Err(err).Msg("WebSocket error")
			delete(clients, ws)
			break
		}

		if messageType == websocket.CloseMessage {
			delete(clients, ws)
			break
		}
	}
}

func handleBroadcast() {
	for {
		deviceList := <-broadcast
		for client := range clients {
			err := client.WriteJSON(deviceList)
			if err != nil {
				log.Err(err).Msg("WebSocket error")
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func sendDeviceList(ws *websocket.Conn) {
	deviceList := getDeviceListCopy(false)
	if err := ws.WriteJSON(deviceList); err != nil {
		log.Err(err).Msg("Error sending update")
	}
}

func broadcastDeviceUpdates(deviceList []simulatedDeviceInfo) {
	broadcast <- deviceList
}
