
# Simple Database pool demo

Without a database pool, each request to the database will open a new connection, execute the query, and close the connection.  
Which is ineffecient due to the overhead of opening and closing connections and could also lead to errors if the database is not able to handle the number of connections.

By running queries wihout a database pool on my local postgres databse, I was able to make 100 connections before I started getting errors. This is because the default number of connections in postgres is 100.    


<img src="https://github.com/user-attachments/assets/94d61b99-364f-4086-bd1d-2ffb0d2a707d" height="300">  

The image above shows the successful query with 100 connections.  
  
<img src="https://github.com/user-attachments/assets/60ebb34c-da24-4889-b957-41287db8eb06" height="300">  

The image above shows the error when I tried to make the 101th connection.  

Next I used pgx library to create a database pool. With which we are able to make multiple queries wihtout exausting the number of connections.

<img src="https://github.com/user-attachments/assets/a2677591-85c1-4904-bb98-aa7d3ba762d9" height="300">  

The image above shows the successful execution of 201 queries with pgx connection pool.

Next I created my own database pool using a go channel and putting a collection of connections in the channel. I matched the number of connection used by pgx and was able to match the execution time of pgx connection pool.

<img src="https://github.com/user-attachments/assets/f74e483d-48fc-40fe-ab11-af477186adc9" height="300">  

The image above shows the successful execution of 201 queries with custom connection pool using go channels.  
