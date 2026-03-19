CREATE TABLE experience_locations (
    location_id   VARCHAR(512) NOT NULL,
    type          ENUM ('start_point', 'pickup_point', 'redemption_point', 'end_point') NOT NULL
);
