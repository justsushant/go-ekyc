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
	dStore      store.WorkerDataStore
	fStore      store.FileStore
	faceMatcher service.FaceMatcher
	ocr         service.OCRPerformer
}

type QueueMessage struct {
	Type string          `json:"type"`
	Msg  json.RawMessage `json:"msg"`
}

func NewWorker(queue service.TaskQueue, dStore store.WorkerDataStore, fStore store.FileStore, faceMatcher service.FaceMatcher, ocr service.OCRPerformer) *Worker {
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
		// unmarshal the payload
		q := QueueMessage{}
		if err := json.Unmarshal(payload.Body, &q); err != nil {
			log.Println("Error while unmarshaling JSON: ", err.Error())
		}

		// call service on the basis of type in payload
		switch q.Type {
		case types.FACE_MATCH_WORK_TYPE:
			var s types.FaceMatchInternalPayload
			if err := json.Unmarshal(q.Msg, &s); err != nil {
				log.Println(err)
			}

			w.ProcessFaceMatchWork(s)
		case types.OCR_WORK_TYPE:
			var s types.OCRInternalPayload
			if err := json.Unmarshal(q.Msg, &s); err != nil {
				log.Println(err)
			}

			w.ProcessOCRWork(s)
		}
	}
}

func (w *Worker) ProcessFaceMatchWork(payload types.FaceMatchInternalPayload) {
	// change state to processing
	err := w.dStore.UpdateFaceMatchJobProcessed(payload.JobID)
	if err != nil {
		log.Printf("Error while updating the face match job (%s) state to 'processing': %s\n", payload.JobID, err.Error())
		w.changeStateToFailed(types.FACE_MATCH_WORK_TYPE, payload.JobID, err.Error())
		return
	}

	// do the work
	p := types.FaceMatchPayload{
		Image1: payload.Image1,
		Image2: payload.Image2,
	}
	score, err := w.faceMatcher.PerformFaceMatch(p)
	if err != nil {
		log.Printf("Error while performing the face match job (%s): %s\n", payload.JobID, err.Error())
		w.changeStateToFailed(types.FACE_MATCH_WORK_TYPE, payload.JobID, err.Error())
		return
	}

	// change state to completed
	err = w.dStore.UpdateFaceMatchJobCompleted(payload.JobID, score)
	if err != nil {
		log.Printf("Error while updating the face match job (%s) state to 'completed': %s\n", payload.JobID, err.Error())
		w.changeStateToFailed(types.FACE_MATCH_WORK_TYPE, payload.JobID, err.Error())
		return
	}
}

func (w *Worker) ProcessOCRWork(payload types.OCRInternalPayload) {
	// change state to processing
	err := w.dStore.UpdateOCRJobProcessed(payload.JobID)
	if err != nil {
		log.Printf("Error while updating the ocr job (%s) state to 'processing': %s\n", payload.JobID, err.Error())
		w.changeStateToFailed(types.OCR_WORK_TYPE, payload.JobID, err.Error())
		return
	}

	// do the work
	p := types.OCRPayload{
		Image: payload.Image,
	}
	resp, err := w.ocr.PerformOCR(p)
	if err != nil {
		log.Printf("Error while performing the ocr job (%s): %s\n", payload.JobID, err.Error())
		w.changeStateToFailed(types.OCR_WORK_TYPE, payload.JobID, err.Error())
		return
	}

	// change state to completed
	err = w.dStore.UpdateOCRJobCompleted(payload.JobID, resp)
	if err != nil {
		log.Printf("Error while updating the face match job (%s) state to 'completed': %s\n", payload.JobID, err.Error())
		w.changeStateToFailed(types.OCR_WORK_TYPE, payload.JobID, err.Error())
		return
	}
}

func (w *Worker) changeStateToFailed(jobType, jobID, errMessage string) {
	switch jobType {
	case types.FACE_MATCH_WORK_TYPE:
		err := w.dStore.UpdateFaceMatchJobFailed(jobID, errMessage)
		if err != nil {
			log.Printf("Error while updating the face match job (%s) state to 'failed': %s\n", jobID, errMessage)
		}
		return
	case types.OCR_WORK_TYPE:
		err := w.dStore.UpdateOCRJobFailed(jobID, errMessage)
		if err != nil {
			log.Printf("Error while updating the ocr job (%s) state to 'failed': %s\n", jobID, errMessage)
		}
	}
}
