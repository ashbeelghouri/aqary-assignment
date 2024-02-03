WITH RankedStudents AS (
    SELECT id, student,
           ROW_NUMBER() OVER (ORDER BY id) AS RowNum
    FROM Seat
),
SwappedStudents AS (
    SELECT id,
           CASE
               WHEN RowNum % 2 = 0 THEN (
                   SELECT id FROM RankedStudents
                   WHERE RowNum = s.RowNum - 1
               )
               ELSE (
                   SELECT id FROM RankedStudents
                   WHERE RowNum = s.RowNum + 1
               )
           END AS student
    FROM RankedStudents s
)
SELECT id, student
FROM SwappedStudents
ORDER BY id;