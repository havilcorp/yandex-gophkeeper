CREATE TABLE IF NOT EXISTS users(
    ID serial PRIMARY KEY,
    password VARCHAR (255) NOT NULL,
    email VARCHAR (300) UNIQUE NOT NULL
)