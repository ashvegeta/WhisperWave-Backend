<h1> WhisperWave</h1>
<h3> A simple chat application written in Golang</h3>

<h2>About</h2>

WhisperWave is an application that emerged as a hobby project to understand and simulate real-time chat applications. To understand in detail, the project structure and data flow in the design and the thinking behind the database design, click on the medium 
[post](https://medium.com/@1ms18cs030/my-experience-in-building-a-chat-application-in-golang-f0b815d7b7ae).

<h2>System Design</h2>
  
![2-user system design](https://i.imgur.com/uTkLGRM.png)

<h2>DataBase Design</h2>
  
![Database Design Overview](https://i.imgur.com/tDzaQdw.png)

<h2>Tech Stack</h2>

<img src="https://i.imgur.com/rZsHj24.png" width="98px" height="48px">   <img src="https://i.imgur.com/OAOXf5W.png" width="200px" height="100px">   <img src="https://i.imgur.com/ZgxcU74.png" width="250px" height="48px">

<h2>Project Setup</h2>

1. To set up the project, clone the repository by typing:
   
    ```
   git clone
    ```

4. Install all the golang modules declared in ```go.mod``` file using:

   ```
   go mod download
   ```

6. Follow the official RabbitMQ documentation or visit dockerhub to download the server image:

   ```
   docker pull rabbitmq:3.12-management
   ```
   
8. Run the RabbitMQ docker container and map it to the port of your liking:

     ```shell
   docker run -d --hostname [your-hostname] --name [docker-server-name] -p [service-port]:5672 -p [management-port]:15672 rabbitmq:3.12-management
     ```

9. Enter the root directory of the project where the MakeFile is present and type:
    ```
    make
    ```
