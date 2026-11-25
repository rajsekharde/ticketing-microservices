package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
	"google.golang.org/grpc"
)

var db *database

type server struct {
	userpb.UnimplementedUserServiceServer
}

var cfg config

func main() {
	cfg = *loadEnv(".env")

	var err error
	db, err = newDatabase(&cfg)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	} else {
		log.Printf("Connected to database: %v\n", cfg.db_name)
	}
	defer db.close()

	// create tables
	err = db.initTables()
	if err != nil {
		log.Fatalf("Failed to initialize DB tables: %v", err)
	}
	log.Println("Initialized DB tables")

	// create a default admin user
	adminPHash, err := hashPassword(cfg.admin_password)
	if err != nil {
		log.Println("Failed to hash password for default admin user")
		adminPHash = ""
	}
	adminId, err := db.createUserQuery(&dbUser{
		name: cfg.admin_name,
		email: cfg.admin_email,
		password_hash: adminPHash,
		role: "ADMIN",
	})
	if err != nil {
		log.Println("Failed to create default admin user, error: %v\n", err)
	} else {
		log.Printf("Created default admin user: id = %v\n", adminId)
	}

	addr := ":" + cfg.port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &server{})

	log.Printf("User gRPC server listening on: %v", cfg.port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type config struct {
	port string

	db_host string
	db_port string
	db_user string
	db_password string
	db_name string

	admin_name string
	admin_email string
	admin_password string

	jwt_secret_key string
	access_token_expiry_mins int
	jwtIssuer string
}

func loadEnv(path string) *config {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("Failed to load env file. Using environment variables.")
	}

	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "50051"
	}

	mins, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINS"))
	if err != nil {
		mins = 15
	}

	return &config{
		port: port,

		db_host: os.Getenv("POSTGRES_HOST"),
		db_port: os.Getenv("POSTGRES_PORT"),
		db_user: os.Getenv("POSTGRES_USER"),
		db_password: os.Getenv("POSTGRES_PASSWORD"),
		db_name: os.Getenv("USER_DB"),

		admin_name: os.Getenv("USER_ADMIN_NAME"),
		admin_email: os.Getenv("USER_ADMIN_EMAIL"),
		admin_password: os.Getenv("USER_ADMIN_PASSWORD"),

		jwt_secret_key: os.Getenv("JWT_SECRET_KEY"),
		access_token_expiry_mins: mins,
		jwtIssuer: os.Getenv("JWT_ISSUER"),
	}
}