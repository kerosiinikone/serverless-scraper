package reddit

import (
	"fmt"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	s3storage "github.com/kerosiinikone/serverless-scraper/infra/blob"
	"github.com/kerosiinikone/serverless-scraper/internal/scraper"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
)

type Manager struct {
	finished chan<- struct{}
	
	request  *scraper.APIRequest

	timer 	*time.Timer
	timerStopped bool

	storage *s3manager.Uploader
}

func NewManager(f chan<- struct{}, req *scraper.APIRequest, storage *s3manager.Uploader) actor.Producer {
	return func() actor.Receiver {
		return &Manager{
			request: req,
			storage: storage,
			finished: f,
		}
	}
}

func (m *Manager) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case models.RedditPostDetails:
		ctx.SpawnChild(NewActor(msg, ctx.PID()), fmt.Sprintf("scraper-%s", msg.Id))
	case models.ForumTree:
		if msg.Id != "" {
			if err := m.storePost(msg); err != nil {
				close(m.finished)
				ctx.Engine().Poison(ctx.PID())
			}
		}
		m.resetTimer(func() error {
			close(m.finished)
			ctx.Engine().Poison(ctx.PID())
			return nil
		})
	}
}

func (m *Manager) resetTimer(fn func() error) {
	if m.timerStopped {
		return
	}
	if m.timer != nil {
		m.timer.Stop()
	}
	m.timer = time.AfterFunc(8*time.Second, func() {
		m.timerStopped = true 
		if err := fn(); err != nil {
			fmt.Println("Failed to reset timer:", err)
		}
	})
}

// Mock the storage in case of tests -> otherwise the manager will try to save the file to S3 and fail
func (m *Manager) storePost(msg models.ForumTree) error {
	d := models.DataEntry{
		Post:      msg,
		ClientID:  m.request.ClientID,
		RequestID: m.request.ID,
	}
	if err := s3storage.SaveFile(m.storage, d); err != nil {
		fmt.Println("Failed to save file:", err)
		return err
	}
	return nil
}
