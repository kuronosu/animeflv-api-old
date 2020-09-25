package server

// IndexPath uri endpoint
const IndexPath = "/"

// TypesPath uri endpoint
const TypesPath = "/api/types/"

// StatesPath uri endpoint
const StatesPath = "/api/states/"

// GenresPath uri endpoint
const GenresPath = "/api/genres/"

// AnimesPath uri endpoint
const AnimesPath = "/api/animes/"

// LatestEpisodesPath uri endpoint
const LatestEpisodesPath = "/api/latest/"

// DirectoryPath uri endpoint
const DirectoryPath = "/api/directory/"

// AnimeDetailsPath uri endpoint
const AnimeDetailsPath = "/api/animes/{flvid:[0-9]+}/"

// EpisodeDetailsPath uri endpoint
const EpisodeDetailsPath = AnimeDetailsPath + "episodes/"

// AllPaths array with all paths (uris)
var AllPaths = []string{IndexPath, TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}

// AllPathsWithoutIndex array with all paths (uris) except the index
var AllPathsWithoutIndex = []string{TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}
