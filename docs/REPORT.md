# Car Store System
## Assignment 3 – Project Design & Setup Milestone

---

## 1. Project Proposal

The project proposal describes the relevance of the system, competitor analysis, target users, planned features, and scope boundaries.

**Proposal document:**
- [Project Proposal](Proposal.md)

---

## 2. Architecture & Design — Monolithic Architecture

The system is designed using a **monolithic architecture**, where all components are part of a single deployable application. The architecture follows a layered approach to separate responsibilities and improve maintainability.

### 2.1 System Architecture Overview
The application is divided into the following layers:
- **Handlers (Controllers):** Handle HTTP requests and responses
- **Services:** Contain business logic
- **Repositories:** Handle database operations
- **Database:** Stores persistent data

### 2.2 Diagrams

The following diagrams are provided to describe the system design:

- **Use-Case Diagram** — shows user interactions with the system  
[Use-Case Diagram](UseCase.png)

- **ERD (Entity Relationship Diagram)** — describes database schema and relationships  
[ERD Diagram](ERD.pdf)

- **UML Diagram** — represents system classes and their relationships  
[UML Diagram](uml.png) (to be added by Ehson)

---

## 3. Project Plan — Gantt Chart (Weeks 7–10)

The project plan outlines the development tasks distributed across team members during weeks 7–10 of the semester.

### 3.1 Task Distribution

- **Nurdaulet Usenbai**
    - Project proposal
    - Cars/Catalog module implementation

- **Nurbol**
    - ERD design
    - Use-case diagram
    - Orders module implementation

- **Ehson Usmanov**
    - System architecture design
    - UML diagram
    - Authentication and authorization module

### 3.2 Gantt Chart

**Gantt Plan:**
- Gantt diagram for weeks 7–10
  `docs/Gantt.pdf`
- [Gantt](Gantt.pdf)

---

## 4. Repository Setup

### 4.1 Git Repository
- Repository URL:  
  https://github.com/usenbai-nur/AdvancedProgramming

### 4.2 Branching Strategy
Each team member works in a separate branch:
- `nurdaulet-proposal`
- `nurbol-erd-usecase`
- `ehson-architecture`

Each branch contains individual commits demonstrating contribution to the project.