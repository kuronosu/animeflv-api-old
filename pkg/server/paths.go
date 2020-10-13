package server

// APIPath uri endpoint
const APIPath = "/api"

// TypesPath uri endpoint
const TypesPath = APIPath + "/types"

// TypeDetailsPath uri endpoint
const TypeDetailsPath = TypesPath + "/{id:[0-9]+}"

// StatesPath uri endpoint
const StatesPath = APIPath + "/states"

// StateDetailsPath uri endpoint
const StateDetailsPath = StatesPath + "/{id:[0-9]+}"

// GenresPath uri endpoint
const GenresPath = APIPath + "/genres"

// GenreDetailsPath uri endpoint
const GenreDetailsPath = GenresPath + "/{id:[0-9]+}"

// AnimesPath uri endpoint
const AnimesPath = APIPath + "/animes"

// LatestEpisodesPath uri endpoint
const LatestEpisodesPath = APIPath + "/latest"

// AnimeDetailsPath uri endpoint
const AnimeDetailsPath = AnimesPath + "/{flvid:[0-9]+}"

// SearchAnimePath uri endpoint
const SearchAnimePath = AnimesPath + "/search"

// EpisodeListPath uri endpoint
const EpisodeListPath = AnimeDetailsPath + "/episodes"

// EpisodeDetailsPath uri endpoint
const EpisodeDetailsPath = EpisodeListPath + `/{eNumber}`

// VideoPath uri endpoint
const VideoPath = EpisodeDetailsPath + `/{server}`

// VideoLangPath uri endpoint
const VideoLangPath = VideoPath + `/{lang}`

// DirectoryPath uri endpoint
const DirectoryPath = APIPath + "/directory"

// AllPaths array with all paths (uris)
var AllPaths = []string{APIPath, TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}

// AllPathsWithoutIndex array with all paths (uris) except the index
var AllPathsWithoutIndex = []string{TypesPath, StatesPath, GenresPath, AnimesPath, LatestEpisodesPath, DirectoryPath}
