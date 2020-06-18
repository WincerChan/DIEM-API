DROP TABLE hito_len_count;
CREATE TABLE hito_len_count (
length INT PRIMARY KEY,
count INT
);
DROP TABLE hito;
CREATE TABLE hito (like hitokoto including all);
INSERT INTO hito SELECT * FROM hitokoto ORDER BY length;
DELETE FROM hitokoto;
INSERT INTO hitokoto SELECT * FROM hito ORDER BY length;
-- SELECT COUNT(id) AS count, length AS length FROM hitokoto GROUP BY length ORDER BY length;
