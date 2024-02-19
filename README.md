**[Internal Document for Developers](https://docs.google.com/document/d/1oxb7pU3UWdIycvHNsNe2JztI_q0RyP08HDL1qKCBYwk/edit)**

## Get Started
I personally recommend using JetBrains' [Goland IDE](https://www.jetbrains.com/go/) for the project. \
We have [free educational licenses](https://www.jetbrains.com/community/education/#students) to the JetBrains' IDEs as Yale students. \
After Goland is installed, open this project in Goland and it will automatically install the necessary dependencies. 

## Start Development
### 1. Create .env files
Create three files: `.env`, `.env.local`, and `.env.dev` environments. \
The content of these files are in [this google doc](https://docs.google.com/document/d/1b8g1Iau7TJo6f2scI5bhIGdKSBgu8lu7Kp8v8QPR0Uc/edit) \
Please request access to the document if you are a developer of the project.

### 2. Install MySQL
Install a version no older than `8.0`.

### 3. Install Redis (Optional)
You don't need to install Redis if you don't need to work with the log-in/sign-up workflows.

### 4. Start the Application
#### Local environment: 
Run the `main.go` program.

#### Dev environment: 
Specify the environment variable `GIN_ENV_MODE=dev`.  
You can do this in Goland as:
![image](https://github.com/dekunma/cpsc-519-project-backend/assets/53892579/ad8ebce2-3e89-4fbc-91e7-a808d7414828)

## Debugging
### View logs of backend application
Download the `full_stack_project.pem` file and find the ip address of the AWS EC2 instance from our internal document provided at the top. \
SSH into the EC2 server: 
```bash
ssh -i full_stack_project.pem ec2-user@OUR_SERVER_IP
```

View the logs: 
```bash
sudo less /var/log/web.stdout.log
```
