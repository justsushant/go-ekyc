package worker

import (
	"encoding/json"
	"log"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type Worker struct {
	queue       service.TaskQueue
	dStore      store.DataStore
	fStore      store.FileStore
	faceMatcher service.FaceMatcher
	ocr         service.OCRPerformer
}

type QueueMessage struct {
	Type string          `json:"type"`
	Msg  json.RawMessage `json:"msg"`
}

func NewWorker(queue service.TaskQueue, dStore store.DataStore, fStore store.FileStore, faceMatcher service.FaceMatcher, ocr service.OCRPerformer) *Worker {
	return &Worker{
		queue,
		dStore,
		fStore,
		faceMatcher,
		ocr,
	}
}

func (w *Worker) ProcessMessages() {
	msgs, err := w.queue.PullJobFromQueue()
	if err != nil {
		log.Fatalf("Error while pulling jobs from queue: %v", err)
	}

	for payload := range msgs {
		// log.Printf("Received payload: %+v\n", payload)
		// log.Printf("payload.Body: %+v", string(payload.Body))

		// var msg json.RawMessage
		q := QueueMessage{
			// Msg: msg,
		}

		if err := json.Unmarshal(payload.Body, &q); err != nil {
			log.Println("Error while unmarshaling JSON: ", err)
		}

		// log.Printf("msg: %+v\n", msg)
		// log.Printf("q: %+v\n", q)

		switch q.Type {
		case types.FaceMatchWorkType:
			var s types.FaceMatchInternalPayload
			if err := json.Unmarshal(q.Msg, &s); err != nil {
				log.Println(err)
			}

			w.ProcessFaceMatchWork(s)
		case types.OCRWorkType:
			var s types.OCRInternalPayload
			if err := json.Unmarshal(q.Msg, &s); err != nil {
				log.Println(err)
			}

			w.ProcessOCRWork(s)
		}
	}
}

func (w *Worker) ProcessFaceMatchWork(payload types.FaceMatchInternalPayload) error {
	// log.Println("INSIDE THE HELPER")
	// log.Printf("%v\n", payload)

	// change state to processing
	err := w.dStore.UpdateFaceMatchJobProcessed(payload.JobID)
	if err != nil {
		log.Println("Error while marking the face match job processing: ", err)
		err := w.dStore.UpdateFaceMatchJobFailed(payload.JobID, err.Error())
		if err != nil {
			log.Println("Error while marking the face match jo failed: ", err.Error())
		}
	}

	// do the work
	p := types.FaceMatchPayload{
		Image1: payload.Image1,
		Image2: payload.Image2,
	}
	score, err := w.faceMatcher.CalcFaceMatchScore(p)
	if err != nil {
		log.Println("Error while performing the face match: ", err)
		err := w.dStore.UpdateFaceMatchJobFailed(payload.JobID, err.Error())
		if err != nil {
			log.Println("Error while marking the face match jo failed: ", err.Error())
		}
	}

	// update in db
	err = w.dStore.UpdateFaceMatchJobCompleted(payload.JobID, score)
	if err != nil {
		log.Println("Error while updating the face match result: ", err)
		err := w.dStore.UpdateFaceMatchJobFailed(payload.JobID, err.Error())
		if err != nil {
			log.Println("Error while marking the face match jo failed: ", err.Error())
		}
	}

	return nil
}

func (w *Worker) ProcessOCRWork(payload types.OCRInternalPayload) error {
	// change state to processing
	err := w.dStore.UpdateOCRJobProcessed(payload.JobID)
	if err != nil {
		log.Println("Error while marking the ocr job processing: ", err)
		err := w.dStore.UpdateOCRJobFailed(payload.JobID, err.Error())
		if err != nil {
			log.Println("Error while marking the ocr job failed: ", err)
		}
	}

	// do the work
	p := types.OCRPayload{
		Image: payload.Image,
	}
	resp, err := w.ocr.PerformOCR(p)
	if err != nil {
		log.Println("Error while performing the ocr: ", err)
		err := w.dStore.UpdateOCRJobFailed(payload.JobID, err.Error())
		if err != nil {
			log.Println("Error while marking the ocr job failed: ", err)
		}
	}

	// update in db
	err = w.dStore.UpdateOCRJobCompleted(payload.JobID, resp)
	if err != nil {
		log.Println("Error while updating the face match result: ", err)
		err := w.dStore.UpdateOCRJobFailed(payload.JobID, err.Error())
		if err != nil {
			log.Println("Error while marking the ocr job failed: ", err)
		}
	}

	return nil
}
