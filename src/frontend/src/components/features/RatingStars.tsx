import './RatingStars.css';

interface RatingStarsProps {
    rating: number;
    max?: number;
    size?: 'sm' | 'md' | 'lg';
    interactive?: boolean;
    onChange?: (rating: number) => void;
}

export function RatingStars({
    rating,
    max = 5,
    size = 'md',
    interactive = false,
    onChange
}: RatingStarsProps) {
    const stars = [];

    // Round to nearest 0.5
    const roundedRating = Math.round(rating * 2) / 2;

    const handleStarClick = (index: number) => {
        if (interactive && onChange) {
            onChange(index + 1);
        }
    };

    for (let i = 0; i < max; i++) {
        let starClass = 'star';

        if (i < Math.floor(roundedRating)) {
            starClass += ' filled';
        } else if (i === Math.floor(roundedRating) && roundedRating % 1 !== 0) {
            starClass += ' half';
        }

        stars.push(
            <span
                key={i}
                className={starClass}
                onClick={() => handleStarClick(i)}
            >
                ★
            </span>
        );
    }

    return (
        <div className={`rating-stars ${size} ${interactive ? 'interactive' : ''}`}>
            {stars}
        </div>
    );
}
