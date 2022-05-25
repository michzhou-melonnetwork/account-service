package address

import (
	"net/http"

	"github.com/Melon-Network-Inc/account-service/internal/errors"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
	"github.com/Melon-Network-Inc/account-service/pkg/pagination"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/addresses/<id>", res.get)
	r.Get("/addresses", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/addresses", res.create)
	r.Put("/addresses/<id>", res.update)
	r.Delete("/addresses/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	address, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(address)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	addresses, err := r.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = addresses
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateAddressRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	address, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(address, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateAddressRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	address, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(address)
}

func (r resource) delete(c *routing.Context) error {
	address, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(address)
}
