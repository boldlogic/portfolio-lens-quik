-- +goose Up
-- +goose StatementBegin
CREATE TABLE moex.security_boards (
    sec_id NVARCHAR(20) NOT NULL,
    board_id NVARCHAR(10) NOT NULL,
    is_traded BIT NOT NULL DEFAULT 0,
    is_primary BIT NOT NULL DEFAULT 0,
    PRIMARY KEY (sec_id, board_id),
    CONSTRAINT FK_security_boards_securities FOREIGN KEY (sec_id) REFERENCES moex.securities(sec_id),
    CONSTRAINT FK_security_boards_boards FOREIGN KEY (board_id) REFERENCES moex.boards(board_id)
);
CREATE INDEX IX_security_boards_board_id ON moex.security_boards(board_id);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS moex.security_boards;
