package server

// IndexPath uri
const IndexPath = "/"

// TypesPath uri
const TypesPath = "/types/"

// StatesPath uri
const StatesPath = "/states/"

// GenresPath uri
const GenresPath = "/genres/"

// AnimesPath uri
const AnimesPath = "/animes/"

// LatestEpisodesPath uri
const LatestEpisodesPath = "/latest/"

// DirectoryPath uri
const DirectoryPath = "/directory/"

// AllPaths array with all paths (uris)
var AllPaths = []string{IndexPath, TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}

// AllPathsWithoutIndex array with all paths (uris) except the index
var AllPathsWithoutIndex = []string{TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}
