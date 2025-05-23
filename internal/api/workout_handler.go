package api

import (
	"encoding/json"
	"fmt"
	"goBackendServer/internal/store"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workOutStore store.WorkoutStore
}

func NewWorkoutHandler(workOutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workOutStore: workOutStore,
	}
}

func (wh *WorkoutHandler) HandleWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workOutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var workOut *store.Workout

	workOut, err = wh.workOutStore.GetWorkoutById(workOutId)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workOut)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workOut store.Workout
	err := json.NewDecoder(r.Body).Decode(&workOut)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}

	savedWorkout, err := wh.workOutStore.CreateWorkout(&workOut)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedWorkout)
}
