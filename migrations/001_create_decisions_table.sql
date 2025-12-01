CREATE TABLE IF NOT EXISTS decisions (
    actor_user_id       VARCHAR(64) NOT NULL,
    recipient_user_id   VARCHAR(64) NOT NULL,
    liked_recipient     TINYINT(1)  NOT NULL,
    decision_unix_ts    BIGINT UNSIGNED NOT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (actor_user_id, recipient_user_id),
    INDEX idx_recipient_liked_ts (
        recipient_user_id,
        liked_recipient,
        decision_unix_ts DESC,
        actor_user_id
    ),
    INDEX idx_actor_recipient_liked (
        actor_user_id,
        recipient_user_id,
        liked_recipient
    )
);
