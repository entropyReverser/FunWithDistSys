package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func noDatabasePool(numberOfQueries int) {
	var wg sync.WaitGroup

	for i := 0; i < numberOfQueries; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
				os.Exit(1)
			}
			defer conn.Close(context.Background())

			fmt.Printf("Connection number %d established\n", i)

			rows, err := conn.Query(context.Background(), "select pg_sleep(2),'Hello, world!'")
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(1)
			}
			defer rows.Close()

			for rows.Next() {
				values, err := rows.Values()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get row values: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(values)
			}
			if err := rows.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Row iteration error: %v\n", err)
				os.Exit(1)
			}
		}(i)
	}
	wg.Wait()
}

func with_library_connection_pool(numberOfQueries int) {
	connPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connection pool max connections:", connPool.Config().MaxConns)

	var wg sync.WaitGroup

	for i := 0; i < numberOfQueries; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			conn, err := connPool.Acquire(context.Background())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to acquire connection: %v\n", err)
				os.Exit(1)
			}
			defer conn.Release()

			fmt.Printf("Connection number %d established\n", i)

			rows, err := conn.Query(context.Background(), "select pg_sleep(1),'Hello, world!'")
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(1)
			}
			defer rows.Close()

			for rows.Next() {
				values, err := rows.Values()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get row values: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(values)
			}
			if err := rows.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Row iteration error: %v\n", err)
				os.Exit(1)
			}
		}(i)
	}
	wg.Wait()
}

func custom_simple_connection_pool(numberOfQueries int) {
	//Create connection pool using pg.connect
	numberOfConnection := 12 //Matching the number of connection to phxpool max connection

	connPool := make(chan *pgx.Conn, numberOfConnection)
	for i := 0; i < numberOfConnection; i++ {
		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		connPool <- conn
	}

	var wg sync.WaitGroup

	for i := 0; i < numberOfQueries; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			conn := <-connPool
			defer func() {
				connPool <- conn
			}()

			fmt.Printf("Connection number %d established\n", i)

			rows, err := conn.Query(context.Background(), "select pg_sleep(1),'Hello, world!'")
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(1)
			}
			defer rows.Close()

			for rows.Next() {
				values, err := rows.Values()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get row values: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(values)
			}
			if err := rows.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Row iteration error: %v\n", err)
				os.Exit(1)
			}
		}(i)
	}
	wg.Wait()

}

func main() {
	defer timer("main")()
	//noDatabasePool(100)
	//with_library_connection_pool(201)
	custom_simple_connection_pool(201)
}
