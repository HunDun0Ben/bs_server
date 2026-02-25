# Project Requirements

1. Use **Go**, **GoGin**, **GoCV**, and **MongoDB** as the main frameworks.

# Goals

1. Build an image recognition system where the frontend interacts with APIs for:
    - Training models
    - Performing user recognition
2. Support **configurable pipelines**, allowing users to combine different:
    - Feature extraction algorithms
    - Clustering methods
    - Training strategies  
      The system should generate corresponding trainers based on user configurations.

# Git Repository

- The primary branch is named **`master`**.

# Git Commit

- Automatically follow the standard format for Git `Conventional Commits`

# GOCV

As a Go expert using the gocv library, you must strictly follow these principles for gocv.Mat memory management. gocv.Mat requires manual memory handling and is not managed by the Go GC.

1. The Principle of Ownership The function that creates a gocv.Mat (via gocv.IMRead, NewMat, Clone, IMDecode, etc.) is solely responsible for its destruction. Always pair a creation with a defer mat.Close() call to guarantee cleanup.

2. The Principle of Borrowing (Intra-Service)

- Pass by Pointer: When passing a Mat to another function within the same process, always use a pointer (\*gocv.Mat). This signals that the callee is only "borrowing" the Mat.
- Borrowers Don't Close: A function that receives a \*gocv.Mat must not close it.
- Returning Transfers Ownership: If a function returns a new gocv.Mat, it transfers ownership to the caller. The caller is now responsible for closing it.

3. The Principle of Serialization (Inter-Service)

- Never Pass `Mat` Objects Directly: Do not attempt to serialize or send a Mat object across process boundaries (e.g., via Kafka or an API). Its internal pointer is meaningless elsewhere.
- Encode and Decode: The sending service must encode the Mat into a byte slice (e.g., JPEG/PNG). The receiving service must decode those bytes into a fresh Mat instance, thereby becoming its new owner and responsible for closing it.

# Process

1. **Analyze code for optimization opportunities**:
    - Identify Go anti-patterns that may limit compiler optimizations or runtime performance.
    - Consider Go and Gin best practices for each suggestion.
2. **Provide actionable guidance**:
    - Explain specific code changes with clear reasoning.
    - Include before/after examples where applicable.
    - Only suggest changes that meaningfully improve performance or maintainability.

# Comment Policy

- Only provide **high-value comments**.
- Avoid unnecessary guidance or conversational remarks in code comments.

# Notes on Prompt Optimization

1. **Professional tone**: Replaced casual phrases with precise technical wording.
2. **Clarity & readability**: Organized sections logically.
3. **Conciseness**: Reduced redundancy.
4. **Explicit technical terms**: Clearly defined expectations for actionable guidance.
5. **Consistency**: Standardized capitalization for technologies and branch names.
