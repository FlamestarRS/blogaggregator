Blog Aggregator

-- SETUP AND INSTALLATION --

This section will take you through command line instructions required to run the Blog Aggregator

Install Go:
    - curl -sS https://webi.sh/golang | sh


Install Postgres for Linux/WSL:
    - sudo apt update
    - sudo apt install postgresql postgresql-contrib

Confirm installation with:
    - psql --version

Start postgres and set a :
    - sudo -u postgres psql

Within postgres:
    - CREATE DATABASE gator;
    - \c gator
    - ALTER USER postgres PASSWORD <password>;
    - exit

Create a config file in your home directory:
    - cd
    - touch .gatorconfig.json

In .gatorconfig.json:
    {"db_url": "postgres://username:@localhost:5432/database?sslmode=disable"}

From the project directory:
    - go build .


-- USAGE --

Commands are run from the command line while in the project directory.
All commands must be prefixed with:
    - blogaggregator

Examples:
    - blogaggregator register Flamestar
    - blogaggregator addfeed "feednamehere" "rsslink.here"
    - blogaggregator feeds

Commands:
    register:
        Usage: "register <name>"
        Desc: creates a user
    login:
        Usage: "login <name>"
        Desc: set current user, must be registered with "register"
    users:
        Usage: "users"
        Desc: print list of registered users
    addfeed:
        Usage: "addfeed <feed_name> <feed_url>"
        Desc: adds and follows a feed
    feeds:
        Usage: "feeds"
        Desc: print list of added feeds
    follow:
        Usage: "follow <feed_url>"
        Desc: adds an existing feed to the current user's following list
    following:
        Usage: "following"
        Desc: print list of all feeds the current user is following
    unfollow:
        Usage: "unfollow <feed_url>"
        Desc: removes a feed from the current user's following list
    agg:
        Usage: "agg <time_between_requests>" Ex: agg 30s | agg 1m | agg 1h
        Desc: scrapes added feeds for new posts
        Warning: agg will continue to make requests on the given interval until manually cancelled with ctrl C. It is intended to run indefinitely on a separete terminal. Be considerate of the third party servers you are requesting data from when choosing your interval.
    browse:
        Usage: "browse <number_of_posts>" defaults to 2 if no input is given
        Desc: prints info for the most recent posts, limited by the input
    reset:
        Usage: "reset"
        Desc: deletes all registered users and feeds

    
