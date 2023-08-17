-- +goose Up
CREATE TABLE IF NOT EXISTS `orders` (
  `id` VARCHAR(255) NOT NULL,
  `price` DECIMAL(10,2) NULL,
  `tax` DECIMAL(10,2) NULL,
  `final_price` DECIMAL(10,2) NULL,
  PRIMARY KEY (`id`)
);

-- +goose Down
DROP TABLE IF EXISTS `orders`;