package impl

import (
	"encoding/json"
	"github.com/blueleafdi2/moving-window/src/api"
	"github.com/blueleafdi2/moving-window/src/common"
	"github.com/blueleafdi2/moving-window/src/server"
	"github.com/blueleafdi2/moving-window/src/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	instance := new(serviceImpl)
	instance.load()
	go instance.probeAndPersist()
	server.InjectService(instance)
}

type serviceImpl struct {
	Buckets  [common.WindowSize]int64 `json:"buckets"`
	Total    int64                    `json:"total"`
	LastSave int64                    `json:"last_save"`
	Mu       sync.Mutex               `json:"-"`
}

func (s *serviceImpl) probeAndPersist() {
	probeTicker := time.NewTicker(common.ProbeInterval)
	for {
		select {
		case <-probeTicker.C:
			s.probe()
			s.persist()
		}
	}
}

func (s *serviceImpl) CountHandler(w http.ResponseWriter, r *http.Request) {
	err := util.Quietly(func() {
		s.incrementCount()
		response := api.Response{
			Status: "Success",
			Code:   http.StatusOK,
			Data:   api.CountDada{TotalRequest: atomic.LoadInt64(&s.Total)},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}, true)

	if err != nil {
		s.buildDefaultErrRsp(w, r)
	}
}

func (s *serviceImpl) ProbeHandler() {
	go s.probeAndPersist()
}
func (s *serviceImpl) buildDefaultErrRsp(w http.ResponseWriter, r *http.Request) {
	response := api.Response{
		Status: "Failed",
		Code:   http.StatusInternalServerError,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *serviceImpl) probe() {
	expiredIndex := s.getExpiredIndex()
	lastCount := atomic.LoadInt64(&s.Buckets[expiredIndex])
	if lastCount > 0 {
		atomic.AddInt64(&s.Total, -lastCount)
		atomic.StoreInt64(&s.Buckets[expiredIndex], 0)
	}
	// Update LastSave to the current time
	s.LastSave = time.Now().Unix()
}

func (s *serviceImpl) incrementCount() {
	unixTime := time.Now().Unix()
	index := unixTime % common.WindowSize

	// Increment the bucket and update the total
	atomic.AddInt64(&s.Buckets[index], 1)
	atomic.AddInt64(&s.Total, 1)
}

func (s *serviceImpl) persist() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	data, err := json.Marshal(s)
	util.CheckErr(err)

	ioutil.WriteFile(common.PersistenceFile, data, 0644)
	util.CheckErr(err)
}

func (s *serviceImpl) load() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	data, err := ioutil.ReadFile(common.PersistenceFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Failed to read persistence file: %v\n", err)
		}
		//return
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, s)
		util.CheckErr(err)
	}
	now := time.Now().Unix()
	gapSize := now - s.LastSave
	resetSize := gapSize
	if resetSize > common.WindowSize {
		resetSize = common.WindowSize
	}
	for i := int64(0); i < gapSize; i++ {
		expiredIndex := (s.LastSave + i + 1) % common.WindowSize
		lastCount := atomic.LoadInt64(&s.Buckets[expiredIndex])
		if lastCount > 0 {
			atomic.AddInt64(&s.Total, -lastCount)
			atomic.StoreInt64(&s.Buckets[expiredIndex], 0)
		}
	}
}

func (s *serviceImpl) getExpiredIndex() int64 {
	return (time.Now().Unix() + 1) % common.WindowSize
}
