-- +goose Up
CREATE TABLE building_floors (
	id bigint(20) NOT NULL AUTO_INCREMENT,
	floor varchar(150) NOT NULL,
	building_id bigint not null,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY fk_building_floors(building_id)
	REFERENCES buildings(id)
	ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
DROP TABLE building_floors;

