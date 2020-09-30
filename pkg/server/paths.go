package server

// APIPath uri endpoint
const APIPath = "/api"

// TypesPath uri endpoint
const TypesPath = APIPath + "/types"

// StatesPath uri endpoint
const StatesPath = APIPath + "/states"

// GenresPath uri endpoint
const GenresPath = APIPath + "/genres"

// AnimesPath uri endpoint
const AnimesPath = APIPath + "/animes"

// LatestEpisodesPath uri endpoint
const LatestEpisodesPath = APIPath + "/latest"

// DirectoryPath uri endpoint
const DirectoryPath = APIPath + "/directory"

// AnimeDetailsPath uri endpoint
const AnimeDetailsPath = APIPath + "/animes/{flvid:[0-9]+}"

// EpisodeListPath uri endpoint
const EpisodeListPath = AnimeDetailsPath + "episodes"

// EpisodeDetailsPath uri endpoint
const EpisodeDetailsPath = EpisodeListPath + `{eNumber}`

// AllPaths array with all paths (uris)
var AllPaths = []string{APIPath, TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}

// AllPathsWithoutIndex array with all paths (uris) except the index
var AllPathsWithoutIndex = []string{TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}
