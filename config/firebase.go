package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	FirebaseApp     *firebase.App
	AuthClient      *auth.Client
	FirestoreClient *firestore.Client
)

func InitFirebase() {
	ctx := context.Background()

	credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credPath == "" {
		credPath = "../finatiol-firebase-adminsdk-fbsvc-3e52add403.json"
	}

	opt := option.WithCredentialsFile(credPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error al inicializar Firebase: %v", err)
	}

	FirebaseApp = app

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error al inicializar Firebase Auth: %v", err)
	}
	AuthClient = authClient

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error al inicializar Firestore: %v", err)
	}
	FirestoreClient = firestoreClient

	log.Println("Firebase inicializado correctamente para proyecto: finatiol")
}
