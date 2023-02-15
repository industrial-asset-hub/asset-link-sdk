# Contributing

We welcome contributions in several forms, e.g.

- Documenting
- Testing
- Coding
- etc.

Please read
[14 Ways to Contribute to Open Source without Being a Programming Genius or a Rock Star](https://smartbear.com/blog/14-ways-to-contribute-to-open-source-without-being/).

As far as possible we are providing templates and automation to simplify the
daily management of the below mentioned guidelines.
Therefore in some cases you won't need to do anything to comply.

## Git Guidelines

### Workflow

We currently recommend the **[Feature Branch Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow)**.

We prefer [**rebasing**](https://www.atlassian.com/git/tutorials/merging-vs-rebasing)
over merging branches.
This approach improves readability of the git history.
Repository configuration might be even enforcing it.
In order to have traceability from the commits to the originating Merge Requests
we are configuring the projects to create
[merge commits linking to the corresponding Merge Request](https://docs.gitlab.com/ee/user/project/merge_requests/methods/index.html#merge-commit-with-semi-linear-history).

The mentioned links from Atlassian are the recommended docs to read and
understand the git workflows.

In order to have traceability from requirements down to the code we are using a
[semi-linear history with merge commit](https://docs.gitlab.com/ee/user/project/merge_requests/methods/#merge-commit-with-semi-linear-history).
The merge commit has to provide a link to the corresponding Jira issue
motivating whatever is being merged.

No commits are allowed to the main (release) branch not coming from a merge
request.

### Git Commit

The cardinal rule for creating good commits is to ensure there is **only one
"logical change" per commit**. There are many reasons why this is an important
rule:

- The smaller the amount of code being changed, the quicker & easier it is to
  review & identify potential flaws.
- If a change is found to be flawed later, it may be necessary to revert the
  broken commit. This is much easier to do if there are not other unrelated code
  changes entangled with the original commit.
- When troubleshooting problems using Git's bisect capability, small well
  defined changes will aid in isolating exactly where the code problem was
  introduced.
- When browsing history using Git annotate/blame, small well defined changes
  also aid in isolating exactly where & why a piece of code came from.

Things to avoid when creating commits

- Mixing whitespace changes with functional code changes.
- Mixing two unrelated functional changes.
- Sending large new features in a single giant commit.

### Git Commit Conventions

More information about the reason for some of these rules is available
[here](https://cbea.ms/git-commit/).

#### Commit Message Subject

Subject line should be separated with a blank line from the body.

Try to limit the subject line to 50 characters.

Don't end the subject with a period character.

Use the imperative mood. For example: "clean ..." instead of "we have cleaned
...".

We use git commit as per [Conventional Changelog](https://github.com/ajoslin/conventional-changelog):

``` text
<type>(<scope>): <subject>
```

Example:

``` text
docs(CONTRIBUTING.md): add commit message guidelines
```

Allowed types:

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space,
  formatting, missing semi-colons, etc)
- **refactor**: A code change that neither fixes a bug or adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests
- **chore**: Changes to the build process or auxiliary tools and libraries such
  as documentation generation

#### Commit Message Body

Try to wrap the lines at 72 characters.

Use the body to explain what and why vs. how.

#### Commit Message Footer

Every commit author must 'sign off' the created commits. Usage of tool support
is recommended for this task. See [here](https://stackoverflow.com/a/1962112)
for an explanation about this feature.

Every reviewer must also either 'sign off' the merge commit or add a
`Reviewed-by: John Doe <john.doe@siemens.com>` at the footer of the commit
message.

#### Commit Messages in Branches

Rewritting the branch commits once ready for merging is allowed.

Commit in branches don't necessarily need to be compliant with the above rules.
Since rewritting the commits is allowed, they can be fixed before merging.
The usage of commit linters to avoid undesired commit message formats in the
default branch is highly recommendable.

## File Headers

All files providing support for comments should provide a header with following
information:

- Copyright Notice
- Authorship Notice

### Copyright Notice

Siemens has the copyright on the created files.
In order to explicitly declare it, following notice must be provided in the
header!

``` text
(c) Siemens <year-of-creation>
```

or

``` text
(c) Siemens <year-of-creation>-<year-of-last-modification>
```

For example for a Python file created 2022 and modified lastly 2023, this must
be the notice:

``` Python
# (c) Siemens 2022-2023
```

### Authorship Notice

Although the Git history helps keeping track of who created/modified a file, it
doesn't keep track of who created the content of the file.
For example, if someone copies a file created by someone else, Git will report
the committer (who copied the file) as the author, but the file itself would
declare the right author.

Authorship notices aren't mandatory, but recommended.

An Authorship notice in Python would look like this:

``` Python
# Author: John Doe <john.doe@siemens.com>
```

