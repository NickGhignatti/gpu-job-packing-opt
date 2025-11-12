# GPU Job Packing Optimizer

### ❓ WHAT

A multi-tenant GPU scheduler that intelligently packs multiple small machine learning jobs onto single GPUs instead of dedicating one GPU per job. It's a container-based HPC job scheduler with ML-powered resource prediction and intelligent GPU allocation that maximizes hardware utilization.

### ❓ WHY
GPUs are expensive and often underutilized—many ML jobs don't fully saturate a GPU's compute or memory capacity. By co-locating compatible jobs on the same GPU, you can:
- Dramatically improve GPU utilization rates
- Reduce infrastructure costs
- Decrease job queue wait times
- Enable more researchers/teams to share limited GPU resources efficiently

### ❓ HOW
The system uses a five-component architecture:

1. Job Submission API (FastAPI): RESTful interface for submitting ML jobs
2. ML Predictor: Analyzes incoming jobs to predict GPU memory usage, runtime duration, and interference levels between jobs
3. Scheduler Core: Implements bin packing algorithms and interference-aware placement logic, leveraging NVIDIA MPS/MIG for isolation
4. Job Executor: Manages containerized execution environments with resource allocation enforcement
5. GPU Monitor: Tracks real-time NVML metrics and provides feedback for dynamic optimization

The predictor estimates resource requirements, the scheduler decides which jobs can safely coexist on the same GPU, and the executor runs them in isolated containers.