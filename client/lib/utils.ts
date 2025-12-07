import { Movie } from "@/api";

export const MovieMapper = (movie:Movie) => {
    return {
        id: movie.id,
        title: movie.title,
        overview: movie.overview,
        posterUrl: movie.poster_path
            ? `https://image.tmdb.org/t/p/w500${movie.poster_path}`
            : null,
        backdropUrl: movie.backdrop_path
            ? `https://image.tmdb.org/t/p/w780${movie.backdrop_path}`
            : null,
        releaseDate: movie.release_date,
    };
};