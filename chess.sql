-- Author: Nathaniel Martes
-- Description : Creates 3 MySQL tables representing chess data from a PGN file
USE chess;

CREATE TABLE events(
    event_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT 'Name of the event',
    site VARCHAR(100) NOT NULL COMMENT 'Site of the event',
    date DATE COMMENT 'Date of the event',
    round INT COMMENT 'Round number'
);
CREATE TABLE games(
    game_id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    white VARCHAR(255) NOT NULL COMMENT 'Name of the player using white pieces',
    black VARCHAR(255) NOT NULL COMMENT 'Name of the player using black pieces',
    result CHAR(1) NOT NULL COMMENT 'Result of the game (either B (black), W (white), T (tie))',
    FOREIGN KEY (event_id) REFERENCES events(event_id)
);
CREATE TABLE moves(
    move_id INT AUTO_INCREMENT PRIMARY KEY,
    game_id INT NOT NULL,
    move_num INT NOT NULL COMMENT 'Number of the move in the game',
    color CHAR(1) NOT NULL COMMENT 'Color of the piece moved (either B (black), W (white))',
    start_pos VARCHAR(2) NOT NULL COMMENT 'Starting postion of piece, ex : c4',
    ending_pos VARCHAR(2) NOT NULL COMMENT 'Ending postion of piece, ex : c5',
    kingside_castl TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'If the move was a kingside castling, this will be 1, otherwise, 0',
    queenside_castl TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'If the move was a queenside castling, this will be 1, otherwise, 0',
    FOREIGN KEY (game_id) REFERENCES games(game_id)
);