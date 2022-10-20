package main

import (
	"context"
	"fmt"
	"log"

	"github.com/margen2/goknition/backend/api"
	"github.com/margen2/goknition/backend/data"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	err := api.RefreshCollections()
	if err != nil {
		log.Fatal(err)
	}

	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) GetCollections(refresh bool) []string {
	if refresh {
		err := api.RefreshCollections()
		if err != nil {
			log.Fatal(err)
		}
	}

	return api.ListCollections()
}

func (a *App) GetCwd() string {
	cwd, err := data.GetCwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

func (a *App) ListFolders(folder string) []string {
	folders, err := data.ListFolders(folder)
	if err != nil {
		log.Fatal(err)
	}

	return folders
}

func (a *App) GoBack(folder string) string {
	fmt.Println(folder)
	folder = data.GoBack(folder)
	fmt.Println(folder)
	return folder
}
