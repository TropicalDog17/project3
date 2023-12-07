package main

import (
	"github.com/TropicalDog17/project3/client/editor"
	"github.com/TropicalDog17/project3/crdt"
	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
)

type UIConfig struct {
	EditorConfig editor.EditorConfig
}

// initUI creates a new editor view and runs the main loop.
func initUI(conn *websocket.Conn, conf UIConfig) error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	e = editor.NewEditor(conf.EditorConfig)
	e.SetSize(termbox.Size())
	e.SetText(crdt.Content(doc))
	e.SendDraw()
	e.IsConnected = true

	go handleStatusMsg()

	go drawLoop()

	err = mainLoop(conn)
	if err != nil {
		return err
	}

	return nil
}

// mainLoop is the main update loop for the UI.
func mainLoop(conn *websocket.Conn) error {
	// Receiving termbox event
	termboxChan := getTermboxChan()

	// Receiving msg event
	msgChan := getMsgChan(conn)

	// Handle events
	for {
		select {
		case termboxEvent := <-termboxChan:
			err := handleTermboxEvent(termboxEvent, conn)
			if err != nil {
				return err
			}
		case msg := <-msgChan:
			handleMsg(msg, conn)
		}
	}
}
