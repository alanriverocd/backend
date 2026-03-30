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

   // Primero intenta leer las credenciales desde la variable de entorno FIREBASE_CREDENTIALS_JSON
   credJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON")
   var opt option.ClientOption

   if credJSON != "" {
	   // Si existe, crea un archivo temporal
	   tmpFile, err := os.CreateTemp("", "firebase-creds-*.json")
	   if err != nil {
		   log.Fatalf("No se pudo crear archivo temporal para credenciales de Firebase: %v", err)
	   }
	   defer tmpFile.Close()

	   _, err = tmpFile.WriteString(credJSON)
	   if err != nil {
		   log.Fatalf("No se pudo escribir credenciales en archivo temporal: %v", err)
	   }
	   opt = option.WithCredentialsFile(tmpFile.Name())
   } else {
	   // Si no existe la variable, usa la ruta por defecto o la variable FIREBASE_CREDENTIALS_PATH
	   credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	   if credPath == "" {
		   credPath = "../finatiol-firebase-adminsdk-fbsvc-3e52add403.json"
	   }
	   opt = option.WithCredentialsFile(credPath)
   }

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
