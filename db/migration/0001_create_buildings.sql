-- +goose Up
CREATE TABLE buildings (
	id bigint(20) NOT NULL AUTO_INCREMENT,
	name varchar(150) NOT NULL,
	address varchar(255) DEFAULT NULL,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	KEY idx_buildings(name),
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
DROP TABLE buildings;

