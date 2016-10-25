package main

import (
	"io"

	"github.com/goadesign/goa"
	"github.com/goadesign/gorma-cellar/app"
	"github.com/goadesign/gorma-cellar/models"
	"golang.org/x/net/websocket"
)

// BottleController implements the bottle resource.
type BottleController struct {
	*goa.Controller
}

// NewBottleController creates a bottle controller.
func NewBottleController(service *goa.Service) *BottleController {
	return &BottleController{Controller: service.NewController("bottle")}
}

// Create runs the create action.
func (c *BottleController) Create(ctx *app.CreateBottleContext) error {
	b := models.Bottle{}
	b.AccountID = ctx.AccountID
	b.Color = ctx.Payload.Color
	b.Country = ctx.Payload.Country
	b.Name = ctx.Payload.Name
	b.Region = ctx.Payload.Region
	b.Review = ctx.Payload.Review
	b.Sweetness = ctx.Payload.Sweetness
	b.Varietal = ctx.Payload.Varietal
	b.Vineyard = ctx.Payload.Vineyard
	b.Vintage = ctx.Payload.Vintage
	err := bdb.Add(ctx.Context, &b)
	if err != nil {
		return ErrDatabaseError(err)
	}
	ctx.ResponseData.Header().Set("Location", app.BottleHref(b.AccountID, b.ID))
	return ctx.Created()
}

// Delete runs the delete action.
func (c *BottleController) Delete(ctx *app.DeleteBottleContext) error {
	err := bottleDBFilterByAccount(ctx.AccountID).Delete(ctx.Context, ctx.BottleID)
	if err != nil {
		return ErrDatabaseError(err)
	}
	return ctx.NoContent()
}

func bottleDBFilterByAccount(accountID int) *models.BottleDB {
	return models.NewBottleDB(bdb.Db.Scopes(models.BottleFilterByAccount(accountID, bdb.Db)))
}

// List runs the list action.
func (c *BottleController) List(ctx *app.ListBottleContext) error {
	bottles := bdb.ListBottle(ctx.Context, ctx.AccountID)
	return ctx.OK(bottles)
}

// Rate runs the rate action.
func (c *BottleController) Rate(ctx *app.RateBottleContext) error {
	// TBD: implement
	return nil
}

// Show runs the show action.
func (c *BottleController) Show(ctx *app.ShowBottleContext) error {
	// TBD: implement
	res := &app.Bottle{}
	return ctx.OK(res)
}

// Update runs the update action.
func (c *BottleController) Update(ctx *app.UpdateBottleContext) error {
	// TBD: implement
	return nil
}

// Watch runs the watch action.
func (c *BottleController) Watch(ctx *app.WatchBottleContext) error {
	c.WatchWSHandler(ctx).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

// WatchWSHandler establishes a websocket connection to run the watch action.
func (c *BottleController) WatchWSHandler(ctx *app.WatchBottleContext) websocket.Handler {
	return func(ws *websocket.Conn) {
		// TBD: implement
		ws.Write([]byte("watch bottle"))
		// Dummy echo websocket server
		io.Copy(ws, ws)
	}
}
