# On Air

**Table of Contents**

- [About the project](#about-the-project)
    - [API docs](#api-docs)
    - [Clone Project](#clone-project) 
      - [Control version](#control-version)
      - [Branches](#branches) 
    - [Status](#status)
- [Getting started](#getting-started)
  - [Layout](#layout)
  - [Notes](#notes)

## About the project

The template is used to create golang project. All golang projects must follow the conventions in the
template. Calling for exceptions must be brought up in the engineering team.

### API docs

The template doesn't have API docs. For web service, please include API docs here, whether it's
auto-generated or hand-written. For auto-generated API docs, you can also give instructions on the
build process.

### Control Version

Project use git for control version and  Git Flow model for manage branches and branches will be merged to get log of all commit in all branch and pull request will be without commit message 

The Git Flow model consists of two main branches: master and develop.

in this approach all branch will merge to dev and after test will be merge to master

#### Branches

master: This branch represents the production-ready code. It should only contain code that has been thoroughly tested and is ready to be deployed to production.

develop: This branch is used to develop new features. It should contain the latest development changes and should be the base branch for all feature branches.

In addition to these two main branches, Git Flow defines three types of supporting branches

feature: These branches are used to develop new features. They are based on develop and are merged back into develop once the feature is complete.

release: These branches are used to prepare the code for a new production release. They are based on develop and are merged back into both develop and master once the release is complete.

hotfix: These branches are used to quickly fix issues in the production code. They are based on master and are merged back into both develop and master.

####  Branch naming

You can name a branch in Git using the command git branch <branch-name>, where <branch-name> is the name you want to give to the branch. For example, to create a new branch called "feature/add-login-page", you can run the following command:

git branch feature/add_login_page

This will create a new branch with the name feature/add-login-page based on the current branch you are on.

#### Commit message

ConventionalCommit message is a specific format for writing commit messages that provides a standardized way of conveying information about changes made to code in a repository.

The format consists of three parts:

    A type that describes the kind of change being made, such as feat for a new feature, fix for a bug fix, docs for documentation updates, refactor for code refactoring, and so on.

    A scope that describes the part of the codebase being modified, such as a specific module, component, or function.

    A short description that summarizes the changes made in the commit.

Optionally, the commit message can also include a longer description that provides more detailed information about the changes, as well as references to related issues, pull requests, or other relevant information.

For example, a conventional commit message for a bug fix in the authentication module of an application might look like this:

fix(auth): Validate user input before authentication

This commit fixes a bug where the authentication module could accept invalid user input, leading to security vulnerabilities. The fix adds input validation checks to the authentication process to ensure that only valid user input is processed.

Closes #123

By using a conventional commit message format, developers can more easily understand the nature and purpose of changes made to code in a repository, which can help improve collaboration, code quality, and maintenance of the codebase over time.

#### Clone Project

```bash
git clone https://github.com/amasoudfam/on-air.git
```

### Status

The template project is in alpha status.

## Getting started

Below we describe the conventions or tools specific to golang project.

### Layout

```tree
├── .github
├── .gitignore
├── .golangci.yml
├── README.md
├── build
├── docs
│   └── README.md
├── pkg
├── release
│   ├── template-admin.yaml
│   └── template-controller.yaml
├── test
│   ├── README.md
├── third_party
│   └── README.md
```

A brief description of the layout:

## Notes
