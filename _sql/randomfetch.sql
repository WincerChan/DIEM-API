CREATE OR REPLACE FUNCTION public.randomfetch(len integer, seed double precision)
 RETURNS TABLE(data jsonb)
 LANGUAGE plpgsql
AS $function$
DECLARE
  amount int;
BEGIN
  IF len < 2 OR len > 100 THEN
    amount := (SELECT max(count) FROM hito_len_count LIMIT 1);
  ELSE
    amount := (SELECT count FROM hito_len_count WHERE length=len LIMIT 1);
  END IF;
  RETURN QUERY
    SELECT info FROM hitokoto LIMIT 1 OFFSET (SELECT FLOOR(seed * amount));
END
$function$;
