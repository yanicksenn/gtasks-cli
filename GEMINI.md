# Agent Coding Protocol

## 1. Core Principles

As a coding agent, you must adhere to the following principles to ensure the delivery of high-quality, well-documented, and robust code. Your primary directive is to follow a structured and transparent development process.

- **Clarity:** Be explicit about your intentions and actions.
- **Simplicity:** Make small, incremental, and easily understandable changes.
- **Quality:** Always verify your work with tests.
- **Consistency:** Follow the established workflow for every change.

## 2. Development Workflow

You must follow this workflow for every new feature, bug fix, or modification.

### Step 1: Planning & Design

For any new functionality, you **MUST** begin by creating a **Software Design Document (SDD)**.

- **Format:** Markdown (`.md`).
- **Location:** Create a new file, for example, in a `docs/` or `design/` directory. The name should be descriptive (e.g., `feature-user-authentication-sdd.md`).
- **Content:** The SDD **MUST** include the following sections:
    1.  **Date:** The date of the document's creation.
    2.  **Goal:** A clear and concise description of the feature's purpose, user value, and intended outcome.
    3.  **Technical Considerations:** A detailed analysis of:
        -   Architectural changes.
        -   Data models.
        -   Potential risks and challenges.
        -   Dependencies (new and existing).
    4.  **Testing Plan:** A comprehensive strategy outlining:
        -   The scope of testing (unit, integration, end-to-end).
        -   Specific cases to be tested.
        -   New test files to be created.

### Step 2: Task & Milestone Definition

After the SDD is finalized, you **MUST** create a **Task List** file.

- **Format:** Markdown (`.md`) with checklists.
- **Naming:** The file name **MUST** correspond to the SDD it is based on (e.g., if the SDD is `feature-user-authentication-sdd.md`, the task list should be `feature-user-authentication-tasks.md`).
- **Content:** Break down the implementation into a series of small, concrete tasks and milestones. Each task can include sub-steps or further details to ensure clarity. You will use this file to track your progress.

**Example (`feature-user-authentication-tasks.md`):**
```markdown
# Feature: User Authentication

- [ ] **Task 1: Create User model and database migration.**
    - Define `User` schema with `username`, `email`, and `password_hash`.
    - Generate and apply the database migration script.
- [ ] **Task 2: Implement registration logic.**
    - Create `POST /api/register` endpoint.
    - Implement password hashing.
    - Add validation for input fields.
- [ ] **Task 3: Write unit tests for registration.**
    - Test successful user creation.
    - Test duplicate username/email prevention.
- [ ] **Task 4: Implement login endpoint.**
- [ ] **Task 5: Write integration tests for the login flow.**
- [ ] **Task 6: Update API documentation.**
```

### Step 3: Implementation & Verification Cycle

You **MUST** follow this iterative cycle for each task in your task list:

1.  **Outline:** Before writing any code, state clearly which task you are starting and provide a brief, high-level summary of your planned changes.
2.  **Implement:** Make small, atomic changes that directly address the current task.
3.  **Verify:** After the change is made, run all relevant tests. For new features, you **MUST** add new tests as defined in your Testing Plan. Ensure all tests pass before proceeding.
4.  **Commit:** Once the task is implemented and verified, create a Git commit. The commit message should be clear and descriptive, referencing the completed task.
5.  **Update Task List:** Mark the task as complete in your task list file.

### Step 4: Documentation Update

After a feature is fully implemented and all tasks are complete, you **MUST** review the project's existing documentation (e.g., READMEs, API docs). If your changes have altered any functionality, behavior, or configuration, you **MUST** update the documentation to reflect these changes accurately.
