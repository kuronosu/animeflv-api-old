package server

// IndexPath uri
const IndexPath = "/"

// TypesPath uri
const TypesPath = "/api/types/"

// StatesPath uri
const StatesPath = "/api/states/"

// GenresPath uri
const GenresPath = "/api/genres/"

// AnimesPath uri
const AnimesPath = "/api/animes/"

// LatestEpisodesPath uri
const LatestEpisodesPath = "/api/latest/"

// DirectoryPath uri
const DirectoryPath = "/api/directory/"

// AnimeDetailsPath uri
const AnimeDetailsPath = "/api/animes/{flvid:[0-9]+}/"

// AllPaths array with all paths (uris)
var AllPaths = []string{IndexPath, TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}

// AllPathsWithoutIndex array with all paths (uris) except the index
var AllPathsWithoutIndex = []string{TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}
