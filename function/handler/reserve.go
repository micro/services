package handler

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/function/proto"
	"github.com/micro/services/pkg/tenant"
)

var (
	mtx sync.Mutex
)

type Reservation struct {
	// The function name
	Name string `json:"name"`
	// The owner e.g tenant id
	Owner string `json:"owner"`
	// Uniq associated token
	Token string `json:"token"`
	// Time of creation
	Created time.Time `json:"created"`
	// The expiry time
	Expires time.Time `json:"expires"`
}

func genToken(name, owner string) string {
	h := sha1.New()
	io.WriteString(h, name+owner)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Call is a single request handler called via client.Call or the generated client code
func (f *Function) Reserve(ctx context.Context, req *pb.ReserveRequest, rsp *pb.ReserveResponse) error {
	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "micro"
	}

	if len(req.Name) == 0 {
		return errors.BadRequest("function.reserve", "missing function name")
	}

	if len(req.Name) < 3 || len(req.Name) > 32 {
		return errors.BadRequest("function.reserve", "name must be longer than 3-32 chars in length")
	}

	if !NameFormat.MatchString(req.Name) {
		return errors.BadRequest("function.reserve", "invalidate name format")
	}

	// to prevent race conditions in reservation lets global lock
	mtx.Lock()
	defer mtx.Unlock()

	// check the store for reservation
	recs, err := store.Read(ReservationKey + req.Name)
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("function.reserve", "failed to reserve name")
	}

	var rsrv *Reservation

	// check if the record exists
	if len(recs) > 0 {
		// existing reservation exists
		rec := recs[0]

		if err := rec.Decode(&rsrv); err != nil {
			return errors.BadRequest("function.reserve", "name already reserved")
		}

		// check the owner matches or if the reservation expired
		if rsrv.Owner != id && rsrv.Expires.After(time.Now()) {
			return errors.BadRequest("function.reserve", "name already reserved")
		}

		// update the owner
		rsrv.Owner = id

		// update the reservation expiry
		rsrv.Expires = time.Now().AddDate(1, 0, 0)
	} else {
		// check if its already running
		key := FunctionKey + req.Name
		recs, err := store.Read(key, store.ReadLimit(1))
		if err != nil && err != store.ErrNotFound {
			return errors.InternalServerError("function.reserve", "failed to reserve name")
		}

		// existing function is running by that name
		if len(recs) > 0 {
			return errors.BadRequest("function.reserve", "function already exists")
		}

		// not reserved
		rsrv = &Reservation{
			Name:    req.Name,
			Owner:   id,
			Created: time.Now(),
			Expires: time.Now().AddDate(1, 0, 0),
			Token:   genToken(req.Name, id),
		}
	}

	rec := store.NewRecord(ReservationKey+req.Name, rsrv)

	if err := store.Write(rec); err != nil {
		return errors.InternalServerError("function.reserve", "error while reserving name")
	}

	rsp.Reservation = &pb.Reservation{
		Name:    rsrv.Name,
		Owner:   rsrv.Owner,
		Created: rsrv.Created.Format(time.RFC3339Nano),
		Expires: rsrv.Expires.Format(time.RFC3339Nano),
		Token:   rsrv.Token,
	}

	return nil
}
