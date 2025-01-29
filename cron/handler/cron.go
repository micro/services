package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/micro/micro/v5/service/errors"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"
	pb "github.com/micro/services/cron/proto"
	"github.com/micro/services/pkg/tenant"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	sync.Mutex
	jobs map[string]*cron.Cron
}

func New() *Cron {
	c := &Cron{
		jobs: make(map[string]*cron.Cron),
	}
	c.Start()
	return c
}

var (
	jobPrefix = "job"
)

func (c *Cron) Start() {
	go func() {
		limit := uint(100)
		offset := uint(0)
		// read all jobs

		for {
			recs, err := store.Read(jobPrefix+"/", store.ReadPrefix(), store.ReadLimit(limit), store.ReadOffset(offset))
			if err != nil {
				log.Errorf("Failed to start: %v", err)
				return
			}

			// when no records are left leave
			if len(recs) == 0 {
				break
			}

			c.Lock()
			for _, rec := range recs {
				job := new(pb.Job)
				rec.Decode(&job)
				cr := c.Setup(job)
				c.jobs[rec.Key] = cr
			}
			c.Unlock()

			// update the offset
			offset += 100
		}

	}()
}

func (c *Cron) Setup(job *pb.Job) *cron.Cron {
	log.Infof("Setting up job id: %s", job.Id)
	// schedule the job
	cr := cron.New()
	cr.AddFunc(job.Interval, func() {
		log.Infof("Running job id: %s", job.Id)
		rsp, err := http.Get(job.Callback)
		if err != nil {
			log.Errorf("Failed job id: %s error: %v", job.Id, err)
			return
		}
		defer rsp.Body.Close()
		b, _ := ioutil.ReadAll(rsp.Body)
		if rsp.StatusCode != 200 {
			log.Errorf("Non 200 job id: %s error: %s", job.Id, string(b))
			return
		}

		// TODO: save job state
		log.Infof("Successful job id: %s", job.Id)
	})
	cr.Start()
	return cr
}

func (c *Cron) Schedule(ctx context.Context, req *pb.ScheduleRequest, rsp *pb.ScheduleResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("cron.schedule", "missing name")
	}

	if len(req.Id) == 0 {
		req.Id = uuid.New().String()
	}

	if len(req.Interval) == 0 {
		return errors.BadRequest("cron.schedule", "missing interval")
	}

	if len(req.Callback) == 0 {
		return errors.BadRequest("cron.schedule", "missing callback")
	}

	tnt, _ := tenant.FromContext(ctx)

	key := fmt.Sprintf("%s/%s/%s", jobPrefix, tnt, req.Id)

	c.Lock()
	defer c.Unlock()

	// check local
	if _, ok := c.jobs[key]; ok {
		return errors.BadRequest("cron.schedule", "job already exists")
	}

	// check if it exists in store and unscheduled
	recs, err := store.Read(key, store.ReadLimit(1))
	if err != store.ErrNotFound || len(recs) > 0 {
		return errors.BadRequest("cron.schedule", "job already exists")
	}

	job := &pb.Job{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Interval:    req.Interval,
		Callback:    req.Callback,
	}

	// start the job
	cr := c.Setup(job)
	// save the job
	c.jobs[key] = cr

	rec := store.NewRecord(
		key,
		job,
	)

	rsp.Job = job

	// save in store
	return store.Write(rec)
}

func (c *Cron) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("cron.delete", "missing id")
	}

	tnt, _ := tenant.FromContext(ctx)
	key := fmt.Sprintf("%s/%s/%s", jobPrefix, tnt, req.Id)

	c.Lock()
	defer c.Unlock()

	if cr, ok := c.jobs[key]; ok {
		cr.Stop()
		delete(c.jobs, key)
	}

	// delete from store
	return store.Delete(key)
}

func (c *Cron) Jobs(ctx context.Context, req *pb.JobsRequest, rsp *pb.JobsResponse) error {
	tnt, _ := tenant.FromContext(ctx)
	key := fmt.Sprintf("%s/%s/", jobPrefix, tnt)

	recs, err := store.Read(key, store.ReadPrefix())
	if err != nil {
		return err
	}

	for _, rec := range recs {
		job := new(pb.Job)
		rec.Decode(&job)
		rsp.Jobs = append(rsp.Jobs, job)
	}

	return nil
}
