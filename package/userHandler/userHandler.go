package userhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Rekciuq/go-bucket/ent"
)

type postPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func CreateUser(ctx context.Context, client *ent.Client, payload postPayload) (*ent.User, error) {
	u, err := client.User.Create().SetAge(payload.Age).SetName(payload.Name).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	log.Println("user was created: ", u)

	return u, nil
}

func QueryAllUsers(ctx context.Context, client *ent.Client) ([]*ent.User, error) {
	u, err := client.User.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed querying all users: %w", err)
	}
	log.Println("all users: ", u)
	return u, nil
}

func DeleteUser(ctx context.Context, client *ent.Client, userId int) error {
	err := client.User.DeleteOneID(userId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete user: %w", err)
	}

	return nil
}

func UserHandler(ctx context.Context, client *ent.Client) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		res, err := QueryAllUsers(ctx, client)
		if err != nil {
			fmt.Fprintln(w, "It's an error!", err)
		}

		fmt.Fprintln(w, "Getted the user")
		fmt.Fprintln(w, "Users is: ", res)
		fmt.Println("it's get!")
	})

	router.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var payload postPayload
		fmt.Println("It's post!")
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		usr, err := CreateUser(ctx, client, payload)
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			http.Error(w, "User wans't created", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Response!")
		fmt.Fprintln(w, usr)
	})
	router.HandleFunc("DELETE /{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		log.Printf("Attempting to delete user with ID: %d", id)
		error := DeleteUser(ctx, client, id)
		if error != nil {
			log.Println("Eror:", error.Error())
			http.Error(w, "User with this id doesn't exists", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User deleted successfully."))
	})

	return router
}
