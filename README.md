## First Commit

1. Ran: "go get -u github.com/gofiber/fiber/v2 go.mongodb.org/mongo-driver/mongo github.com/joho/godotenv github.com/go-playground/validator/v10"
2. Set up project and database on mongoDB atlas.
3. Got mongo connection uri from db atlas and created an env var in .env.

## Second Commit

1. Created a function to get the MongoURI from the env.
2. Created a function to connect to the database.
3. Created a function to get collections in the database.

## Third Commit

1. Created a user model.
2. Created a function that will handle user routes.
3. Created a user response struct.

## Fourth Commit

1. Created a route to create a user.