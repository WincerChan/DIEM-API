CREATE OR REPLACE FUNCTION public.randomfetch(len integer)
 RETURNS TABLE(data jsonb)
 LANGUAGE plpgsql
AS $function$
DECLARE
  amount int;
  lengths smallint;
BEGIN
  IF len < 2 OR len > 100 THEN
    amount := (SELECT reltuples::bigint FROM pg_catalog.pg_class WHERE relname = 'hitokoto');
    lengths := 32767;
  ELSE
    amount := (SELECT count(id) FROM hitokoto WHERE length < len);
    lengths := len;
  END IF;
  RETURN QUERY
    SELECT info FROM hitokoto WHERE length < lengths LIMIT 1 OFFSET (SELECT FLOOR(RANDOM() * amount));
END
$function$;
