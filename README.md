# Chatbox Project

Welcome to the Chatbox project! This is a chat application built with Golang using the Echo framework, structured to follow clean architecture principles and ensure modular design.

## Table of Contents

- [Project Structure](#project-structure)
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)

## Project Structure

```plaintext
├───api
│   ├───controller
│   │   ├───log        # Handles logging
│   │   ├───message    # Manages messages
│   │   ├───room       # Manages chat rooms
│   │   └───user       # Manages user operations
│   ├───middlewares    # Custom middlewares
│   └───router         # Routes definition
├───assets             # Static assets
├───bootstrap          # Initialization and setup
├───cmd                # Command line scripts
├───config             # Configuration files
├───domain             # Core business logic
├───infrastructure     # Infrastructure-related code
├───pkg
│   ├───cache          # Caching utilities
│   ├───cloudinary     # Cloudinary integration
│   ├───const          # Constants
│   ├───cron           # Scheduled tasks
│   ├───helper         # Helper functions
│   ├───jwt            # JWT authentication
│   ├───mail           # Email services
│   ├───oauth2         # OAuth2 services
│   ├───review         # Review and rating system
│   └───websocket      # WebSocket implementation
├───repository         # Data access layer
├───templates          # HTML templates
└───usecase            # Application use cases
```

## Installation

To get started with the Chatbox project, follow these steps:

### Clone the repository:

- sh
- Copy code
- git clone https://github.com/yourusername/chatbox.git

### Navigate to the project directory:

- sh
- Copy code
- cd chatbox

### Install the dependencies:

sh
Copy code
go mod download

### Set up your environment variables as needed. You can use the provided .env.example file as a template:

- sh
- Copy code
- cp .env.example .env

## Run the application:

- sh
- Copy code
- go run main.go

## Usage
Once the application is running, you can access the chatbox through your web browser at http://localhost:8080.

## Features
- User Management: Register, login, and manage user accounts.
- Chat Rooms: Create and join chat rooms.
- Messaging: Real-time messaging with WebSocket.
- Logging: Detailed logging of activities.
- Email Services: Integration with email services for notifications.
- OAuth2: Support for OAuth2 authentication.
- Cloudinary: Integration with Cloudinary for media uploads.

