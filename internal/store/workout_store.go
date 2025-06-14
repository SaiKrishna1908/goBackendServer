// Repository Layer
package store

import (
	"database/sql"
	"fmt"
)

type Workout struct {
	ID              int            `json:"id"`
	UserId          int            `json:"user_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight" `
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutById(int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkoutById(int64) error
	GetWorkoutOwner(id int64) (int, error)
}

// TODO: Generate id by auto increment
func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query :=
		`
	INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id
	`

	err = tx.QueryRow(query, workout.UserId, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	// insert the entries

	for _, entry := range workout.Entries {
		query := `
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`

		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)

		if err != nil {
			return nil, err
		}

	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	query := `
		SELECT id, title, description, duration_minutes, calories_burned
		FROM workouts
		WHERE id = $1
	`
	// get the entry

	var workOut Workout

	err = tx.QueryRow(query, id).Scan(&workOut.ID, &workOut.Title, &workOut.Description, &workOut.DurationMinutes, &workOut.CaloriesBurned)

	if err == sql.ErrNoRows {
		fmt.Println(err)
		return nil, nil
	}

	if err != nil {
		fmt.Println("Error fetching workouts", err)
		return nil, nil
	}

	entryQuery := `
			SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
			FROM workout_entries
			WHERE workout_id = $1
			ORDER BY order_index
			`

	rows, err := pg.db.Query(entryQuery, workOut.ID)

	if err != nil {
		fmt.Println("Error fetching workout entries")
		return &workOut, nil
	}

	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry

		err := rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex)

		if err != nil {
			return nil, err
		}

		workOut.Entries = append(workOut.Entries, entry)
	}

	// fmt.Println("workout from database ", workOut)

	return &workOut, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workOut *Workout) error {
	tx, err := pg.db.Begin()

	if err != nil {
		fmt.Println("Error with db Begin ", err)
	}

	query := `
		UPDATE workouts
		SET title = $2, description = $3, duration_minutes = $4, calories_burned = $5
		where id = $1
	`

	result, err := tx.Exec(query, workOut.ID, workOut.Title, workOut.Description, workOut.DurationMinutes, workOut.CaloriesBurned)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// patch update workout_entries

	for _, entry := range workOut.Entries {
		query =
			`
		 INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight ,notes, order_index)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)

		`

		_, err := tx.Exec(query,
			workOut.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pg *PostgresWorkoutStore) DeleteWorkoutById(id int64) error {

	workout, err := pg.GetWorkoutById(id)

	if err != nil {
		fmt.Println("error fetching workout for delete ", err)
		return err
	}

	query := `
		DELETE from workouts
		where id = $1
	`

	result, err := pg.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsEffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsEffected == 0 {
		return sql.ErrNoRows
	}

	// delete workout entries
	fmt.Println("Deleting from workout_entries")

	entryDeleteQuery := `
		DELETE FROM workout_entries WHERE workout_id = $1
	`

	_, err = pg.db.Exec(entryDeleteQuery, workout.ID)

	if err != nil {
		fmt.Println("Error deleting tables from` workout_entries")
	}

	return nil
}

func (pg *PostgresWorkoutStore) GetWorkoutOwner(workoutID int64) (int, error) {
	var userID int

	query := `
		SELECT user_id 
		FROM workouts
		WHERE id = $1
	`

	err := pg.db.QueryRow(query, workoutID).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, err
}
