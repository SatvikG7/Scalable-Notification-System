CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE,
    status VARCHAR(50) NOT NULL,
    preference_low_channel VARCHAR(50) NOT NULL,
    preference_medium_channel VARCHAR(50) NOT NULL,
    preference_high_channel VARCHAR(50) NOT NULL
);