export const profile = {
  name: 'Raphael San Andres',
  title: 'ML & AI Systems Engineer',
  tagline: 'ML and AI Systems Engineer with 3+ years building and deploying ML systems at scale. First hire on a customer-facing ML support team managing multi-node NVIDIA H100 GPU clusters at 99.9% uptime.',
  avatar: '/IMG_3333.jpg',
  location: 'Pleasanton, CA',
}

export const about = `Hello! I'm Raphael. I'm a person who gets excited about linguistics, loves diving into video game statistics, and has a bunch of other random interests floating around. Professionally, I'm working in the cool world of ML/AI.`

export interface ExperiencePosition {
  title: string
  period: string
  details: string[]
}

export interface ExperienceItem {
  company: string
  subtitle?: string
  title?: string
  period?: string
  details?: string[]
  positions?: ExperiencePosition[]
}

export const experience: ExperienceItem[] = [
  {
    company: 'Stealth Startup',
    subtitle: 'GPU Cloud Infrastructure Provider',
    title: 'ML Solutions Engineer (First Hire) | Software Engineer',
    period: 'Feb 2024 - Present',
    details: [
      'First engineer on a 5-person customer-facing ML support team supporting several early-stage AI startups on GPU infrastructure',
      'Managed multi-node NVIDIA H100 GPU Kubernetes clusters, maintaining 99.9% uptime across enterprise AI infrastructure',
      'Helped design and implement automated node repair processes, reducing mean time to repair from 36-72 hours to under 1 hour',
      'Built internal knowledge base and ticketing codebase RAG system to accelerate issue resolution across the support team',
      'Leading infrastructure readiness for next-generation NVIDIA GB100 (Blackwell) GPU deployments',
      'Developed bidirectional Jira-Kubernetes operators in Go, automating incident response for ~20 weekly node failures',
      'Optimized SLURM and Kubernetes ML workflows, improving training job throughput for customers running distributed workloads',
    ],
  },
  {
    company: 'Weights & Biases',
    title: 'Machine Learning Support Engineer',
    period: 'Jan 2023 - Jan 2024',
    details: [
      'Debugged and resolved 600+ technical issues for ML practitioners at OpenAI, NVIDIA, and Microsoft, covering model integrations, LLM deployments, and on-premise instances',
      'Triaged and traced 50+ bugs across the Weights & Biases SDK, web application, backend services, and managed instances, contributing bugfix PRs and new integrations',
      'Managed ~20 customer requests daily while running debugging sessions, cross-team syncs, and building internal tooling (W&B integrations, frontend features)',
    ],
  },
]

export const skillCategories = [
  {
    name: 'Languages',
    skills: ['Python', 'Go', 'TypeScript', 'SQL', 'R', 'C++'],
  },
  {
    name: 'ML & AI',
    skills: ['PyTorch', 'TensorFlow', 'Keras', 'XGBoost', 'HuggingFace', 'LangChain', 'LangGraph', 'Ray', 'OpenAI API', 'Ollama'],
  },
  {
    name: 'Infrastructure',
    skills: ['Kubernetes', 'Docker', 'SLURM', 'AWS SageMaker', 'GCP Vertex', 'Azure', 'Helm', 'Kustomize', 'GitHub Actions'],
  },
  {
    name: 'Data & Observability',
    skills: ['PostgreSQL', 'pgvector', 'DataDog', 'Prometheus', 'Grafana', 'Jupyter', 'Weights & Biases'],
  },
  {
    name: 'Certifications',
    skills: ['Google Data Analytics Professional Certificate'],
  },
]

export const education = [
  {
    degree: 'Masters in Artificial Intelligence',
    school: 'Penn State',
    year: 'May 2025',
  },
  {
    degree: 'Bachelor of Science in Statistics',
    school: 'UCLA',
    year: 'June 2022',
  },
]

export interface ProjectItem {
  name: string
  description: string
  period: string
  tech?: string
  url?: string
  github?: string
}

export const projects: ProjectItem[] = [
  {
    name: 'Atlas AI',
    description: 'Clinical intelligence platform — query patient records in natural language via a 3-node multi-agent RAG pipeline with 20+ medical tools.',
    period: 'Active',
    tech: 'LangGraph, FastAPI, Next.js, pgvector, AWS Bedrock, ECS, RDS, Ollama',
    url: 'https://hcai.rsanandres.com',
    github: 'https://github.com/rsanandres/hc_ai',
  },
  {
    name: 'Atlas AI MCP',
    description: 'MCP server exposing healthcare AI tools for RAG-powered clinical queries, document reranking, and FHIR data ingestion. Published on PyPI and Smithery.',
    period: 'Active',
    tech: 'Python, MCP',
    url: 'https://smithery.ai/servers/rsanandres/hc_ai_mcp',
    github: 'https://github.com/rsanandres/hc_ai_mcp',
  },
  {
    name: 'Corium',
    description: 'Kubernetes operator with 3 custom CRDs for automated metrics collection, threshold alerting, and a monitoring dashboard.',
    period: 'Recent',
    tech: 'Go, Kubebuilder, Prometheus, Grafana, Next.js',
    github: 'https://github.com/rsanandres/corium',
  },
  {
    name: 'JaxStats',
    description: 'Game performance analyzer with XGBoost ML scoring, 8-metric GPI breakdown, AI coaching, and live game overlays.',
    period: 'Active',
    tech: 'FastAPI, XGBoost, React, Chart.js, Ollama',
    github: 'https://github.com/rsanandres/jaxstats',
  },
  {
    name: 'Aphae',
    description: 'AI agent office simulation — procedurally generated personalities (Big Five), emergent relationships, LLM-driven conversations, and a drama director.',
    period: 'Active',
    tech: 'Godot 4, GDScript, Ollama',
    github: 'https://github.com/rsanandres/aphae',
  },
  {
    name: 'Models from Scratch',
    description: 'MLP, CNN, and Transformer self-attention implemented from scratch with hand-derived forward and backward passes — no frameworks, no autograd.',
    period: 'Recent',
    tech: 'Python, NumPy',
    github: 'https://github.com/rsanandres/models-from-scratch',
  },
  {
    name: 'Capstone (PSU AI 894)',
    description: 'Predicts NFL formation based on player positions and coordinates.',
    period: 'Spring 2024',
    github: 'https://github.com/JohnnyZ67/AA894-Group-4-Capstone',
  },
]

export const publications = [
  {
    title: 'RL with Mario',
    url: 'https://wandb.ai/raphael-sanandres/cleanRL/reports/Can-We-Beat-Mario-Bros-with-Gymnasium-and-CleanRL-on-a-Laptop---Vmlldzo0NTcxNTcw',
    type: 'Article',
  },
]

export const socialLinks = [
  {
    href: 'https://github.com/rsanandres',
    label: 'GitHub',
  },
  {
    href: 'https://www.linkedin.com/in/raphael-san-andres/',
    label: 'LinkedIn',
  },
]
