-- +goose Up

IF NOT EXISTS (
    SELECT 1
    FROM ref.instrument_attribute_priorities p
    JOIN ref.external_systems es ON es.ext_system_id = p.ext_system_id
    WHERE p.attribute_code = 'board.trade_point_id'
      AND es.ext_system = 'MOEX'
)
    INSERT INTO ref.instrument_attribute_priorities (attribute_code, ext_system_id, priority)
    SELECT 'board.trade_point_id', ext_system_id, 100
    FROM ref.external_systems
    WHERE ext_system = 'MOEX';

-- +goose Down

DELETE p
FROM ref.instrument_attribute_priorities p
JOIN ref.external_systems es ON es.ext_system_id = p.ext_system_id
WHERE p.attribute_code = 'board.trade_point_id'
  AND es.ext_system = 'MOEX';
