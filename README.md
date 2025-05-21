# JaxStats - League of Legends Stats Analysis (Kubernetes Edition)

JaxStats is a Python-based League of Legends stats analysis tool that uses Jax for machine learning-powered performance analysis. This repository is focused on deploying JaxStats as a scalable web service using Kubernetes (KinD or other clusters).

## Features

- Riot Games API integration for fetching match data
- Detailed player statistics analysis
- Machine learning-powered performance rating using Jax
- Interactive web interface (FastAPI + Jinja2)
- Ready for local Kubernetes deployment

## Quick Start (Kubernetes)

For a step-by-step guide to running JaxStats locally with Kubernetes and KinD, see [DEPLOYMENT.md](./DEPLOYMENT.md).

## Project Structure

- `app/` - Main application directory
  - `api/` - Riot Games API client
  - `analysis/` - Data processing and analysis
  - `ml/` - Jax-based machine learning models
  - `static/` - Frontend assets
  - `templates/` - HTML templates
  - `main.py` - FastAPI application entry point
- `k8s-deployment.yaml` - Kubernetes deployment and service manifest
- `k8s-secret.yaml` - Kubernetes secret manifest for Riot API key
- `Dockerfile` - Containerization for Kubernetes

## Local Development (Optional)

You can also run JaxStats locally without Kubernetes:

1. Create a virtual environment:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```
2. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```
3. Set up your Riot Games API key in a `.env` file:
   ```
   RIOT_API_KEY=your_api_key_here
   ```
4. Run the application:
   ```bash
   uvicorn app.main:app --reload
   ```
5. Open your browser and navigate to `http://localhost:8000`

## Note

This application requires a valid Riot Games API key. You can obtain one from the [Riot Games Developer Portal](https://developer.riotgames.com/).

## License

MIT 