# EduScroll

EduScroll is an innovative educational platform that generates short-form educational videos from textbook content. This project aims to make learning more accessible and engaging by transforming traditional textbook material into bite-sized video content.

## Project Overview

- **Backend**: Django-based service that handles user management and content generation
- **Frontend**: Swift application providing a user interface for interacting with the generated content
- **Deployment**: Docker Compose for easy orchestration of backend services

## High-Level Functionality

1. **Content Generation**: EduScroll takes textbook content as input and generates short-form educational videos.
2. **User Management**: The backend handles user accounts, preferences, and interactions.
3. **Video Delivery**: The Swift frontend app allows users to view and interact with the generated video content.

## Getting Started

### Backend

The backend services can be easily run using Docker Compose. This orchestrates all the necessary services for the Django-based backend.

To start the backend:

```bash
docker-compose up
