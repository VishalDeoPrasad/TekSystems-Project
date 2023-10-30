package main

import (
	"context"
	"fmt"
	"golang/internal/auth"
	"golang/internal/database"
	"golang/internal/handlers"
	"golang/internal/services"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

func main() {
	err := startApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("hello this is our app")
}
func startApp() error {

	// ==========================================================================
	// Initialize authentication support
	log.Info().Msg("main : Started : Initializing authentication support")
	/*
		It reads the content of a file named "private.pem" and stores it in the privatePEM variable as a byte slice.
		If there is an error reading the file, it sets err to the error that occurred.
	*/
	privatePEM, err := os.ReadFile("private.pem")

	//This is an error check. If there was an error in reading the file
	if err != nil {
		return fmt.Errorf("reading auth private key %w", err)
		//it returns an error message using fmt.
		//The %w verb is used to wrap the original error within a new error message for context.
	}
	//PEM is a text-based encoding format that's commonly used for various types of data, including cryptographic keys.
	//converting it into a format that can be used for cryptographic operations, such as signing or decrypting data.
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing auth private key %w", err)
	}

	publicPEM, err := os.ReadFile("pubkey.pem")
	if err != nil {
		return fmt.Errorf("reading auth public key %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("parsing auth public key %w", err)
	}

	//received the instance of Auth struct and a nil error
	a, err := auth.NewAuth(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("constructing auth %w", err)
	}

	// =========================================================================
	// Start Database
	log.Info().Msg("main : Started : Initializing db support")
	db, err := database.Open() // @ Return a connection of db
	/*
		db: This is a pointer to a gorm.DB object. It represents a connection
		to your database. You can use this object to perform various database
		operations such as querying, inserting, updating, or deleting data in
		your database.*/
	if err != nil {
		return fmt.Errorf("connecting to db %w", err)
	}
	//pg : This variable will hold the pointer to the underlying PostgreSQL database connection.
	pg, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w ", err)
	}

	/*
		We use context in Go to control and manage how long operations should take
		and how they can be canceled or timed out. It helps ensure that operations
		don't run indefinitely and allows for better coordination in concurrent
		or networked applications. It's a way to handle timeouts, cancellations,
		and data sharing between different parts of a program.

		creating a context with a timeout. The operation should complete within
		the specified time, or it will be canceled.	*/
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	/*
		context.Background() : To create a basic context, start with the "background"
		context. This is typically the parent context from which other contexts are
		derived. Use context.Background() to create a background context.

		.WithTimeout(context.Background(), time.Second*5) :
		If you want to set a timeout for the context, use context.WithTimeout.
		This is helpful when you want to ensure that a function or operation doesn't
		take longer than a specified duration.
		context.WithTimeout(parentContext, timeout) creates a new context derived
		from a parent context with the specified timeout.

		The timeout specifies how long the context should be valid before it's
		automatically canceled. If the timeout is reached, the context becomes "done,"
		and any operations using this context should be canceled or stopped.*/
	/*
	   err = pg.PingContext(ctx) is checking the accessibility of a database by
	   sending a ping request with a specified context, and it captures any error
	   that might occur during the operation. This is a common practice to ensure
	   that the database is reachable before performing other database operations.*/
	err = pg.PingContext(ctx)
	/*
		ctx: A context can be used to set a timeout for the operation. If the ping
		operation takes longer than the context's timeout, it can be canceled, and
		an error will be returned.*/
	if err != nil {
		return fmt.Errorf("database is not connected: %w ", err)
	}

	// =========================================================================
	//Initialize Conn layer support
	ms, err := services.NewConn(db)
	if err != nil {
		return err
	}
	err = ms.AutoMigrate()
	if err != nil {
		return err
	}

	// Initialize http service
	api := http.Server{
		Addr:         ":8080",
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
		IdleTimeout:  800 * time.Second,
		Handler:      handlers.API(a, ms),
	}

	// channel to store any errors while setting up the service
	serverErrors := make(chan error, 1)
	go func() {
		log.Info().Str("port", api.Addr).Msg("main: API listening")
		serverErrors <- api.ListenAndServe()
	}()
	//shutdown channel intercepts ctrl+c signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error %w", err)
	case sig := <-shutdown:
		log.Info().Msgf("main: Start shutdown %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		//Shutdown gracefully shuts down the server without interrupting any active connections.
		//Shutdown works by first closing all open listeners, then closing all idle connections,
		//and then waiting indefinitely for connections to return to idle and then shut down.
		err := api.Shutdown(ctx)
		if err != nil {
			//Close immediately closes all active net.Listeners
			err = api.Close() // forcing shutdown
			return fmt.Errorf("could not stop server gracefully %w", err)
		}

	}
	return nil

}
