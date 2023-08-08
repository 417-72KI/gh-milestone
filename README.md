# gh-milestone

gh-milestone is a [gh](https://github.com/cli/cli) extension to create/list/view/close milestones.

## Installation

```
gh extension install 417-72KI/gh-milestone
```

## Usage

### Authorization

You need to set your personal access token of github to `GITHUB_TOKEN` environment variable.

If you use github enterprise, you need to set your api base url to `GITHUB_BASE_URL` environment variable.

### List milestones
```
gh milestone list
```

### Create a milestone
```
gh milestone create -t <title>
```

### View a milestone
```
gh milestone view <milestone_number>
```

### Close a milestone
```
gh milestone close <milestone_number>
```
