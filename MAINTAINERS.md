# Maintainers Guide

This guide is intended for anybody with maintainer rights in the project.
The term *project* refers within this document to the code.siemens.com project
where this document can be found.

Some of the configurations are automatically applied if you create the project
by importing the "project template".

Those configurations automatically provided by the template start with the text
"*Automatically applied by the* **project template**" in the corresponding
section.

Some tooling might exist to automatically check conformance of a project
with this guide.
Automatic enforcement is not desired as of this writting.

Any documentation about the motivation for these configurations is provided in
the [Contributors Guidelines](CONTRIBUTING.md) since they have impact on the
contributors and they should be aware of the rationals behind the configuration.

## Merge Configuration

*Automatically applied by the* **project template**

*Manual configuration:*

Go to the merge requests configuration (`Settings` ->
`General` -> `Merge requests`) and ensure that `Merge commit with semi-linear
history` is the selectecd `Merge method` in case something else but the default
is configured, the [Contributions Guidelines](CONTRIBUTING.md) should be
accordingly updated.

## Default Branch

*Automatically applied by the* **project template**

*Manual configuration:*

The default branch (`Settings` -> `Repository` ->
`Default branch`) is called `main`.

## Branches and Tags Protection

*Automatically applied by the* **project template**

*Manual configuration:*

The default branch should be protected (`Settings` -> `Repository` ->
`Protected branches`) and configured so:

* `Allowed to merge` -> `Maintainers`
* `Allowed to push` -> `No one`
* `Allowed to force push` -> *disabled*

This way no commit get into the main branch without a merge request, what
implies a review process (see [below](#review-merge-requests-to-main).
Additionally at least one project maintainer (make sure to provide at least two)
should be involved in the merging process, making sure that the merge request
complies with the guidelines.

## Review Merge Requests to Main

We work with a 4-eyes principle for all the code getting into the main branch.
Meaning that every merge request must be reviewed by at least one person that is
not the merge request author.
This rule also **applies to maintainers**!

A merge request should be assigned to a project maintainer, also a reviewer
other than the merge request author should be selected.
This is optional if the maintainer that is the assignee of the merge request
should also review the merge request.

## Merging Branches

*Automatically applied by the* **project template**

*Manual configuration:*

Go to the merge requests configuration (`Settings` ->
`General` -> `Merge requests`) and ensure that following options are active in
the `Merge checks` section:

* `Pipelines must succeed` (must be disabled if the project doesn't have any
  pipelines)
* `All threads must be resolved`

## Merge Commit Message

*Automatically applied by the* **project template**

*Manual configuration:*

The message of the commit being added by a merge request should contain
following information:

* A reference to the JIRA issue that requires the changes (no merge request
    without its corresponding JIRA issue).
* A `Merged-by: <maintainer name and e-mail>` footer.
* One or more `Approved-by: <reviewer name and e-mail>` footers (optional if
  maintainer merging is also reviewer).

## README.md

The maintainers are responsible for ensuring that a README.md file providing
following information:

* It has to contain an `Introduction` section which briefly explains what the
  repository provides.
* The `Introduction` has to be written in a language that can be easily
  understood by all potential stakeholders.
* It has to be kept up-to-date. Move to other files sections that are hard to
  be kept up-to-date.
* It should refer to any other files relevant to any potential stakeholders
  (for example to `CONTRIBUTING.md` for developers and to `MAINTAINERS.md` for
  maintainers and quality assurance).

### License

Should it be needed to license the reuse of the content of this repository, add
to the repository the file `LICENSE` and replace the content of the
[License section of the file README.md](README.md#license) with following text:

``` markdown
Reuse of the content of this repository has been licensed under the conditions
described in the file [LICENSE](LICENSE).
```

