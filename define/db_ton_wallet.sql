CREATE TABLE t_wallet_address (
    dauth_id varchar(255) NOT NULL DEFAULT '' COMMENT 'id of dauth',
		address varchar(255) NOT NULL DEFAULT '' COMMENT 'The user address ',
    secret varchar(255) NOT NULL DEFAULT '' COMMENT 'The user secret ',
    words varchar(255) NOT NULL DEFAULT '' COMMENT 'The user mnemonics' ,
		create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (dauth_id)
)ENGINE = InnoDB COLLATE utf8mb4_general_ci COMMENT 'wallet address table';