// Service Layer
package api

import (
	"database/sql"
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

func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside Handle Wokrout")
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	fmt.Println("id is ", paramsWorkoutId)
	workOutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)

	fmt.Println("int id is ", workOutId)
	if err != nil {
		http.NotFound(w, r)
	}

	workout, err := wh.workOutStore.GetWorkoutById(workOutId)

	fmt.Printf("%v", workout)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to fetch the workout", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(workout)

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

func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")

	if paramsWorkoutId == "" {
		fmt.Println("Could not extract id")
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	existingWorkout, err := wh.workOutStore.GetWorkoutById(workoutId)

	if err != nil {
		http.Error(w, "failed to fetch  workout", http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	// at this point we have a existing workout
	var updatedWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duraiton_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedWorkoutRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedWorkoutRequest.Title != nil {
		existingWorkout.Title = *updatedWorkoutRequest.Title
	}

	if updatedWorkoutRequest.Description != nil {
		existingWorkout.Description = *updatedWorkoutRequest.Description
	}

	if updatedWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updatedWorkoutRequest.CaloriesBurned
	}

	if updatedWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updatedWorkoutRequest.DurationMinutes
	}

	if updatedWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updatedWorkoutRequest.Entries
	}

	err = wh.workOutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		fmt.Println("update workout error", err)
		http.Error(w, "failed to update the workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(existingWorkout)
}

func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutParamId := chi.URLParam(r, "id")

	workoutId, err := strconv.ParseInt(workoutParamId, 10, 64)

	if err != nil {
		fmt.Println("Could not covert to int, ", workoutParamId)
		http.NotFound(w, r)
		return
	}

	err = wh.workOutStore.DeleteWorkoutById(workoutId)

	if err != nil {
		fmt.Println("Could not delete workout")
		http.Error(w, "failed to delete workout", http.StatusInternalServerError)
		return
	}

	if err == sql.ErrNoRows {
		fmt.Println("No workouts deleted")
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
