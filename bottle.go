package main

import (
	"io"

	"github.com/goadesign/goa"
	"github.com/goadesign/gorma-cellar/app"
	"github.com/goadesign/gorma-cellar/models"
	"github.com/jinzhu/gorm"
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
	b, err := bdb.Get(ctx.Context, ctx.BottleID)
	if err == gorm.ErrRecordNotFound || (err == nil && b.AccountID != ctx.AccountID) {
		return ctx.NotFound()
	} else if err != nil {
		return ErrDatabaseError(err)
	}
	b.Rating = ctx.Payload.Rating
	err = bdb.Update(ctx, b)
	if err != nil {
		return ErrDatabaseError(err)
	}
	return ctx.NoContent()
}

// Show runs the show action.
func (c *BottleController) Show(ctx *app.ShowBottleContext) error {
	bottle, err := bdb.OneBottleFull(ctx.Context, ctx.BottleID, ctx.AccountID)
	if err == gorm.ErrRecordNotFound {
		return ctx.NotFound()
	} else if err != nil {
		return ErrDatabaseError(err)
	}
	bottle.Href = app.BottleHref(ctx.AccountID, ctx.BottleID)
	return ctx.OKFull(bottle)
}

// Update runs the update action.
func (c *BottleController) Update(ctx *app.UpdateBottleContext) error {
	b, err := bdb.Get(ctx.Context, ctx.BottleID)
	if err == gorm.ErrRecordNotFound || (err == nil && b.AccountID != ctx.AccountID) {
		return ctx.NotFound()
	} else if err != nil {
		return ErrDatabaseError(err)
	}
	if ctx.Payload.Color != nil {
		b.Color = *ctx.Payload.Color
	}
	b.Country = ctx.Payload.Country
	if ctx.Payload.Name != nil {
		b.Name = *ctx.Payload.Name
	}
	b.Region = ctx.Payload.Region
	b.Review = ctx.Payload.Review
	b.Sweetness = ctx.Payload.Sweetness
	if ctx.Payload.Varietal != nil {
		b.Varietal = *ctx.Payload.Varietal
	}
	if ctx.Payload.Vineyard != nil {
		b.Vineyard = *ctx.Payload.Vineyard
	}
	if ctx.Payload.Vintage != nil {
		b.Vintage = *ctx.Payload.Vintage
	}
	err = bdb.Update(ctx, b)
	if err != nil {
		return ErrDatabaseError(err)
	}
	return ctx.NoContent()
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
