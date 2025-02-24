DELETE FROM users;

UPDATE sqlite_sequence SET seq = 0 WHERE name = 'users';

DROP TABLE users

INSERT INTO
    `users` (
        `uuid`,
        `uuid_exp`,
        `nickname`,
        `email`,
        `password`,
        `first_name`,
        `last_name`,
        `age`,
        `gender`,
        `created_at`,
        `last_seen`
    )
VALUES (
        'uuid1',
        1672531199,
        'john_doe',
        'john.doe@example.com',
        'password123',
        'John',
        'Doe',
        25,
        'male',
        1672444800,
        1672444800
    ),
    (
        'uuid2',
        1672531199,
        'jane_smith',
        'jane.smith@example.com',
        'password123',
        'Jane',
        'Smith',
        30,
        'female',
        1672444800,
        1672444800
    ),
    (
        'uuid3',
        1672531199,
        'alice_johnson',
        'alice.johnson@example.com',
        'password123',
        'Alice',
        'Johnson',
        22,
        'female',
        1672444800,
        1672444800
    ),
    (
        'uuid4',
        1672531199,
        'bob_brown',
        'bob.brown@example.com',
        'password123',
        'Bob',
        'Brown',
        28,
        'male',
        1672444800,
        1672444800
    ),
    (
        'uuid5',
        1672531199,
        'charlie_davis',
        'charlie.davis@example.com',
        'password123',
        'Charlie',
        'Davis',
        35,
        'male',
        1672444800,
        1672444800
    ),
    (
        'uuid6',
        1672531199,
        'diana_evans',
        'diana.evans@example.com',
        'password123',
        'Diana',
        'Evans',
        27,
        'female',
        1672444800,
        1672444800
    ),
    (
        'uuid7',
        1672531199,
        'edward_garcia',
        'edward.garcia@example.com',
        'password123',
        'Edward',
        'Garcia',
        40,
        'male',
        1672444800,
        1672444800
    ),
    (
        'uuid8',
        1672531199,
        'fiona_hall',
        'fiona.hall@example.com',
        'password123',
        'Fiona',
        'Hall',
        33,
        'female',
        1672444800,
        1672444800
    ),
    (
        'uuid9',
        1672531199,
        'george_lee',
        'george.lee@example.com',
        'password123',
        'George',
        'Lee',
        29,
        'male',
        1672444800,
        1672444800
    ),
    (
        'uuid10',
        1672531199,
        'hannah_martin',
        'hannah.martin@example.com',
        'password123',
        'Hannah',
        'Martin',
        26,
        'female',
        1672444800,
        1672444800
    );

INSERT INTO
    `posts` (
        `user_id`,
        `title`,
        `content`,
        `categories`,
        `created_at`,
        `image`
    )
VALUES (
        1,
        'First Post',
        'This is the content of the first post.',
        'Technology',
        1672444800,
        'image1.jpg'
    ),
    (
        2,
        'Second Post',
        'This is the content of the second post.',
        'Science',
        1672444800,
        'image2.jpg'
    ),
    (
        3,
        'Third Post',
        'This is the content of the third post.',
        'Health',
        1672444800,
        'image3.jpg'
    ),
    (
        4,
        'Fourth Post',
        'This is the content of the fourth post.',
        'Education',
        1672444800,
        'image4.jpg'
    ),
    (
        5,
        'Fifth Post',
        'This is the content of the fifth post.',
        'Travel',
        1672444800,
        'image5.jpg'
    ),
    (
        6,
        'Sixth Post',
        'This is the content of the sixth post.',
        'Food',
        1672444800,
        'image6.jpg'
    ),
    (
        7,
        'Seventh Post',
        'This is the content of the seventh post.',
        'Fashion',
        1672444800,
        'image7.jpg'
    ),
    (
        8,
        'Eighth Post',
        'This is the content of the eighth post.',
        'Sports',
        1672444800,
        'image8.jpg'
    ),
    (
        9,
        'Ninth Post',
        'This is the content of the ninth post.',
        'Music',
        1672444800,
        'image9.jpg'
    ),
    (
        10,
        'Tenth Post',
        'This is the content of the tenth post.',
        'Art',
        1672444800,
        'image10.jpg'
    );

INSERT INTO
    `comments` (
        `user_id`,
        `post_id`,
        `content`,
        `created_at`
    )
VALUES (
        1,
        1,
        'Great post!',
        1672444800
    ),
    (
        2,
        2,
        'Interesting read.',
        1672444800
    ),
    (
        3,
        3,
        'Very informative.',
        1672444800
    ),
    (
        4,
        4,
        'Thanks for sharing.',
        1672444800
    ),
    (
        5,
        5,
        'Awesome content!',
        1672444800
    ),
    (6, 6, 'Loved it!', 1672444800),
    (
        7,
        7,
        'Well written.',
        1672444800
    ),
    (
        8,
        8,
        'Keep it up!',
        1672444800
    ),
    (
        9,
        9,
        'Fantastic!',
        1672444800
    ),
    (
        10,
        10,
        'Amazing post!',
        1672444800
    );

INSERT INTO
    `likes` (
        `user_id`,
        `post_id`,
        `comment_id`,
        `like`
    )
VALUES (2, 14, NULL, 1),
    (1, 1, NULL, 1),
    (3, 3, NULL, 1),
    (4, 4, NULL, 1),
    (5, 5, NULL, 1),
    (6, 6, NULL, 1),
    (7, 7, NULL, 1),
    (8, 8, NULL, 1),
    (9, 9, NULL, 1),
    (10, 10, NULL, 1);

UPDATE users
set
    password = "$2a$10$bVPk0LM0bUE6UiL5MM91d.LOidE.EsKdzi.L0i5nJO5v5EpzO6Rwa"

    UPDATE users
set
    image = ""