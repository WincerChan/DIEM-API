set count 0
for v in (psql -p 5433 -U postgres -At -d api -c "SELECT COUNT(id) AS count, length AS length FROM hitokoto GROUP BY length ORDER BY length;")
    set data (string split "|" $v)
    set count (expr $data[1] + $count)
    psql -p 5433 -U postgres -At -d api -c "INSERT INTO hito_len_count (length, count) VALUES($data[2], $count);"
    # echo "len: $data[2], count: $count"
end
