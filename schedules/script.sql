INSERT INTO developers (name, department, geolocation, last_known_ip, is_available)
SELECT
    CONCAT(
            (ARRAY['James','Mary','John','Patricia','Robert'])[floor(random() * 5) + 1],
            ' ',
            (ARRAY['Smith','Johnson','Williams','Brown','Jones'])[floor(random() * 5) + 1]
    ) AS name,

    (ARRAY['backend','frontend','ios','android'])[floor(random() * 4) + 1] AS department,

    ST_SetSRID(
            ST_MakePoint(
                    10 + random() * 30,
                    45 + random() * 15
            ),
            4326
    )::GEOGRAPHY AS geolocation,

    CONCAT(
            (random()*255)::INT, '.',
            (random()*255)::INT, '.',
            (random()*255)::INT, '.',
            (random()*255)::INT
    )::INET AS last_known_ip,

    random() > 0.5 AS is_available
FROM generate_series(1, 1000);