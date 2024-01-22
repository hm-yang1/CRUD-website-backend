**Go api backend for a web forum**

Link: https://cvwo-assignment-backend.onrender.com/ 
(Might not work with other domains due to CORS)

Features:
1. Account-based authentication
2. Hash password with bcrypt and stores sessions with gorilla sessions
3. CRUD posts and comments
4. middleware to check authentication
5. Upvotes for posts and comments
6. Filter posts by tags
7. Search posts
8. Sort posts and comments by time/upvotes

Requirements:
1. Go 1.21.4
2. PostgreSQL 15 server
    a. Schema named cvwo
    b. Write the DB connection url
3. .env file for sensitive information (See sample below)

Sample Env:

    # Below are neccessary for the code to run, pls don't change the names. 
    
    # Database variables
    DB_URL = your_db_url
    
    # Sessions variables
    SESSION_SECRET = your_key
    JWT_SECRET = your_key
    
    #Web
    PORT = your_own_port
    FRONTEND_URL = your_frontend_url
