-- DB sqlite:test.db < query.sql

--select count(*) from habit;
SELECT 
    COUNT(*) AS total_number_of_habits,
    SUM(CASE WHEN is_completed = 1 THEN 1 ELSE 0 END) AS completed_habits,
    strftime('%d', day) AS day,
    strftime('%m', day) AS month,
    strftime('%Y', day) AS year
FROM habit
GROUP BY day, month, year
HAVING completed_habits > 0
ORDER BY year, month, day
