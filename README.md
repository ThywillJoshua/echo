# Git Commit Message Generator

## Overview

This application generates Git commit messages based on the differences (diff) in staged changes. It features a web interface and uses AI to create meaningful commit messages.

## Features

- **AI-Generated Commit Messages**: Automatically generates commit messages using an AI model.
- **Web Interface**: Provides a simple web interface to interact with the application.
- **Streaming Output**: Offers real-time updates of the Git commit process.

## Project Structure

- **`main.go`**: The entry point of the application. Sets up and starts the server and web view.
- **`server` package**: Handles HTTP server setup and request handling.
- **`webview` package**: Manages the creation and display of the web view.
- **`gitops` package**: Handles Git operations and generates commit messages.
- **`generate` package**: Interfaces with AI models to generate commit messages.

## Prerequisites

- Go 1.18 or later
- Git
- OpenAI API key (if using OpenAI for generating messages)
- Ollama model (if using Ollama for message generation)

## Setup

Clone the repository:

```bash
git clone https://github.com/ThywillJoshua/echo
cd echo
```

Install dependencies:

Run the following command to download the required Go modules:

```bash
go mod tidy
```

Build the application:

To build the application, run:

```bash
go build -o echo
```

Run the application:

Start the application with:

```bash
<path-to>/echo
```

The server will start on `http://localhost:8080`, and a web view will be created.

## Usage

Access the Web Interface:

A web view will be created or open your web browser and navigate to `http://localhost:8080`. You will see the web interface where you can interact with the application.

Commit Messages:

The application will automatically check for staged changes, generate a commit message using AI, and display it in the web interface.
Then you can edit the commit messages.

Control the Application:

To close the application, navigate to `http://localhost:8080/close` or simply close the web view window.

## Troubleshooting

No Staged Changes Detected:

If you see the message "No changes staged.", ensure you have staged changes in your Git repository. Use `git add` to stage changes.

API Key Issues:

If the application fails to generate messages, check that your API key is correctly set and valid.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements. For larger changes, please discuss them with the project maintainers first.
