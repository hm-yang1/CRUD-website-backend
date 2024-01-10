**Go api backend for a web forum**

Features:
1. Account-based authentication
   a. Hash password with bcrypt and stores sessions with gorilla sessions
2. CRUD posts and comments
3. middleware to check authentication
4. Upvotes for posts and comments
5. Filter posts by tags
6. Search posts
7. Sort posts and comments by time/upvotes

Requirements:
1. Go 1.21.4
2. MySql 8.2.0
3. .env file for sensitive information (See sample below)

Sample Env:

    # Below are absolutely neccessary for the code to run, pls don't change the names. 
    #If u want to add more, u probably know more than me so feel free ig.
    
    # Database variables
    DB_URL = your_db_url
    
    # Sessions variables
    SESSION_SECRET = your_key
    JWT_SECRET = your_key
    
    #Web
    PORT = your_own_port
