INSERT INTO nickname
    VALUES (12111, 'B''z', TRUE);
COPY  (
    SELECT *
        FROM sample_table WHERE id = '123
'
) TO 'sample_dump'; INSERT INTO title VALUES (1, 'double " quote');