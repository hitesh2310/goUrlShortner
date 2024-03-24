# goUrlShortner

A URL shortener written in Go.

## Requirements

- Go 1.16 or higher
- MySQL database
- Redis service  



## Installation

1. Clone the repository:

```bash
git clone https://github.com/hitesh2310/goUrlShortner.git
```

2. Install dependencies:

```bash
cd goUrlShortner
go mod download
```

3. Create a MySQL database and a user with privileges to create and modify tables.
   mysql> use urlShortner 
Database changed
mysql> show tables;
+-----------------------+
| Tables_in_urlShortner |
+-----------------------+
| host_stats            |
| url_mapping           |
+-----------------------+
2 rows in set (0.07 sec)

mysql> desc url_mapping ;
+-----------+--------------+------+-----+---------+----------------+
| Field     | Type         | Null | Key | Default | Extra          |
+-----------+--------------+------+-----+---------+----------------+
| id        | int          | NO   | PRI | NULL    | auto_increment |
| longUrl   | varchar(700) | YES  | UNI | NULL    |                |
| shortUrl  | varchar(15)  | YES  | UNI | NULL    |                |
| createdAt | bigint       | YES  |     | NULL    |                |
+-----------+--------------+------+-----+---------+----------------+
4 rows in set (0.03 sec)

mysql> show create table  url_mapping \G;
*************************** 1. row ***************************
       Table: url_mapping
Create Table: CREATE TABLE `url_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `longUrl` varchar(700) DEFAULT NULL,
  `shortUrl` varchar(15) DEFAULT NULL,
  `createdAt` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `longUrl` (`longUrl`),
  UNIQUE KEY `shortUrl` (`shortUrl`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
1 row in set (0.00 sec)




mysql> desc host_stats ;
+-------+--------------+------+-----+---------+----------------+
| Field | Type         | Null | Key | Default | Extra          |
+-------+--------------+------+-----+---------+----------------+
| id    | int          | NO   | PRI | NULL    | auto_increment |
| host  | varchar(100) | YES  | UNI | NULL    |                |
| count | int          | YES  |     | NULL    |                |
+-------+--------------+------+-----+---------+----------------+
3 rows in set (0.00 sec)


mysql> show create table  host_stats \G;
*************************** 1. row ***************************
       Table: host_stats
Create Table: CREATE TABLE `host_stats` (
  `id` int NOT NULL AUTO_INCREMENT,
  `host` varchar(100) DEFAULT NULL,
  `count` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `host` (`host`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
1 row in set (0.00 sec)




4. Update the `config/config.json` file with your MySQL and Redis database credentials.


5. Start the application:

```bash
go run main.go
```

The server will be available at http://localhost:8081.

## Usage

You can use the API to create short URLs or redirect to the original long URLs.

### Creating Short URLs

 To create a short URL, make a `POST` request to the `/shorten` endpoint with the following parameters:
 - `url`: The long URL you want to shorten
   Example -  
       curl --location 'localhost:8081/shorten' --header 'Content-Type: application/json' --data '{
         "url": "https://app.slack.com/client/T046XDVMU49/D048U24R0JG"}' 
       Response - {"shortUrl": "localhost:8081/O"}


### Metrics API
 
  To return top 3 domains that has been shortened most number of times, make a `GET` request to the `/stat`     endpoint with the following parameters:
   Example - 
    curl --location 'localhost:8081/stat' 
    Response - {"result": ["developers.facebook.com","app.slack.com","www.iplt20.com"]} 


#### Docker Deployment
 - Step 1: Load the provided Docker image onto your system
   `docker load < urlShortenerImage.gz`

 - Step 2: Run a Docker container with the loaded image
   `docker run --name <name-the-container> -d --net=host -v  <log-path-on-host>:/app/logs/ urlShortenerImage`
   Replace <name-the-container> with a suitable name for your container
   Replace <log-path-on-host> with a suitable log path for your container
   Note: This command assumes that MySQL, Redis, and configuration files are correctly set up
   Make sure your MySQL and Redis instances are running and configuration file is properly configured for the URL shortener application
   In case, required to adjust the Docker run command as needed to accommodate your specific environment and configuration requirements

  









































