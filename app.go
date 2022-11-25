package main

import (
	"context"
	"log"

	"github.com/margen2/goknition/backend/api"
	"github.com/margen2/goknition/backend/controllers"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	load()
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

// Collections
// =====================

func (a *App) GetCollections(refresh bool) ([]string, error) {
	if refresh {
		err := api.RefreshCollections()
		if err != nil {
			return nil, err
		}
	}

	return api.ListCollections(), nil
}

func (a *App) CreateCollection(collectionID string) error {
	err := controllers.CreateCollection(collectionID)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteCollection(collectionID string) error {
	err := controllers.DeleteCollection(collectionID)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) IndexFaces(collectionID, path string) error {
	err := controllers.IndexFaces(collectionID, path)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) GetFaces(collectionID string) ([]string, error) {
	faces, err := controllers.GetFaces(collectionID)
	if err != nil {
		return nil, err
	}
	return faces, nil
}

// Search Faces
// =====================
func (a *App) SearchFaces(collectionID string) error {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Data Folder",
	})
	if err != nil {
		return err
	}

	err = controllers.SearchImages(collectionID, dir)
	if err != nil {
		return err
	}
	return nil
}
