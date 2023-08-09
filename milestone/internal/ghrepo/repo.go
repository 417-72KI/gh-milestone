package ghrepo

type Interface interface {
	RepoName() string
	RepoOwner() string
	RepoHost() string
}
