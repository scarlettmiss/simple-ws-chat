package movement

var Schema = `
CREATE TABLE movements (
    id uuid primary key ,
    user_id uuid,
    movement_type text,
    acceleration_x numeric,
    acceleration_y numeric,
    position_x numeric,
    position_y numeric,
    created_at timestamp
);
