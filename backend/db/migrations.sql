-- Users
CREATE TABLE IF NOT EXISTS `users` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `uuid` TEXT UNIQUE,
    `uuid_exp` INT ,
    `nickname` TEXT NOT NULL UNIQUE,
    `email` TEXT NOT NULL UNIQUE,
    `password` TEXT NOT NULL,
    `first_name` TEXT NOT NULL,
    `last_name` TEXT NOT NULL,
    `age` INTEGER NOT NULL,
    `gender` TEXT CHECK(`gender` IN ('male', 'female')),
    `created_at` INT ,
    `last_seen` INT,
    `image` TEXT
);


-- Posts with categories
CREATE TABLE IF NOT EXISTS `posts` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `title` TEXT NOT NULL CHECK(LENGTH(`title`) BETWEEN 3 AND 100),
    `content` TEXT NOT NULL CHECK(LENGTH(`content`) BETWEEN 10 AND 2000),
    `categories` TEXT NOT NULL,
    `created_at` INT DEFAULT CURRENT_TIMESTAMP,
    `image` TEXT,
    FOREIGN KEY(`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

-- Private Messages with read status
CREATE TABLE IF NOT EXISTS `messages` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `sender_id` INTEGER NOT NULL,
    `receiver_id` INTEGER NOT NULL,
    `content` TEXT NOT NULL CHECK(LENGTH(`content`) BETWEEN 1 AND 500),
    `created_at` INT DEFAULT CURRENT_TIMESTAMP,
    `is_read` BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(`sender_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY(`receiver_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `comments` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `post_id` INTEGER NOT NULL,
    `content` TEXT NOT NULL,
    `created_at` INTEGER NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `likes` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `post_id` INTEGER,
    `comment_id` INTEGER,
    `like` INT NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`) ON DELETE CASCADE
);